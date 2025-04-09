package main

// #include <stdlib.h>
// #include <string.h>
import "C"

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/hikacsuser-go/internal/models"
	"github.com/hikacsuser-go/internal/sdk"
	"github.com/hikacsuser-go/internal/utils"
	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

var (
	hcnetsdk sdk.HCNetSDK
	lUserID  = -1 // 用户句柄
)

// 初始化SDK
func initSDK() bool {
	// 获取SDK实例
	instance, err := sdk.GetSDKInstance()
	if err != nil {
		fmt.Printf("获取SDK实例失败: %v\n", err)
		return false
	}

	hcnetsdk = instance

	// 设置日志
	if !hcnetsdk.NET_DVR_SetLogToFile(3, "./sdklog", false) {
		fmt.Printf("设置日志失败: %d\n", hcnetsdk.NET_DVR_GetLastError())
		return false
	}

	fmt.Println("SDK初始化成功")
	return true
}

// 设备登录
func loginDevice(ip string, port uint16, username, password string) int {
	// 创建登录信息
	var loginInfo sdk.NET_DVR_USER_LOGIN_INFO
	var deviceInfo sdk.NET_DVR_DEVICEINFO_V40

	// 设置设备IP地址
	copy(loginInfo.SDeviceAddress[:], []byte(ip))

	// 传输模式设置
	loginInfo.ByUseTransport = 0 // 默认为0

	// 设置端口
	loginInfo.WPort = port

	// 设置用户名和密码
	copy(loginInfo.SUserName[:], []byte(username))
	copy(loginInfo.SPassword[:], []byte(password))

	// 登录模式设置
	loginInfo.BUseAsynLogin = false // 同步登录
	loginInfo.ByLoginMode = 0       // 使用SDK私有协议
	loginInfo.ByHttps = 0           // 不使用HTTPS

	// 执行登录
	userID := hcnetsdk.NET_DVR_Login_V40(&loginInfo, &deviceInfo)
	if userID == -1 {
		fmt.Printf("登录失败，错误码为: %d\n", hcnetsdk.NET_DVR_GetLastError())
	} else {
		fmt.Printf("%s 设备登录成功！\n", ip)
	}

	return userID
}

// 设备登出
func logout() {
	if lUserID >= 0 {
		if !hcnetsdk.NET_DVR_Logout(lUserID) {
			fmt.Printf("设备注销失败，错误码：%d\n", hcnetsdk.NET_DVR_GetLastError())
			return
		}
		fmt.Println("设备注销成功！！！")
	}
}

// 打印菜单
func printMenu() {
	fmt.Println("\n请输入您想要执行的demo实例! （退出请输入yes）")
	fmt.Println("1 - 获取门禁参数示例代码")
	fmt.Println("2 - 获取门禁状态示例代码")
	fmt.Println("3 - 远程控门示例代码")
	fmt.Println("4 - 下发人员示例代码")
	fmt.Println("5 - 查询人员示例代码")
	fmt.Println("6 - 删除人员代码")
	fmt.Println("7 - 下发卡号代码")
	fmt.Println("8 - 查询卡号代码")
	fmt.Println("9 - 删除卡号代码")
	fmt.Println("10 - 二进制方式下发人脸代码")
	fmt.Println("12 - URL方式下发人脸代码")
	fmt.Println("13 - 查询人脸代码")
	fmt.Println("14 - 删除人脸代码")
	fmt.Println("15 - 采集人脸代码")
	fmt.Println("16 - 门禁历史事件查询代码")
	fmt.Println("17 - 设置计划模板代码")
	fmt.Println("18 - 报警布防回调演示")
	fmt.Println("--------------------------------------------------")
}

func main() {
	// 初始化SDK
	if !initSDK() {
		fmt.Println("初始化SDK失败")
		return
	}

	// 登录设备
	lUserID = loginDevice("192.168.1.5", 8000, "admin", "zswlzswl2025")
	if lUserID == -1 {
		fmt.Println("登录设备失败")
		return
	}

	// 读取用户输入
	scanner := bufio.NewScanner(os.Stdin)

	for {
		printMenu()

		// 读取用户输入
		scanner.Scan()
		input := scanner.Text()
		input = strings.ToLower(input)

		if input == "yes" {
			break
		}

		// 处理用户选择
		handleUserChoice(input)
	}

	// 休眠1秒，确保程序完全执行
	time.Sleep(time.Second)

	// 注销设备
	logout()

	// 释放SDK资源
	hcnetsdk.NET_DVR_Cleanup()
}

