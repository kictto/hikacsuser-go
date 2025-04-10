package models

import (
	"fmt"
	"github.com/clockworkchen/hikacsuser-go/internal/sdk"
	"time"
	"unsafe"
)

// NET_DVR_ACS_CFG 门禁主机参数
type NET_DVR_ACS_CFG struct {
	DwSize             uint32    // 结构体大小
	ByRS485Backup      byte      // 是否启用下行RS485通信备份功能，0-不启用，1-启用
	ByShowCapPic       byte      // 是否显示抓拍图片 0-不显示，1-显示
	ByShowCardNo       byte      // 是否显示卡号 0-不显示，1-显示
	ByShowUserInfo     byte      // 是否显示用户信息，0-不显示，1-显示
	ByOverlayUserInfo  byte      // 是否叠加用户信息，0-不叠加，1-叠加
	ByVoicePrompt      byte      // 是否语音提示 0-不提示，1-提示
	ByUploadCapPic     byte      // 联动抓拍是否上传图片 0-不上传，1-上传
	BySaveCapPic       byte      // 是否保存抓拍图片，0-不保存，1-保存
	ByInputCardNo      byte      // 是否是否允许按键输入卡号，0-不允许，1-允许
	ByEnableWifiDetect byte      // 是否启动wifi，0-不启动，1-启动
	ByEnable3G4G       byte      // 3G4G使能，0-不使能，1-使能
	ByProtocol         byte      // 读卡器通信协议类型，0-私有协议（默认），1-OSDP协议
	ByRes              [500]byte // 保留字节
}

// NET_DVR_ACS_WORK_STATUS_V50 门禁主机工作状态V50
type NET_DVR_ACS_WORK_STATUS_V50 struct {
	DwSize                          uint32                               // 结构体大小
	ByDoorLockStatus                [sdk.MAX_DOOR_NUM_256]byte           // 门锁状态(继电器开合状态)，0-正常关，1-正常开，2-短路报警，3-断路报警，4-异常报警
	ByDoorStatus                    [sdk.MAX_DOOR_NUM_256]byte           // 门状态(楼层状态)，1-休眠，2-常开状态(自由)，3-常闭状态(禁用)，4-普通状态(受控)
	ByMagneticStatus                [sdk.MAX_DOOR_NUM_256]byte           // 门磁状态，0-正常关，1-正常开，2-短路报警，3-断路报警，4-异常报警
	ByCaseStatus                    [sdk.MAX_CASE_SENSOR_NUM]byte        // 事件触发器状态，0-无输入，1-有输入
	WBatteryVoltage                 uint16                               // 蓄电池电压值，实际值乘10，单位：伏特
	ByBatteryLowVoltage             byte                                 // 蓄电池是否处于低压状态，0-否，1-是
	ByPowerSupplyStatus             byte                                 // 设备供电状态，1-交流电供电，2-蓄电池供电
	ByMultiDoorInterlockStatus      byte                                 // 多门互锁状态，0-关闭，1-开启
	ByAntiSneakStatus               byte                                 // 反潜回状态，0-关闭，1-开启
	ByHostAntiDismantleStatus       byte                                 // 主机防拆状态，0-关闭，1-开启
	ByIndicatorLightStatus          byte                                 // 指示灯状态，0-掉线，1-在线
	ByCardReaderOnlineStatus        [sdk.MAX_CARD_READER_NUM_512]byte    // 读卡器在线状态，0-不在线，1-在线
	ByCardReaderAntiDismantleStatus [sdk.MAX_CARD_READER_NUM_512]byte    // 读卡器防拆状态，0-关闭，1-开启
	ByCardReaderVerifyMode          [sdk.MAX_CARD_READER_NUM_512]byte    // 读卡器当前验证方式，1-休眠，2-刷卡+密码，3-刷卡，4-刷卡或密码
	BySetupAlarmStatus              [sdk.MAX_ALARMHOST_ALARMIN_NUM]byte  // 报警输入口布防状态，0-对应报警输入口处于撤防状态，1-对应报警输入口处于布防状态
	ByAlarmInStatus                 [sdk.MAX_ALARMHOST_ALARMIN_NUM]byte  // 按位表示报警输入口报警状态，0-对应报警输入口当前无报警，1-对应报警输入口当前有报警
	ByAlarmOutStatus                [sdk.MAX_ALARMHOST_ALARMOUT_NUM]byte // 按位表示报警输出口状态，0-对应报警输出口无报警，1-对应报警输出口有报警
	DwCardNum                       uint32                               // 已添加的卡数量
	ByFireAlarmStatus               byte                                 // 消防报警状态显示：0-正常、1-短路报警、2-断开报警
	ByBatteryChargeStatus           byte                                 // 电池充电状态：0-无效；1-充电中；2-未充电
	ByMasterChannelControllerStatus byte                                 // 主通道控制器在线状态：0-无效；1-不在线；2-在线
	BySlaveChannelControllerStatus  byte                                 // 从通道控制器在线状态：0-无效；1-不在线；2-在线
	ByAntiSneakServerStatus         byte                                 // 反潜回服务器状态：0-无效，1-未启用，2-正常，3-断开
	ByRes3                          [3]byte                              // 保留
	DwAllowFaceNum                  uint32                               // 已添加的允许名单人脸数量
	DwBlockFaceNum                  uint32                               // 已添加的禁止名单人脸数量
	ByRes2                          [108]byte                            // 保留字节
}

