package acsapi

// #include <stdlib.h>
// #include <string.h>
import "C"
import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"unsafe"

	"github.com/clockworkchen/hikacsuser-go/internal/models"
	"github.com/clockworkchen/hikacsuser-go/internal/sdk"
)

// AlarmCallback 布控回调函数类型定义
type AlarmCallback func(lCommand int, alarmInfo interface{}) error

// AlarmSession 布控会话信息
type AlarmSession struct {
	DeviceID           string        // 设备标识（IP+端口）
	AlarmHandle        int           // 布控句柄
	Callback           AlarmCallback // 回调函数
	AutoDownloadPicUrl bool          // 是否自动下载图片URL数据
}

// 全局布控会话管理
var (
	alarmSessions     = make(map[string]*AlarmSession) // 设备ID -> 布控会话
	alarmSessionsLock sync.RWMutex                     // 会话锁
)

// SetupAlarm 设置布控
// callback: 布控回调函数，当有报警事件时会调用此函数
// forceReplace: 是否强制替换现有布控，true表示如果已存在布控则先关闭再重新布控，false表示如果已存在则跳过
// 返回值: 布控句柄和错误信息
func (c *ACSClient) SetupAlarm(callback AlarmCallback, forceReplace bool, autoDownloadPicUrl bool) (int, error) {
	if c.lUserID < 0 {
		return -1, fmt.Errorf("未登录设备")
	}

	// 生成设备ID
	deviceID := fmt.Sprintf("%s:%d", c.deviceIP, c.devicePort)

	// 设置SDK级别的报警回调函数
	// 使用NET_DVR_SetDVRMessageCallBack_V50设置回调函数
	// 索引使用0，这是SDK的默认索引
	if !c.hcnetsdk.NET_DVR_SetDVRMessageCallBack_V50(0, sdkMsgCallback, nil) {
		return -1, fmt.Errorf("设置SDK报警回调函数失败，错误码: %d", c.hcnetsdk.NET_DVR_GetLastError())
	}
	fmt.Println("SDK报警回调函数设置成功")

	// 检查是否已经布控
	alarmSessionsLock.RLock()
	existingSession, exists := alarmSessions[deviceID]
	alarmSessionsLock.RUnlock()

	if exists {
		// 已存在布控会话
		if !forceReplace {
			// 不强制替换，但更新回调函数
			alarmSessionsLock.Lock()
			existingSession.Callback = callback
			alarmSessionsLock.Unlock()
			return existingSession.AlarmHandle, nil
		}

		// 强制替换，先关闭现有布控
		err := c.CloseAlarm()
		if err != nil {
			return -1, fmt.Errorf("关闭现有布控失败: %v", err)
		}
	}

	// 创建报警布防管理实例（如果不存在）
	if c.alarmManage == nil {
		c.alarmManage = models.NewAlarmManage(c.hcnetsdk)
	}

	// 设置布控
	err := c.alarmManage.SetupAlarm(c.lUserID)
	if err != nil {
		return -1, fmt.Errorf("设置布控失败: %v", err)
	}

	// 保存布控会话
	session := &AlarmSession{
		DeviceID:           deviceID,
		AlarmHandle:        c.alarmManage.AlarmHandle,
		Callback:           callback,
		AutoDownloadPicUrl: autoDownloadPicUrl,
	}

	alarmSessionsLock.Lock()
	alarmSessions[deviceID] = session
	alarmSessionsLock.Unlock()

	return c.alarmManage.AlarmHandle, nil
}

// CloseAlarm 关闭布控
func (c *ACSClient) CloseAlarm() error {
	if c.lUserID < 0 {
		return fmt.Errorf("未登录设备")
	}

	// 生成设备ID
	deviceID := fmt.Sprintf("%s:%d", c.deviceIP, c.devicePort)

	// 检查是否存在布控会话
	alarmSessionsLock.RLock()
	_, exists := alarmSessions[deviceID]
	alarmSessionsLock.RUnlock()

	if !exists {
		return fmt.Errorf("设备未进行布控")
	}

	// 关闭布控
	if c.alarmManage == nil {
		return fmt.Errorf("报警管理模块未初始化")
	}

	err := c.alarmManage.CloseAlarm()
	if err != nil {
		return fmt.Errorf("关闭布控失败: %v", err)
	}

	// 移除布控会话
	alarmSessionsLock.Lock()
	delete(alarmSessions, deviceID)
	alarmSessionsLock.Unlock()

	return nil
}