// 处理用户选择
func handleUserChoice(choice string) {
	// 创建各个管理模块的实例
	acsManage := models.NewACSManage(hcnetsdk)
	userManage := models.NewUserManage(hcnetsdk)
	cardManage := models.NewCardManage(hcnetsdk)
	faceManage := models.NewFaceManage(hcnetsdk)
	eventSearch := models.NewEventSearch(hcnetsdk)

	// 将选择转换为数字
	num, err := strconv.Atoi(choice)
	if err != nil {
		fmt.Println("\n未知的指令操作!请重新输入!\n")
		return
	}

	switch num {
	case 1:
		fmt.Println("\n[Module]获取门禁参数示例代码")
		if err := acsManage.AcsCfg(lUserID); err != nil {
			fmt.Printf("获取门禁参数失败: %v\n", err)
		}
	case 2:
		fmt.Println("\n[Module]获取门禁状态示例代码")
		if err := acsManage.GetAcsStatus(lUserID); err != nil {
			fmt.Printf("获取门禁状态失败: %v\n", err)
		}
	case 3:
		fmt.Println("\n[Module]远程控门示例代码")
		if err := acsManage.RemoteControlGate(lUserID); err != nil {
			fmt.Printf("远程控门失败: %v\n", err)
		}
	case 4:
		fmt.Println("\n[Module]下发人员示例代码")
		if err := userManage.AddUserInfo(lUserID, "12345"); err != nil {
			fmt.Printf("下发人员失败: %v\n", err)
		}
	case 5:
		fmt.Println("\n[Module]查询人员示例代码")
		if err := userManage.SearchUserInfo(lUserID); err != nil {
			fmt.Printf("查询人员失败: %v\n", err)
		}
	case 6:
		fmt.Println("\n[Module]删除人员代码")
		if err := userManage.DeleteUserInfo(lUserID); err != nil {
			fmt.Printf("删除人员失败: %v\n", err)
		}
	case 7:
		fmt.Println("\n[Module]下发卡号代码")
		if err := cardManage.AddCardInfo(lUserID, "12345", "12345"); err != nil {
			fmt.Printf("下发卡号失败: %v\n", err)
		}
	case 8:
		fmt.Println("\n[Module]查询卡号代码")
		if err := cardManage.SearchCardInfo(lUserID, "1"); err != nil {
			fmt.Printf("查询卡号失败: %v\n", err)
		}
	case 9:
		fmt.Println("\n[Module]删除卡号代码")
		if err := cardManage.DeleteCardInfo(lUserID, "12345"); err != nil {
			fmt.Printf("删除卡号失败: %v\n", err)
		}
	case 10:
		fmt.Println("\n[Module]二进制方式下发人脸代码")
		if err := faceManage.AddFaceByBinary(lUserID, "1"); err != nil {
			fmt.Printf("下发人脸失败: %v\n", err)
		}
	case 12:
		fmt.Println("\n[Module]URL方式下发人脸代码")
		if err := faceManage.AddFaceByUrl(lUserID, "12345"); err != nil {
			fmt.Printf("下发人脸失败: %v\n", err)
		}
	case 13:
		fmt.Println("\n[Module]查询人脸代码")
		if err := faceManage.SearchFaceInfo(lUserID, "12345"); err != nil {
			fmt.Printf("查询人脸失败: %v\n", err)
		}
	case 14:
		fmt.Println("\n[Module]删除人脸代码")
		if err := faceManage.DeleteFaceInfo(lUserID, "12345"); err != nil {
			fmt.Printf("删除人脸失败: %v\n", err)
		}
	case 15:
		fmt.Println("\n[Module]采集人脸代码")
		if err := faceManage.CaptureFaceInfo(lUserID); err != nil {
			fmt.Printf("采集人脸失败: %v\n", err)
		}
	case 16:
		fmt.Println("\n[Module]门禁历史事件查询代码")
		if err := eventSearch.SearchAllEvent(lUserID); err != nil {
			fmt.Printf("查询事件失败: %v\n", err)
		}
	case 17:
		fmt.Println("\n[Module]设置计划模板代码")
		if err := userManage.SetCardTemplate(lUserID, 2); err != nil {
			fmt.Printf("设置计划模板失败: %v\n", err)
		}
	case 18:
		fmt.Println("\n[Module]报警布防回调演示")
		handleAlarmSetup(lUserID)
	default:
		fmt.Println("\n未知的指令操作!请重新输入!\n")
	}
}