// ACSManage 门禁管理
type ACSManage struct {
	SDK sdk.HCNetSDK
}

// NewACSManage 创建门禁管理实例
func NewACSManage(sdk sdk.HCNetSDK) *ACSManage {
	return &ACSManage{
		SDK: sdk,
	}
}

// AcsCfg 获取门禁参数
func (am *ACSManage) AcsCfg(lUserID int) error {
	if lUserID < 0 {
		return fmt.Errorf("无效的用户ID")
	}

	// 创建门禁参数结构体并初始化大小
	var acsCfg NET_DVR_ACS_CFG
	acsCfg.DwSize = uint32(unsafe.Sizeof(acsCfg))

	// 创建接收数据长度变量
	var bytesReturned uint32

	// 尝试获取门禁参数
	if !am.SDK.NET_DVR_GetDVRConfig(lUserID, sdk.NET_DVR_GET_ACS_CFG, 0xFFFFFFFF, unsafe.Pointer(&acsCfg), uint32(unsafe.Sizeof(acsCfg)), &bytesReturned) {
		errCode := am.SDK.NET_DVR_GetLastError()
		return fmt.Errorf("获取门禁参数失败，错误码: %d", errCode)
	}

	// 打印门禁参数信息
	fmt.Printf("获取门禁参数成功\n")
	fmt.Printf("1.是否启用下行RS485通信备份功能：%d\n", acsCfg.ByRS485Backup)
	fmt.Printf("2.是否显示抓拍图片：%d\n", acsCfg.ByShowCapPic)
	fmt.Printf("3.是否显示卡号：%d\n", acsCfg.ByShowCardNo)
	fmt.Printf("4.是否显示用户信息：%d\n", acsCfg.ByShowUserInfo)
	fmt.Printf("5.是否叠加用户信息：%d\n", acsCfg.ByOverlayUserInfo)
	fmt.Printf("6.是否开启语音提示：%d\n", acsCfg.ByVoicePrompt)
	fmt.Printf("7.联动抓图是否上传：%d\n", acsCfg.ByUploadCapPic)
	fmt.Printf("8.是否保存抓拍图片：%d\n", acsCfg.BySaveCapPic)
	fmt.Printf("9.是否允许按键输入卡号：%d\n", acsCfg.ByInputCardNo)
	fmt.Printf("10.是否启动wifi：%d\n", acsCfg.ByEnableWifiDetect)
	fmt.Printf("11.3G4G使能：%d\n", acsCfg.ByEnable3G4G)
	fmt.Printf("12.读卡器通信协议类型：%d\n", acsCfg.ByProtocol)

	// 设置门禁参数
	acsCfg.ByShowCardNo = 1   // 开启显示卡号
	acsCfg.ByVoicePrompt = 0  // 关闭语音提示
	acsCfg.ByUploadCapPic = 1 // 开启联动抓图上传
	acsCfg.ByShowCapPic = 1   // 开启显示抓拍图片

	// 保存门禁参数
	if !am.SDK.NET_DVR_SetDVRConfig(lUserID, sdk.NET_DVR_SET_ACS_CFG, 0xFFFFFFFF, unsafe.Pointer(&acsCfg), uint32(unsafe.Sizeof(acsCfg))) {
		errCode := am.SDK.NET_DVR_GetLastError()
		return fmt.Errorf("设置门禁参数失败，错误码: %d", errCode)
	}

	fmt.Println("设置门禁参数成功！！！")
	return nil
}

