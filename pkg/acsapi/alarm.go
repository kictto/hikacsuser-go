package acsapi

import (
	"fmt"
	"strings"
	"sync"
	"unsafe"

	"github.com/clockworkchen/hikacsuser-go/internal/models"
	"github.com/clockworkchen/hikacsuser-go/internal/sdk"
)

// AlarmCallback 布控回调函数类型定义
type AlarmCallback func(alarmType int, alarmInfo interface{}) error

// AlarmSession 布控会话信息
type AlarmSession struct {
	DeviceID    string        // 设备标识（IP+端口）
	AlarmHandle int           // 布控句柄
	Callback    AlarmCallback // 回调函数
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
func (c *ACSClient) SetupAlarm(callback AlarmCallback, forceReplace bool) (int, error) {
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
		DeviceID:    deviceID,
		AlarmHandle: c.alarmManage.AlarmHandle,
		Callback:    callback,
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
	var alarmType int

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
			}

			// 处理图片数据
			if acsAlarmInfo.DwPicDataLen > 0 && acsAlarmInfo.PPicData != nil {
				exportAcsAlarmInfo.PicData = make([]byte, acsAlarmInfo.DwPicDataLen)
				slice := unsafe.Slice(acsAlarmInfo.PPicData, acsAlarmInfo.DwPicDataLen)
				copy(exportAcsAlarmInfo.PicData, slice)
			}

			exportAlarmInfo = exportAcsAlarmInfo
			alarmType = int(acsAlarmInfo.DwMajor)
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
			alarmType = int(alarmInfoV30.DwAlarmType)
		}
	default:
		// 其他类型报警，直接使用命令码作为报警类型
		alarmType = lCommand
	}

	// 调用用户定义的回调函数
	if session.Callback != nil && exportAlarmInfo != nil {
		err := session.Callback(alarmType, exportAlarmInfo)
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