// MsgCallback 报警回调函数
func MsgCallback(lCommand int, pAlarmer *sdk.NET_DVR_ALARMER, pAlarmInfo unsafe.Pointer, dwBufLen uint32, pUser unsafe.Pointer) bool {
	fmt.Printf("\n--------------------\n")
	fmt.Printf("报警回调触发!\n")
	fmt.Printf("报警命令: %d (0x%X)\n", lCommand, lCommand)

	// 打印报警设备信息
	if pAlarmer != nil {
		fmt.Printf("报警设备信息:\n")
		fmt.Printf("  用户ID有效: %d\n", pAlarmer.ByUserIDValid)
		if pAlarmer.ByUserIDValid > 0 {
			fmt.Printf("  用户ID: %d\n", pAlarmer.LUserID)
		}
		fmt.Printf("  序列号有效: %d\n", pAlarmer.BySerialValid)
		if pAlarmer.BySerialValid > 0 {
			serialNumberBytes := bytes.Trim(pAlarmer.SSerialNumber[:], "\x00")
			serialNumber := string(serialNumberBytes)
			fmt.Printf("  序列号: %s\n", serialNumber)
		}
		fmt.Printf("  设备IP有效: %d\n", pAlarmer.ByDeviceIPValid)
		if pAlarmer.ByDeviceIPValid > 0 {
			deviceIPBytes := bytes.Trim(pAlarmer.SDeviceIP[:], "\x00")
			deviceIP := string(deviceIPBytes)
			fmt.Printf("  设备IP: %s\n", deviceIP)
		}
		fmt.Printf("  设备名称有效: %d\n", pAlarmer.ByDeviceNameValid)
		if pAlarmer.ByDeviceNameValid > 0 {
			deviceNameBytes := bytes.Trim(pAlarmer.SDeviceName[:], "\x00")
			deviceName := string(deviceNameBytes)
			fmt.Printf("  设备名称: %s\n", deviceName)
		}
		fmt.Printf("  MAC地址有效: %d\n", pAlarmer.ByMacAddrValid)
		if pAlarmer.ByMacAddrValid > 0 {
			fmt.Printf("  MAC地址: %02X:%02X:%02X:%02X:%02X:%02X\n",
				pAlarmer.ByMacAddr[0], pAlarmer.ByMacAddr[1], pAlarmer.ByMacAddr[2],
				pAlarmer.ByMacAddr[3], pAlarmer.ByMacAddr[4], pAlarmer.ByMacAddr[5])
		}
	}

	// 处理不同类型的报警信息 (根据 lCommand 解析 pAlarmInfo)
	fmt.Printf("报警信息长度: %d\n", dwBufLen)

	switch lCommand {
	case sdk.COMM_ALARM_V30:
		fmt.Println("  报警类型: COMM_ALARM_V30 (通用报警)")
		if pAlarmInfo != nil && dwBufLen >= uint32(unsafe.Sizeof(sdk.NET_DVR_ALARMINFO_V30{})) {
			alarmInfoV30 := (*sdk.NET_DVR_ALARMINFO_V30)(pAlarmInfo)
			fmt.Printf("    报警类型代码: %d\n", alarmInfoV30.DwAlarmType)
			fmt.Printf("    报警输入号: %d\n", alarmInfoV30.DwAlarmInputNumber)
			// 根据 alarmInfoV30.DwAlarmType 进一步解析其他信息
			// 例如：打印触发的通道、硬盘号等
		} else {
			fmt.Println("    报警信息数据无效或长度不足")
		}

	case sdk.COMM_ALARM_ACS:
		fmt.Println("  报警类型: COMM_ALARM_ACS (门禁事件)")
		if pAlarmInfo != nil && dwBufLen >= uint32(unsafe.Sizeof(sdk.NET_DVR_ACS_ALARM_INFO{})) {
			acsAlarmInfo := (*sdk.NET_DVR_ACS_ALARM_INFO)(pAlarmInfo)

			// 检查结构体大小是否正确
			if acsAlarmInfo.DwSize != uint32(unsafe.Sizeof(*acsAlarmInfo)) {
				fmt.Printf("    警告: 结构体大小不匹配，收到: %d, 预期: %d\n",
					acsAlarmInfo.DwSize, unsafe.Sizeof(*acsAlarmInfo))
			}

			// 获取报警主类型描述
			majorDesc := utils.GetAlarmMajorTypeDesc(int(acsAlarmInfo.DwMajor))
			// 获取报警次类型描述
			minorDesc := utils.GetAlarmMinorTypeDesc(int(acsAlarmInfo.DwMajor), int(acsAlarmInfo.DwMinor))

			fmt.Printf("    报警主类型: %d (0x%X) - %s\n", acsAlarmInfo.DwMajor, acsAlarmInfo.DwMajor, majorDesc)
			fmt.Printf("    报警次类型: %d (0x%X) - %s\n", acsAlarmInfo.DwMinor, acsAlarmInfo.DwMinor, minorDesc)

			// 打印报警时间
			t := acsAlarmInfo.StruTime
			fmt.Printf("    报警时间: %d-%02d-%02d %02d:%02d:%02d\n",
				t.DwYear, t.DwMonth, t.DwDay, t.DwHour, t.DwMinute, t.DwSecond)

			// 打印事件详情
			eventInfo := acsAlarmInfo.StruAcsEventInfo
			fmt.Printf("原始卡号字节str数据: %v\n", string(eventInfo.ByCardNo[:]))
			fmt.Printf("原始卡号字节数组数据: %v\n", eventInfo.ByCardNo)
			cardNoBytes := bytes.Trim(eventInfo.ByCardNo[:], "\x00")
			fmt.Printf("原始卡号字节数据剪切过: %v\n", cardNoBytes)
			fmt.Printf("原始卡号(十六进制): %X\n", cardNoBytes)
			cardNo := string(cardNoBytes)
			fmt.Printf("原始卡号(剪切后字符串): %s\n", cardNo)
			if len(cardNo) > 0 {
				fmt.Printf("    卡号: %s\n", cardNo)
			}
			fmt.Printf("    卡类型: %d\n", eventInfo.ByCardType)
			fmt.Printf("    读卡器编号: %d\n", eventInfo.DwCardReaderNo)
			fmt.Printf("    门编号: %d\n", eventInfo.DwDoorNo)
			if eventInfo.DwEmployeeNo > 0 {
				fmt.Printf("    工号: %d\n", eventInfo.DwEmployeeNo)
			}
			fmt.Printf("    继续打印是否有图片: acsAlarmInfo.DwPicDataLen=%d\n", acsAlarmInfo.DwPicDataLen)
			fmt.Printf("    图片传输类型(ByPicTransType): %d (0-二进制数据, 1-URL)\n", acsAlarmInfo.ByPicTransType)
			// 根据 DwMajor 和 DwMinor 提供更具体的事件描述 (需要查阅文档)
			// 例如：
			//  if acsAlarmInfo.DwMajor == sdk.MAJOR_ALARM && acsAlarmInfo.DwMinor == sdk.MINOR_LEGAL_CARD_PASS {
			//  	fmt.Println("    事件描述: 合法卡通过")
			//  }

			// 如果有图片信息，可以保存或处理
			if acsAlarmInfo.DwPicDataLen > 0 && acsAlarmInfo.PPicData != nil {
				fmt.Printf("    包含图片数据，长度: %d，指针地址: %p\n", acsAlarmInfo.DwPicDataLen, acsAlarmInfo.PPicData)

				// 拷贝图片数据
				// 注意：这里使用了三种方法处理指针，如果一种失败可以尝试另一种
				var picData []byte

				// 方法1: 使用GoBytes直接复制指针内容到Go slice
				picData = C.GoBytes(unsafe.Pointer(acsAlarmInfo.PPicData), C.int(acsAlarmInfo.DwPicDataLen))

				// 方法2: 如果方法1失败，尝试分配内存并使用memcpy复制
				if len(picData) == 0 && acsAlarmInfo.DwPicDataLen > 0 {
					fmt.Println("    方法1获取图片数据失败，尝试方法2")
					// 分配Go内存
					picData = make([]byte, acsAlarmInfo.DwPicDataLen)
					// 使用C.memcpy复制内存
					C.memcpy(unsafe.Pointer(&picData[0]), unsafe.Pointer(acsAlarmInfo.PPicData), C.size_t(acsAlarmInfo.DwPicDataLen))
				}

				// 方法3: 如果方法1和方法2都失败，尝试使用unsafe指针操作
				if len(picData) == 0 && acsAlarmInfo.DwPicDataLen > 0 {
					fmt.Println("    方法2获取图片数据失败，尝试方法3")
					picData = make([]byte, acsAlarmInfo.DwPicDataLen)

					// 获取内存起始地址
					picPtr := unsafe.Pointer(acsAlarmInfo.PPicData)

					// 手动按字节复制
					picSlice := (*[1 << 30]byte)(picPtr)[:acsAlarmInfo.DwPicDataLen:acsAlarmInfo.DwPicDataLen]
					copy(picData, picSlice)
				}

				if len(picData) == 0 {
					fmt.Println("    警告: 无法从指针获取图片数据!")
				} else {
					// 创建带时间戳的文件名
					timestamp := time.Now().UnixNano()
					picFilename := fmt.Sprintf("alarm_pic_%d.jpg", timestamp)

					// 保存图片数据到文件
					err := os.WriteFile(picFilename, picData, 0644)
					if err == nil {
						fmt.Printf("    图片已保存为: %s (数据大小: %d字节)\n", picFilename, len(picData))
						// 打印图片数据的前几个字节，用于调试
						if len(picData) > 16 {
							fmt.Printf("    图片数据头: %X\n", picData[:16])
						}
					} else {
						fmt.Printf("    保存图片失败: %v\n", err)
					}
				}
			} else if acsAlarmInfo.DwPicDataLen > 0 && acsAlarmInfo.PPicData == nil {
				fmt.Println("    图片数据指针为空，但DwPicDataLen > 0")

				// 处理URL方式(ByPicTransType=1)的图片传输
				if acsAlarmInfo.ByPicTransType == 1 {
					fmt.Println("    图片为URL传输方式，尝试获取URL...")

					// 对于URL方式，我们可能需要从其他字段获取URL
					// 根据SDK文档，通常是从pPicData或其他字段获取URL字符串
					// 这里我们尝试从不同的地方获取URL

					if acsAlarmInfo.PPicData != nil {
						urlBytes := C.GoBytes(unsafe.Pointer(acsAlarmInfo.PPicData), C.int(acsAlarmInfo.DwPicDataLen))
						// 去除可能的空字节
						urlBytes = bytes.Trim(urlBytes, "\x00")
						if len(urlBytes) > 0 {
							urlStr := string(urlBytes)
							fmt.Printf("    获取到图片URL: %s\n", urlStr)
							// 这里可以添加代码从URL下载图片
						} else {
							fmt.Println("    无法获取有效的URL")
						}
					}
				}

				// 特殊处理DwPicDataLen=1的情况
				if acsAlarmInfo.DwPicDataLen == 1 {
					fmt.Println("    检测到DwPicDataLen=1，这可能表示需要通过查找图片接口获取图片")

					// 获取报警时间，并计算前后5分钟的时间范围用于搜索图片
					alarmTime := acsAlarmInfo.StruTime

					// 创建查找图片参数结构
					var findPicture sdk.NET_DVR_FIND_PICTURE_PARAM
					findPicture.DwSize = uint32(unsafe.Sizeof(findPicture))

					// 设置通道号 - 一般门禁事件使用1或门编号
					findPicture.LChannel = 1

					// 设置图片类型 - 门禁事件一般使用0xFF (所有类型)
					findPicture.ByFileType = 0xFF

					// 设置查找时间范围 (报警时间前后5分钟)
					startTime := alarmTime
					endTime := alarmTime

					// 开始时间设为报警时间前5分钟
					if startTime.DwMinute >= 5 {
						startTime.DwMinute -= 5
					} else {
						// 处理分钟减法的借位
						if startTime.DwHour > 0 {
							startTime.DwHour -= 1
							startTime.DwMinute = startTime.DwMinute + 60 - 5
						} else {
							// 如果小时也是0，不再往前推
							startTime.DwMinute = 0
						}
					}

					// 结束时间设为报警时间后5分钟
					if endTime.DwMinute <= 55 {
						endTime.DwMinute += 5
					} else {
						// 处理分钟加法的进位
						endTime.DwMinute = (endTime.DwMinute + 5) % 60
						if endTime.DwHour < 23 {
							endTime.DwHour += 1
						} else {
							// 如果已经是23点，不再往后推
							endTime.DwHour = 23
							endTime.DwMinute = 59
						}
					}

					// 复制开始和结束时间到查找参数
					findPicture.StruStartTime = startTime
					findPicture.StruStopTime = endTime

					// 如果有卡号，设置卡号作为查找条件
					cardNo := string(bytes.Trim(acsAlarmInfo.StruAcsEventInfo.ByCardNo[:], "\x00"))
					if len(cardNo) > 0 {
						fmt.Printf("    使用卡号 %s 查找相关图片\n", cardNo)
						// 复制卡号到查找参数 (如果SDK支持)
						// copy(findPicture.SCardNum[:], []byte(cardNo))
					}

					// 开始查找图片
					lFindHandle := hcnetsdk.NET_DVR_FindPicture(lUserID, &findPicture)
					if lFindHandle < 0 {
						fmt.Printf("    查找图片失败，错误码: %d\n", hcnetsdk.NET_DVR_GetLastError())
					} else {
						fmt.Printf("    开始查找图片，句柄: %d\n", lFindHandle)

						// 循环获取查找结果
						maxAttempts := 20 // 最多尝试20次
						for i := 0; i < maxAttempts; i++ {
							var findData sdk.NET_DVR_FIND_PICTURE
							ret := hcnetsdk.NET_DVR_FindNextPicture(lFindHandle, &findData)

							switch ret {
							case sdk.NET_DVR_FILE_SUCCESS:
								// 找到图片
								fileName := string(bytes.Trim(findData.SFileName[:], "\x00"))
								fmt.Printf("    找到图片: %s\n", fileName)

								// 准备接收图片数据的参数
								var getPicParam sdk.NET_DVR_GETPIC_PARAM
								getPicParam.DwSize = uint32(unsafe.Sizeof(getPicParam))

								// 设置图片保存方式 - 0: 二进制数据
								getPicParam.ByPictype = 0

								// 创建图片文件名
								picFilename := fmt.Sprintf("alarm_pic_%d_%d.jpg", time.Now().UnixNano(), i)

								// 将Go字符串转换为C字符串，用于设置文件名
								cFilename := C.CString(picFilename)
								defer C.free(unsafe.Pointer(cFilename))

								// 设置保存文件名
								getPicParam.PicName = (*byte)(unsafe.Pointer(cFilename))

								// 获取图片
								if hcnetsdk.NET_DVR_GetPicture_V50(lUserID, &findData, &getPicParam) {
									fmt.Printf("    成功保存图片到: %s\n", picFilename)
								} else {
									fmt.Printf("    保存图片失败，错误码: %d\n", hcnetsdk.NET_DVR_GetLastError())
								}

							case sdk.NET_DVR_ISFINDING:
								// 正在查找，继续等待
								fmt.Println("    正在查找图片，请等待...")
								time.Sleep(200 * time.Millisecond)
								continue

							case sdk.NET_DVR_FILE_NOFIND:
								// 未找到更多图片
								fmt.Println("    未找到更多图片")
								break

							case sdk.NET_DVR_NOMOREFILE:
								// 没有更多图片
								fmt.Println("    没有更多图片")
								break

							case sdk.NET_DVR_FILE_EXCEPTION:
								// 查找图片异常
								fmt.Printf("    查找图片异常，错误码: %d\n", hcnetsdk.NET_DVR_GetLastError())
								break

							default:
								fmt.Printf("    查找图片返回未知状态: %d\n", ret)
								break
							}

							// 如果不是正在查找状态，跳出循环
							if ret != sdk.NET_DVR_ISFINDING {
								break
							}
						}

						// 关闭查找句柄
						if !hcnetsdk.NET_DVR_CloseFindPicture(lFindHandle) {
							fmt.Printf("    关闭查找图片句柄失败，错误码: %d\n", hcnetsdk.NET_DVR_GetLastError())
						}
					}
				}
			} else {
				if acsAlarmInfo.DwPicDataLen == 0 {
					fmt.Println("    不包含图片数据 (数据长度为0)")
				} else {
					fmt.Println("    图片数据指针为空")
				}
			}

		} else {
			fmt.Println("    报警信息数据无效或长度不足")
		}

	// 在这里添加更多 case 来处理其他报警类型...
	// case sdk.COMM_ALARM_MOTION:
	//    ...

	default:
		fmt.Printf("  未处理的报警类型: %d\n", lCommand)
		// 可以选择将原始数据打印出来以便调试
		if pAlarmInfo != nil && dwBufLen > 0 {
			fmt.Printf("    报警信息数据长度: %d\n", dwBufLen)
			//rawData := C.GoBytes(pAlarmInfo, C.int(dwBufLen))
			//fmt.Printf("    原始数据 (前%d字节): %X\n", min(int(dwBufLen), 128), rawData[:min(int(dwBufLen), 128)])
		}
	}

	fmt.Printf("--------------------\n\n")

	// 返回 true 表示成功处理了该消息
	return true
}

