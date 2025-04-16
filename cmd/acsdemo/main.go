package main

// #include <stdlib.h>
// #include <string.h>
import "C"

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/clockworkchen/hikacsuser-go/internal/models"
	"github.com/clockworkchen/hikacsuser-go/internal/sdk"
	"github.com/clockworkchen/hikacsuser-go/internal/utils"
	"io"
	"net/http"
	"os"
	"path/filepath"
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
			// 根据SDK文档，DwPicDataLen=1只表示有图片，而不是图片的实际大小
			if acsAlarmInfo.DwPicDataLen > 0 && acsAlarmInfo.PPicData != nil {
				fmt.Printf("    包含图片数据，指针地址: %p\n", acsAlarmInfo.PPicData)

				// 拷贝图片数据
				// 注意：这里不能使用DwPicDataLen作为图片大小，因为它可能只是一个标志值
				var picData []byte
				var actualSize uint32

				// 首先检查是否是URL传输方式
				if acsAlarmInfo.ByPicTransType == 1 {
					fmt.Println("    图片为URL传输方式，尝试获取URL...")

					// 对于URL方式，我们需要获取URL字符串
					// 尝试获取URL字符串，这里需要找到字符串的结束位置
					// 假设URL是以null结尾的C字符串
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
						fmt.Printf("    获取到图片URL: %s\n", urlStr)

						// 从URL下载图片
						fmt.Println("    正在从URL下载图片...")
						resp, httpErr := http.Get(urlStr)
						if httpErr != nil {
							fmt.Printf("    下载图片失败: %v\n", httpErr)
						} else {
							defer resp.Body.Close()

							// 检查HTTP状态码
							if resp.StatusCode != http.StatusOK {
								fmt.Printf("    下载图片失败，HTTP状态码: %d\n", resp.StatusCode)
							} else {
								// 读取图片数据
								picData, readErr := io.ReadAll(resp.Body)
								if readErr != nil {
									fmt.Printf("    读取图片数据失败: %v\n", readErr)
								} else if len(picData) > 0 {
									// 创建带时间戳的文件名
									timestamp := time.Now().UnixNano()
									picFilename := fmt.Sprintf("alarm_pic_%d.jpg", timestamp)

									// 保存图片数据到文件
									saveErr := os.WriteFile(picFilename, picData, 0644)
									if saveErr == nil {
										// 获取文件的绝对路径
										absPath, pathErr := filepath.Abs(picFilename)
										if pathErr != nil {
											absPath = picFilename // 如果获取绝对路径失败，使用相对路径
										}
										fmt.Printf("    URL图片已下载并保存为: %s (数据大小: %d字节)\n", absPath, len(picData))
									} else {
										fmt.Printf("    保存下载的图片失败: %v\n", saveErr)
									}
								} else {
									fmt.Println("    下载的图片数据为空")
								}
							}
						}
					} else {
						fmt.Println("    无法获取有效的URL")
					}
				} else {
					// 二进制数据方式
					fmt.Println("    图片为二进制数据方式")

					// 尝试确定图片的实际大小
					// 方法1：检查JPEG文件头和尾部标记
					// JPEG文件以FF D8开始，以FF D9结束
					// 先读取一个较大的缓冲区，然后查找JPEG结束标记
					maxSize := uint32(10 * 1024 * 1024) // 最大10MB，防止无限读取
					tempBuf := make([]byte, maxSize)

					// 从PPicData复制数据到tempBuf
					picPtr := unsafe.Pointer(acsAlarmInfo.PPicData)
					for i := uint32(0); i < maxSize; i++ {
						tempBuf[i] = *(*byte)(unsafe.Pointer(uintptr(picPtr) + uintptr(i)))

						// 检查是否找到JPEG结束标记(FF D9)
						if i > 1 && tempBuf[i-1] == 0xFF && tempBuf[i] == 0xD9 {
							actualSize = i + 1 // 包括结束标记
							break
						}
					}

					if actualSize > 0 {
						fmt.Printf("    检测到JPEG图片，实际大小: %d字节\n", actualSize)
						picData = tempBuf[:actualSize]
					} else {
						fmt.Println("    无法确定图片大小，尝试使用固定大小读取")
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
					// 创建带时间戳的文件名
					timestamp := time.Now().UnixNano()
					picFilename := fmt.Sprintf("alarm_pic_%d.jpg", timestamp)

					// 保存图片数据到文件
					err := os.WriteFile(picFilename, picData, 0644)
					if err == nil {
						// 获取文件的绝对路径
						absPath, pathErr := filepath.Abs(picFilename)
						if pathErr != nil {
							absPath = picFilename // 如果获取绝对路径失败，使用相对路径
						}
						fmt.Printf("    图片已保存为: %s (数据大小: %d字节)\n", absPath, len(picData))
						// 打印图片数据的前几个字节，用于调试
						if len(picData) > 16 {
							fmt.Printf("    图片数据头: %X\n", picData[:16])
						}
					} else {
						fmt.Printf("    保存图片失败: %v\n", err)
					}
				} else {
					fmt.Println("    警告: 无法获取有效的图片数据!")
				}
			} else if acsAlarmInfo.DwPicDataLen > 0 && acsAlarmInfo.PPicData == nil {
				fmt.Println("    图片数据指针为空，但DwPicDataLen > 0")

				// 处理URL方式(ByPicTransType=1)的图片传输
				if acsAlarmInfo.ByPicTransType == 1 {
					fmt.Println("    图片为URL传输方式，但指针为空，无法获取URL")
				} else if acsAlarmInfo.DwPicDataLen == 1 {
					fmt.Println("    检测到DwPicDataLen=1，这表示有图片但需要通过其他方式获取")
					fmt.Println("    根据用户要求，不使用查询事件日志的方式获取图片")
				}
			} else {
				if acsAlarmInfo.DwPicDataLen == 0 {
					fmt.Println("    不包含图片数据 (DwPicDataLen=0)")
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