// GetAcsStatus 获取门禁状态
func (am *ACSManage) GetAcsStatus(lUserID int) error {
	if lUserID < 0 {
		return fmt.Errorf("无效的用户ID")
	}

	// 创建门禁状态结构体并初始化大小
	var acsStatus NET_DVR_ACS_WORK_STATUS_V50
	// 确保结构体所有字段都被初始化为0
	acsStatus = NET_DVR_ACS_WORK_STATUS_V50{}
	// 设置结构体大小
	acsStatus.DwSize = uint32(unsafe.Sizeof(acsStatus))

	// 创建接收数据长度变量
	var bytesReturned uint32

	// 尝试获取门禁状态
	fmt.Printf("结构体大小: %d 字节\n", acsStatus.DwSize)
	fmt.Printf("调用参数: lUserID=%d, command=%d, channel=0xFFFFFFFF\n", lUserID, sdk.NET_DVR_GET_ACS_WORK_STATUS_V50)

	// 打印结构体内存地址，用于调试
	fmt.Printf("结构体内存地址: %p\n", &acsStatus)

	// 使用指针传递结构体
	if !am.SDK.NET_DVR_GetDVRConfig(lUserID, sdk.NET_DVR_GET_ACS_WORK_STATUS_V50, 0xFFFFFFFF, unsafe.Pointer(&acsStatus), uint32(unsafe.Sizeof(acsStatus)), &bytesReturned) {
		errCode := am.SDK.NET_DVR_GetLastError()
		var errNo int32
		errMsg := am.SDK.NET_DVR_GetErrorMsg(&errNo)
		return fmt.Errorf("获取门禁状态失败，错误码: %d, 错误信息: %s", errCode, errMsg)
	}

	// 打印返回的数据长度
	fmt.Printf("获取门禁状态成功，返回数据长度: %d 字节\n", bytesReturned)

	// 打印门禁状态信息，参照Java版本
	fmt.Printf("获取门禁主机工作状态成功！！！\n")
	fmt.Printf("1.门锁状态（或者梯控的继电器开合状态）：%d (0-正常关，1-正常开，2-短路报警，3-断路报警，4-异常报警)\n", acsStatus.ByDoorLockStatus[0])
	fmt.Printf("2.门状态（或者梯控的楼层状态）：%d (1-休眠，2-常开状态(自由)，3-常闭状态(禁用)，4-普通状态(受控))\n", acsStatus.ByDoorStatus[0])
	fmt.Printf("3.门磁状态：%d (0-正常关，1-正常开，2-短路报警，3-断路报警，4-异常报警)\n", acsStatus.ByMagneticStatus[0])
	fmt.Printf("4.事件报警输入状态：%d (0-无输入，1-有输入)\n", acsStatus.ByCaseStatus[0])
	fmt.Printf("5.蓄电池电压值：%.1f V\n", float64(acsStatus.WBatteryVoltage)/10.0)
	fmt.Printf("6.蓄电池是否处于低压状态：%d (0-否，1-是)\n", acsStatus.ByBatteryLowVoltage)
	fmt.Printf("7.设备供电状态：%d (1-交流电供电，2-蓄电池供电)\n", acsStatus.ByPowerSupplyStatus)
	fmt.Printf("8.多门互锁状态：%d (0-关闭，1-开启)\n", acsStatus.ByMultiDoorInterlockStatus)
	fmt.Printf("9.反潜回状态：%d (0-关闭，1-开启)\n", acsStatus.ByAntiSneakStatus)
	fmt.Printf("10.主机防拆状态：%d (0-关闭，1-开启)\n", acsStatus.ByHostAntiDismantleStatus)
	fmt.Printf("11.指示灯状态：%d (0-掉线，1-在线)\n", acsStatus.ByIndicatorLightStatus)
	fmt.Printf("12.读卡器在线状态：%d (0-不在线，1-在线)\n", acsStatus.ByCardReaderOnlineStatus[0])
	fmt.Printf("13.读卡器防拆状态：%d (0-关闭，1-开启)\n", acsStatus.ByCardReaderAntiDismantleStatus[0])
	fmt.Printf("14.读卡器当前验证方式：%d (1-休眠，2-刷卡+密码，3-刷卡，4-刷卡或密码)\n", acsStatus.ByCardReaderVerifyMode[0])
	fmt.Printf("15.已添加的卡数量：%d\n", acsStatus.DwCardNum)
	fmt.Printf("16.消防报警状态：%d (0-正常、1-短路报警、2-断开报警)\n", acsStatus.ByFireAlarmStatus)
	fmt.Printf("17.电池充电状态：%d (0-无效；1-充电中；2-未充电)\n", acsStatus.ByBatteryChargeStatus)
	fmt.Printf("18.主通道控制器在线状态：%d (0-无效；1-不在线；2-在线)\n", acsStatus.ByMasterChannelControllerStatus)
	fmt.Printf("19.从通道控制器在线状态：%d (0-无效；1-不在线；2-在线)\n", acsStatus.BySlaveChannelControllerStatus)
	fmt.Printf("20.反潜回服务器状态：%d (0-无效，1-未启用，2-正常，3-断开)\n", acsStatus.ByAntiSneakServerStatus)
	fmt.Printf("21.已添加的允许名单人脸数量：%d\n", acsStatus.DwAllowFaceNum)
	fmt.Printf("22.已添加的禁止名单人脸数量：%d\n", acsStatus.DwBlockFaceNum)

	return nil
}