// handleAlarmSetup 处理报警布防逻辑
func handleAlarmSetup(lUserID int) {
	if lUserID < 0 {
		fmt.Println("请先登录设备")
		return
	}

	// 1. 设置报警回调函数
	// 将 MsgCallback 函数注册给 SDK，索引使用 0，pUser 传递 nil
	if !hcnetsdk.NET_DVR_SetDVRMessageCallBack_V50(0, MsgCallback, nil) {
		fmt.Printf("设置报警回调函数失败，错误码: %d\n", hcnetsdk.NET_DVR_GetLastError())
		return
	}
	fmt.Println("报警回调函数设置成功")

	// 2. 设置布防参数
	var setupParam sdk.NET_DVR_SETUPALARM_PARAM
	setupParam.DwSize = uint32(unsafe.Sizeof(setupParam))
	setupParam.ByLevel = 1         // 布防优先级：中
	setupParam.ByAlarmInfoType = 1 // 上传报警信息类型：V30结构
	// 根据需要设置其他布防参数...

	// 3. 建立报警上传通道 (布防)
	lAlarmHandle := hcnetsdk.NET_DVR_SetupAlarmChan_V41(lUserID, &setupParam)
	if lAlarmHandle < 0 {
		fmt.Printf("布防失败，错误码: %d\n", hcnetsdk.NET_DVR_GetLastError())
		// 如果布防失败，最好也注销回调函数（虽然这里省略了）
		return
	}
	fmt.Printf("布防成功, 报警句柄: %d\n", lAlarmHandle)

	// 4. 等待报警触发 (实际应用中这里应该是阻塞或异步等待)
	fmt.Println("布防已启动，等待报警触发... (按 Enter 键撤防并退出)")
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n') // 等待用户输入回车

	// 5. 关闭报警上传通道 (撤防)
	if !hcnetsdk.NET_DVR_CloseAlarmChan_V30(lAlarmHandle) {
		fmt.Printf("撤防失败，错误码: %d\n", hcnetsdk.NET_DVR_GetLastError())
	} else {
		fmt.Println("撤防成功")
	}

	// 注销回调函数 (可选，取决于应用逻辑)
	// hcnetsdk.NET_DVR_SetDVRMessageCallBack_V50(0, nil, nil)
	// fmt.Println("报警回调函数已注销")
}