// GetAlarmStatus 获取布控状态
// 返回值: 是否已布控、布控句柄
func (c *ACSClient) GetAlarmStatus() (bool, int) {
	if c.lUserID < 0 {
		return false, -1
	}

	// 生成设备ID
	deviceID := fmt.Sprintf("%s:%d", c.deviceIP, c.devicePort)

	// 检查是否存在布控会话
	alarmSessionsLock.RLock()
	session, exists := alarmSessions[deviceID]
	alarmSessionsLock.RUnlock()

	if !exists {
		return false, -1
	}

	return true, session.AlarmHandle
}

// sdkMsgCallback SDK级别的报警回调函数
// 该函数符合sdk.MSGCallBack_V31类型，用于接收SDK的报警消息
// 并根据设备ID查找对应的用户回调函数进行处理
func sdkMsgCallback(lCommand int, pAlarmer *sdk.NET_DVR_ALARMER, pAlarmInfo unsafe.Pointer, dwBufLen uint32, pUser unsafe.Pointer) bool {
	// 检查报警设备信息
	if pAlarmer == nil {
		fmt.Println("报警回调: 报警设备信息为空")
		return false
	}

	// 获取设备标识
	var deviceID string
	if pAlarmer.ByDeviceIPValid > 0 {
		// 使用设备IP作为标识
		deviceIP := sdk.BytesToString(pAlarmer.SDeviceIP[:])
		deviceID = deviceIP
	} else if pAlarmer.BySerialValid > 0 {
		// 使用设备序列号作为标识
		serialNumber := sdk.BytesToString(pAlarmer.SSerialNumber[:])
		deviceID = serialNumber
	} else if pAlarmer.ByUserIDValid > 0 {
		// 使用用户ID作为标识
		deviceID = fmt.Sprintf("UserID:%d", pAlarmer.LUserID)
	} else {
		fmt.Println("报警回调: 无法识别设备标识")
		return false
	}

	// 查找对应的布控会话
	alarmSessionsLock.RLock()
	var session *AlarmSession
	for _, s := range alarmSessions {
		if s.DeviceID == deviceID || strings.Contains(s.DeviceID, deviceID) {
			session = s
			break
		}
	}
	alarmSessionsLock.RUnlock()

	if session == nil {
		fmt.Printf("报警回调: 未找到设备 %s 的布控会话\n", deviceID)
		return false
	}

	// 根据报警类型处理报警信息
	var exportAlarmInfo interface{}

	switch lCommand {
	case sdk.COMM_ALARM_ACS: // 门禁主机报警
		if pAlarmInfo != nil && dwBufLen >= uint32(unsafe.Sizeof(sdk.NET_DVR_ACS_ALARM_INFO{})) {
			acsAlarmInfo := (*sdk.NET_DVR_ACS_ALARM_INFO)(pAlarmInfo)

			// 转换为导出类型
			exportAcsAlarmInfo := ACSAlarmInfo{
				Size:  acsAlarmInfo.DwSize,
				Major: acsAlarmInfo.DwMajor,
				Minor: acsAlarmInfo.DwMinor,
				Time: Time{
					Year:   acsAlarmInfo.StruTime.DwYear,
					Month:  acsAlarmInfo.StruTime.DwMonth,
					Day:    acsAlarmInfo.StruTime.DwDay,
					Hour:   acsAlarmInfo.StruTime.DwHour,
					Minute: acsAlarmInfo.StruTime.DwMinute,
					Second: acsAlarmInfo.StruTime.DwSecond,
				},
				NetUser:        sdk.BytesToString(acsAlarmInfo.SNetUser[:]),
				RemoteHostAddr: IPAddr{IPV4: sdk.BytesToString(acsAlarmInfo.StruRemoteHostAddr.SIpV4[:])},
				AcsEventInfo: ACSEventInfo{
					CardNo:            sdk.BytesToString(acsAlarmInfo.StruAcsEventInfo.ByCardNo[:]),
					CardType:          acsAlarmInfo.StruAcsEventInfo.ByCardType,
					WhiteListNo:       acsAlarmInfo.StruAcsEventInfo.ByWhiteListNo,
					ReportChannel:     acsAlarmInfo.StruAcsEventInfo.ByReportChannel,
					CardReaderKind:    acsAlarmInfo.StruAcsEventInfo.ByCardReaderKind,
					CardReaderNo:      acsAlarmInfo.StruAcsEventInfo.DwCardReaderNo,
					DoorNo:            acsAlarmInfo.StruAcsEventInfo.DwDoorNo,
					VerifyNo:          acsAlarmInfo.StruAcsEventInfo.DwVerifyNo,
					AlarmInNo:         acsAlarmInfo.StruAcsEventInfo.DwAlarmInNo,
					AlarmOutNo:        acsAlarmInfo.StruAcsEventInfo.DwAlarmOutNo,
					CaseSensorNo:      acsAlarmInfo.StruAcsEventInfo.DwCaseSensorNo,
					Rs485No:           acsAlarmInfo.StruAcsEventInfo.DwRs485No,
					MultiCardGroupNo:  acsAlarmInfo.StruAcsEventInfo.DwMultiCardGroupNo,
					AccessChannel:     acsAlarmInfo.StruAcsEventInfo.WAccessChannel,
					DeviceNo:          acsAlarmInfo.StruAcsEventInfo.ByDeviceNo,
					DistractControlNo: acsAlarmInfo.StruAcsEventInfo.ByDistractControlNo,
					EmployeeNo:        acsAlarmInfo.StruAcsEventInfo.DwEmployeeNo,
					LocalControllerID: acsAlarmInfo.StruAcsEventInfo.WLocalControllerID,
					InternetAccess:    acsAlarmInfo.StruAcsEventInfo.ByInternetAccess,
					Type:              acsAlarmInfo.StruAcsEventInfo.ByType,
				},
				PicDataLen:         acsAlarmInfo.DwPicDataLen,
				InductiveEventType: acsAlarmInfo.WInductiveEventType,
				PicTransType:       acsAlarmInfo.ByPicTransType,
				IOTChannelNo:       acsAlarmInfo.DwIOTChannelNo,
				TimeType:           acsAlarmInfo.ByTimeType,
			}

			// 处理图片数据
			if acsAlarmInfo.DwPicDataLen > 0 && acsAlarmInfo.PPicData != nil {
				var picData []byte

				// 检查是否是URL传输方式
				if acsAlarmInfo.ByPicTransType == 1 {
					// 对于URL方式，获取URL字符串
					urlBytes := make([]byte, 1024) // 假设URL不会超过1024字节

					// 从PPicData复制数据到urlBytes
					picPtr := unsafe.Pointer(acsAlarmInfo.PPicData)
					for i := 0; i < 1024; i++ {
						b := *(*byte)(unsafe.Pointer(uintptr(picPtr) + uintptr(i)))
						urlBytes[i] = b
						if b == 0 { // 找到字符串结束符
							break
						}
					}

					// 去除可能的空字节
					urlBytes = bytes.Trim(urlBytes, "\x00")
					if len(urlBytes) > 0 {
						urlStr := string(urlBytes)

						if len(urlStr) > 0 {
							exportAcsAlarmInfo.PicUri = urlStr
						}

						// 根据autoDownloadPicUrl参数决定是否下载图片
						if session.AutoDownloadPicUrl {
							// 从URL下载图片
							resp, httpErr := http.Get(urlStr)
							if httpErr == nil {
								defer resp.Body.Close()

								// 检查HTTP状态码
								if resp.StatusCode == http.StatusOK {
									// 读取图片数据
									picData, _ = io.ReadAll(resp.Body)
								}
							}
						}
					}
				} else {
					// 二进制数据方式
					// 尝试确定图片的实际大小
					// JPEG文件以FF D8开始，以FF D9结束
					maxSize := uint32(10 * 1024 * 1024) // 最大10MB
					tempBuf := make([]byte, maxSize)

					// 从PPicData复制数据到tempBuf
					picPtr := unsafe.Pointer(acsAlarmInfo.PPicData)
					var actualSize uint32
					for i := uint32(0); i < maxSize; i++ {
						tempBuf[i] = *(*byte)(unsafe.Pointer(uintptr(picPtr) + uintptr(i)))

						// 检查是否找到JPEG结束标记(FF D9)
						if i > 1 && tempBuf[i-1] == 0xFF && tempBuf[i] == 0xD9 {
							actualSize = i + 1 // 包括结束标记
							break
						}
					}

					if actualSize > 0 {
						picData = tempBuf[:actualSize]
					} else {
						// 如果无法确定大小，使用一个合理的固定大小
						fixedSize := uint32(100 * 1024) // 100KB
						picData = C.GoBytes(unsafe.Pointer(acsAlarmInfo.PPicData), C.int(fixedSize))

						// 尝试查找JPEG结束标记来截断数据
						for i := uint32(0); i < uint32(len(picData))-1; i++ {
							if picData[i] == 0xFF && picData[i+1] == 0xD9 {
								picData = picData[:i+2] // 截断到JPEG结束标记
								break
							}
						}
					}
				}

				if len(picData) > 0 {
					exportAcsAlarmInfo.PicData = picData
				}
			}

			// 处理扩展数据 - ACS_EVENT_INFO_EXTEND
			if acsAlarmInfo.ByAcsEventInfoExtend == 1 && acsAlarmInfo.PAcsEventInfoExtend != nil {
				// 解析扩展数据
				acsEventInfoExtend := (*sdk.NET_DVR_ACS_EVENT_INFO_EXTEND)(unsafe.Pointer(acsAlarmInfo.PAcsEventInfoExtend))

				// 转换为导出类型
				exportAcsEventInfoExtend := &ACSEventInfoExtend{
					FrontSerialNo:       acsEventInfoExtend.DwFrontSerialNo,
					UserType:            acsEventInfoExtend.ByUserType,
					CurrentVerifyMode:   acsEventInfoExtend.ByCurrentVerifyMode,
					CurrentEvent:        acsEventInfoExtend.ByCurrentEvent,
					PurePwdVerifyEnable: acsEventInfoExtend.ByPurePwdVerifyEnable,
					EmployeeNo:          sdk.BytesToString(bytes.Trim(acsEventInfoExtend.ByEmployeeNo[:], "\x00")), // 修正工号字段，去除可能的空字节
					AttendanceStatus:    acsEventInfoExtend.ByAttendanceStatus,
					StatusValue:         acsEventInfoExtend.ByStatusValue,
					UUID:                sdk.BytesToString(bytes.Trim(acsEventInfoExtend.ByUUID[:], "\x00")),       // 修正UUID字段，去除可能的空字节
					DeviceName:          sdk.BytesToString(bytes.Trim(acsEventInfoExtend.ByDeviceName[:], "\x00")), // 修正设备名称字段，去除可能的空字节
				}

				exportAcsAlarmInfo.AcsEventInfoExtend = exportAcsEventInfoExtend

				// 打印扩展信息日志
				fmt.Printf("ACS_EVENT_INFO_EXTEND: EmployeeNo=%s, UserType=%d, CurrentVerifyMode=%d\n",
					exportAcsEventInfoExtend.EmployeeNo, exportAcsEventInfoExtend.UserType, exportAcsEventInfoExtend.CurrentVerifyMode)
			}

			// 处理扩展数据 - ACS_EVENT_INFO_EXTEND_V20
			if acsAlarmInfo.ByAcsEventInfoExtendV20 == 1 && acsAlarmInfo.PAcsEventInfoExtendV20 != nil {
				// 解析扩展数据
				acsEventInfoExtendV20 := (*sdk.NET_DVR_ACS_EVENT_INFO_EXTEND_V20)(unsafe.Pointer(acsAlarmInfo.PAcsEventInfoExtendV20))

				// 转换为导出类型
				exportAcsEventInfoExtendV20 := &ACSEventInfoExtendV20{
					RemoteCheck:          acsEventInfoExtendV20.ByRemoteCheck,
					ThermometryUnit:      acsEventInfoExtendV20.ByThermometryUnit,
					IsAbnomalTemperature: acsEventInfoExtendV20.ByIsAbnomalTemperature,
					CurrTemperature:      acsEventInfoExtendV20.FCurrTemperature,
					RegionCoordinates: Point{
						X: acsEventInfoExtendV20.StruRegionCoordinates.FX,
						Y: acsEventInfoExtendV20.StruRegionCoordinates.FY,
					},
					XCoordinate:   acsEventInfoExtendV20.WXCoordinate,
					YCoordinate:   acsEventInfoExtendV20.WYCoordinate,
					Width:         acsEventInfoExtendV20.WWidth,
					Height:        acsEventInfoExtendV20.WHeight,
					HealthCode:    acsEventInfoExtendV20.ByHealthCode,
					NADCode:       acsEventInfoExtendV20.ByNADCode,
					TravelCode:    acsEventInfoExtendV20.ByTravelCode,
					VaccineStatus: acsEventInfoExtendV20.ByVaccineStatus,
				}

				// 输出温度信息和测温单位日志
				tempUnitStr := "未知"
				switch acsEventInfoExtendV20.ByThermometryUnit {
				case 0:
					tempUnitStr = "摄氏度"
				case 1:
					tempUnitStr = "华氏度"
				case 2:
					tempUnitStr = "开尔文"
				}

				fmt.Printf("ACS_EVENT_INFO_EXTEND_V20: 温度=%.1f%s, 是否异常=%d, 测温坐标=(%.3f,%.3f)\n",
					acsEventInfoExtendV20.FCurrTemperature, tempUnitStr,
					acsEventInfoExtendV20.ByIsAbnomalTemperature,
					acsEventInfoExtendV20.StruRegionCoordinates.FX,
					acsEventInfoExtendV20.StruRegionCoordinates.FY)

				// 处理二维码信息
				if acsEventInfoExtendV20.DwQRCodeInfoLen > 0 && acsEventInfoExtendV20.PQRCodeInfo != nil {
					qrCodeBytes := make([]byte, acsEventInfoExtendV20.DwQRCodeInfoLen)
					qrCodePtr := unsafe.Pointer(acsEventInfoExtendV20.PQRCodeInfo)
					for i := uint32(0); i < acsEventInfoExtendV20.DwQRCodeInfoLen; i++ {
						qrCodeBytes[i] = *(*byte)(unsafe.Pointer(uintptr(qrCodePtr) + uintptr(i)))
					}
					exportAcsEventInfoExtendV20.QRCodeInfo = string(bytes.Trim(qrCodeBytes, "\x00"))
					fmt.Printf("ACS_EVENT_INFO_EXTEND_V20: 解析二维码信息，长度=%d\n", acsEventInfoExtendV20.DwQRCodeInfoLen)
				}

				// 处理可见光图片数据
				if acsEventInfoExtendV20.DwVisibleLightDataLen > 0 && acsEventInfoExtendV20.PVisibleLightData != nil {
					visibleLightData := make([]byte, acsEventInfoExtendV20.DwVisibleLightDataLen)
					visibleLightPtr := unsafe.Pointer(acsEventInfoExtendV20.PVisibleLightData)
					for i := uint32(0); i < acsEventInfoExtendV20.DwVisibleLightDataLen; i++ {
						visibleLightData[i] = *(*byte)(unsafe.Pointer(uintptr(visibleLightPtr) + uintptr(i)))
					}
					exportAcsEventInfoExtendV20.VisibleLightData = visibleLightData
					fmt.Printf("ACS_EVENT_INFO_EXTEND_V20: 解析可见光图片数据，长度=%d\n", acsEventInfoExtendV20.DwVisibleLightDataLen)
				}

				// 处理热成像图片数据
				if acsEventInfoExtendV20.DwThermalDataLen > 0 && acsEventInfoExtendV20.PThermalData != nil {
					thermalData := make([]byte, acsEventInfoExtendV20.DwThermalDataLen)
					thermalPtr := unsafe.Pointer(acsEventInfoExtendV20.PThermalData)
					for i := uint32(0); i < acsEventInfoExtendV20.DwThermalDataLen; i++ {
						thermalData[i] = *(*byte)(unsafe.Pointer(uintptr(thermalPtr) + uintptr(i)))
					}
					exportAcsEventInfoExtendV20.ThermalData = thermalData
					fmt.Printf("ACS_EVENT_INFO_EXTEND_V20: 解析热成像图片数据，长度=%d\n", acsEventInfoExtendV20.DwThermalDataLen)
				}

				// 处理考勤自定义标签
				if len(acsEventInfoExtendV20.ByAttendanceLabel) > 0 {
					exportAcsEventInfoExtendV20.AttendanceLabel = sdk.BytesToString(bytes.Trim(acsEventInfoExtendV20.ByAttendanceLabel[:], "\x00"))
				}

				// 记录健康码、核酸码和疫苗状态信息
				healthCodeStatus := "未知"
				switch exportAcsEventInfoExtendV20.HealthCode {
				case 0:
					healthCodeStatus = "无效"
				case 1:
					healthCodeStatus = "未申领"
				case 2:
					healthCodeStatus = "未查询"
				case 3:
					healthCodeStatus = "绿码"
				case 4:
					healthCodeStatus = "黄码"
				case 5:
					healthCodeStatus = "红码"
				case 6:
					healthCodeStatus = "无此人员"
				case 7:
					healthCodeStatus = "查询失败"
				case 8:
					healthCodeStatus = "查询超时"
				}

				nadCodeStatus := "未知"
				switch exportAcsEventInfoExtendV20.NADCode {
				case 0:
					nadCodeStatus = "无效"
				case 1:
					nadCodeStatus = "未查询到"
				case 2:
					nadCodeStatus = "阴性(未过期)"
				case 3:
					nadCodeStatus = "阴性(已过期)"
				case 4:
					nadCodeStatus = "无效(已过期)"
				}

				travelCodeStatus := "未知"
				switch exportAcsEventInfoExtendV20.TravelCode {
				case 0:
					travelCodeStatus = "无效"
				case 1:
					travelCodeStatus = "14天内一直在低风险地区"
				case 2:
					travelCodeStatus = "14天内离开过低风险地区"
				case 3:
					travelCodeStatus = "14天内到达过中风险地区"
				case 4:
					travelCodeStatus = "其他"
				}

				vaccineStatus := "未知"
				switch exportAcsEventInfoExtendV20.VaccineStatus {
				case 0:
					vaccineStatus = "无效"
				case 1:
					vaccineStatus = "未接种"
				case 2:
					vaccineStatus = "接种中"
				case 3:
					vaccineStatus = "已完成接种"
				}

				fmt.Printf("ACS_EVENT_INFO_EXTEND_V20: 健康码=%s, 核酸码=%s, 行程码=%s, 疫苗状态=%s, 人脸坐标=(%d,%d), 宽高=(%d,%d)\n",
					healthCodeStatus, nadCodeStatus, travelCodeStatus, vaccineStatus,
					exportAcsEventInfoExtendV20.XCoordinate, exportAcsEventInfoExtendV20.YCoordinate,
					exportAcsEventInfoExtendV20.Width, exportAcsEventInfoExtendV20.Height)

				exportAcsAlarmInfo.AcsEventInfoExtendV20 = exportAcsEventInfoExtendV20
			}

			exportAlarmInfo = exportAcsAlarmInfo
			//alarmType = int(acsAlarmInfo.DwMajor)
		}
	case sdk.COMM_ALARM_V30: // 通用报警
		if pAlarmInfo != nil && dwBufLen >= uint32(unsafe.Sizeof(sdk.NET_DVR_ALARMINFO_V30{})) {
			alarmInfoV30 := (*sdk.NET_DVR_ALARMINFO_V30)(pAlarmInfo)

			// 转换为导出类型
			exportAlarmInfoV30 := AlarmInfoV30{
				Size:               alarmInfoV30.DwSize,
				AlarmType:          alarmInfoV30.DwAlarmType,
				AlarmInputNumber:   alarmInfoV30.DwAlarmInputNumber,
				AlarmOutputNumber:  make([]byte, len(alarmInfoV30.ByAlarmOutputNumber)),
				AlarmRelateChannel: make([]byte, len(alarmInfoV30.ByAlarmRelateChannel)),
				Channel:            make([]byte, len(alarmInfoV30.ByChannel)),
				DiskNumber:         make([]byte, len(alarmInfoV30.ByDiskNumber)),
			}

			// 复制数组数据
			copy(exportAlarmInfoV30.AlarmOutputNumber, alarmInfoV30.ByAlarmOutputNumber[:])
			copy(exportAlarmInfoV30.AlarmRelateChannel, alarmInfoV30.ByAlarmRelateChannel[:])
			copy(exportAlarmInfoV30.Channel, alarmInfoV30.ByChannel[:])
			copy(exportAlarmInfoV30.DiskNumber, alarmInfoV30.ByDiskNumber[:])

			exportAlarmInfo = exportAlarmInfoV30
			//alarmType = int(alarmInfoV30.DwAlarmType)
		}
	default:
		// 其他类型报警，直接使用命令码作为报警类型
	}

	// 调用用户定义的回调函数
	if session.Callback != nil && exportAlarmInfo != nil {
		err := session.Callback(lCommand, exportAlarmInfo)
		if err != nil {
			fmt.Printf("报警回调: 用户回调函数处理失败: %v\n", err)
			return false
		}
		return true
	}

	return false
}

// SearchAlarmEvent 查询报警事件
func (c *ACSClient) SearchAlarmEvent() error {
	if c.lUserID < 0 {
		return fmt.Errorf("未登录设备")
	}

	// 创建报警布防管理实例（如果不存在）
	if c.alarmManage == nil {
		c.alarmManage = models.NewAlarmManage(c.hcnetsdk)
	}

	return c.alarmManage.SearchAlarmEvent(c.lUserID)
}