// RemoteControlGate 远程控门
func (am *ACSManage) RemoteControlGate(lUserID int) error {
	if lUserID < 0 {
		return fmt.Errorf("无效的用户ID")
	}

	/**
	 * NET_DVR_ControlGateway参数说明:
	 * 第二个参数lGatewayIndex: 门禁序号（楼层编号、锁ID），从1开始，-1表示对所有门（或者梯控的所有楼层）进行操作
	 * 第三个参数dwStaic: 命令值：
	 *   0-关闭（对于梯控，表示受控）
	 *   1-打开（对于梯控，表示开门）
	 *   2-常开（对于梯控，表示自由、通道状态）
	 *   3-常关（对于梯控，表示禁用）
	 *   4-恢复（梯控，普通状态）
	 *   5-访客呼梯（梯控）
	 *   6-住户呼梯（梯控）
	 */

	// 远程控门：1号门，打开状态
	if !am.SDK.NET_DVR_ControlGateway(lUserID, 1, 1) {
		errCode := am.SDK.NET_DVR_GetLastError()
		return fmt.Errorf("远程控门失败，错误码: %d", errCode)
	}

	fmt.Println("远程控门成功")
	return nil
}

// sendISAPIRequest 发送ISAPI请求
func (am *ACSManage) sendISAPIRequest(lUserID int, url string, requestData []byte) ([]byte, error) {
	// 创建ISAPI输入参数
	inputParam := make([]byte, len(url)+1)
	copy(inputParam, []byte(url))

	var outputBuf [10240]byte

	// 调用远程配置
	lHandle := am.SDK.NET_DVR_StartRemoteConfig(lUserID, sdk.COMM_ISAPI_CONFIG, unsafe.Pointer(&inputParam[0]), uint32(len(inputParam)), 0, nil)
	if lHandle < 0 {
		return nil, fmt.Errorf("NET_DVR_StartRemoteConfig失败，错误码: %d", am.SDK.NET_DVR_GetLastError())
	}
	defer am.SDK.NET_DVR_StopRemoteConfig(lHandle)

	// 发送数据
	var resultLen uint32
	var result int

	for {
		if requestData == nil || len(requestData) == 0 {
			// 如果没有请求数据，只发送空数据
			emptyData := []byte{}
			result = am.SDK.NET_DVR_SendWithRecvRemoteConfig(lHandle, unsafe.Pointer(&emptyData[0]), 0, unsafe.Pointer(&outputBuf[0]), uint32(len(outputBuf)), &resultLen)
		} else {
			// 发送请求数据
			result = am.SDK.NET_DVR_SendWithRecvRemoteConfig(lHandle, unsafe.Pointer(&requestData[0]), uint32(len(requestData)), unsafe.Pointer(&outputBuf[0]), uint32(len(outputBuf)), &resultLen)
		}

		// 处理配置结果
		if result == -1 {
			return nil, fmt.Errorf("发送ISAPI请求失败，错误码: %d", am.SDK.NET_DVR_GetLastError())
		} else if result == sdk.NET_SDK_CONFIG_STATUS_NEED_WAIT {
			// 配置等待，继续等待
			time.Sleep(10 * time.Millisecond)
			continue
		} else if result == sdk.NET_SDK_CONFIG_STATUS_FAILED {
			return nil, fmt.Errorf("配置失败")
		} else if result == sdk.NET_SDK_CONFIG_STATUS_EXCEPTION {
			return nil, fmt.Errorf("配置异常")
		} else if result == sdk.NET_SDK_CONFIG_STATUS_SUCCESS {
			// 提取响应数据
			response := make([]byte, resultLen)
			copy(response, outputBuf[:resultLen])
			return response, nil
		} else if result == sdk.NET_SDK_CONFIG_STATUS_FINISH {
			break
		}
	}

	// 提取响应数据
	response := make([]byte, resultLen)
	copy(response, outputBuf[:resultLen])
	return response, nil
}
