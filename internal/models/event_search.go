package models

import (
	"fmt"
	"time"
	"unsafe"
	"github.com/hikacsuser-go/internal/sdk"
)

// EventSearch 事件查询
type EventSearch struct {
	SDK sdk.HCNetSDK
}

// NewEventSearch 创建事件查询实例
func NewEventSearch(sdk sdk.HCNetSDK) *EventSearch {
	return &EventSearch{
		SDK: sdk,
	}
}

// SearchAllEvent 查询所有事件
func (es *EventSearch) SearchAllEvent(lUserID int) error {
	if lUserID < 0 {
		return fmt.Errorf("无效的用户ID")
	}

	// 创建查询条件结构体
	// 定义NET_DVR_TIME结构体，与Java保持一致
	type NET_DVR_TIME struct {
		DwYear   uint32
		DwMonth  uint32
		DwDay    uint32
		DwHour   uint32
		DwMinute uint32
		DwSecond uint32
	}

	// 创建查询条件结构体，与Java保持一致
	type NET_DVR_ACS_EVENT_COND struct {
		DwSize              uint32
		DwMajor             uint32 // 主类型，0表示全部
		DwMinor             uint32 // 次类型，0表示全部
		StructStartTime     NET_DVR_TIME // 开始时间
		StructEndTime       NET_DVR_TIME // 结束时间
		ByCardNo            [sdk.ACS_CARD_NO_LEN]byte // 卡号
		ByName              [sdk.NAME_LEN]byte        // 持卡人姓名
		ByPicEnable         byte   // 是否带图片，0-不带图片，1-带图片
		ByTimeType          byte   // 时间类型：0-设备本地时间（默认），1-UTC时间
		ByRes2              [2]byte // 保留
		DwBeginSerialNo     uint32 // 起始流水号（为0时默认全部）
		DwEndSerialNo       uint32 // 结束流水号（为0时默认全部）
		DwIOTChannelNo      uint32 // IOT通道号，0-无效
		WInductiveEventType uint16 // 归纳事件类型，0-无效，其他值参见2.2章节
		BySearchType        byte   // 搜索方式：0-保留，1-按事件源搜索
		ByEventAttribute    byte   // 事件属性：0-未定义，1-合法事件，2-其它
		SzMonitorID         [sdk.NET_SDK_MONITOR_ID_LEN]byte // 布防点ID
		ByEmployeeNo        [sdk.NET_SDK_EMPLOYEE_NO_LEN]byte // 工号
		ByRes               [140]byte // 保留
	}

	// 创建并初始化查询条件
	var struAcsEventCond NET_DVR_ACS_EVENT_COND
	
	// 设置结构体大小
	struAcsEventCond.DwSize = uint32(unsafe.Sizeof(struAcsEventCond))
	
	// 设置查询所有主次类型的报警
	struAcsEventCond.DwMajor = 0 // 查询所有主类型事件
	struAcsEventCond.DwMinor = 0 // 查询所有次类型事件
	
	// 设置开始时间 - 与Java代码保持一致，使用固定时间而不是动态计算
	struAcsEventCond.StructStartTime.DwYear = 2024
	struAcsEventCond.StructStartTime.DwMonth = 8
	struAcsEventCond.StructStartTime.DwDay = 1
	struAcsEventCond.StructStartTime.DwHour = 0
	struAcsEventCond.StructStartTime.DwMinute = 0
	struAcsEventCond.StructStartTime.DwSecond = 0
	
	// 设置结束时间 - 与Java代码保持一致
	struAcsEventCond.StructEndTime.DwYear = 2024
	struAcsEventCond.StructEndTime.DwMonth = 8
	struAcsEventCond.StructEndTime.DwDay = 9
	struAcsEventCond.StructEndTime.DwHour = 23
	struAcsEventCond.StructEndTime.DwMinute = 59
	struAcsEventCond.StructEndTime.DwSecond = 59
	
	// 设置其他参数 - 与Java代码保持完全一致
	struAcsEventCond.WInductiveEventType = 1 // 归纳事件类型
	struAcsEventCond.ByPicEnable = 1        // 带图片

	// 启动远程配置
	lHandle := es.SDK.NET_DVR_StartRemoteConfig(lUserID, sdk.NET_DVR_GET_ACS_EVENT, unsafe.Pointer(&struAcsEventCond), uint32(unsafe.Sizeof(struAcsEventCond)), 0, nil)
	if lHandle < 0 {
		return fmt.Errorf("NET_DVR_StartRemoteConfig失败，错误码: %d", es.SDK.NET_DVR_GetLastError())
	}
	defer es.SDK.NET_DVR_StopRemoteConfig(lHandle)

	// 创建接收事件的结构体 - 与Java版本保持一致
	type NET_DVR_ACS_EVENT_CFG struct {
		DwSize         uint32
		DwMajor        uint32 // 主类型
		DwMinor        uint32 // 次类型
		StructTime     NET_DVR_TIME // 时间
		StructAcsEventInfo struct { // 门禁事件信息
			ByCardNo      [sdk.MAX_CARDNO_LEN]byte // 卡号
			ByEmployeeNo  [sdk.NET_SDK_EMPLOYEE_NO_LEN]byte // 工号
			DwEmployeeNo  uint32 // 工号（数值）
			ByRes         [108]byte
		}
		DwPicDataLen   uint32 // 图片数据长度
		PPicData       unsafe.Pointer // 图片数据
		ByRes          [40]byte
	}

	// 初始化接收事件的结构体
	var struAcsEventCfg NET_DVR_ACS_EVENT_CFG
	struAcsEventCfg.DwSize = uint32(unsafe.Sizeof(struAcsEventCfg))

	// 循环获取事件
	var i int = 0
	for {
		var resultLen uint32
		// 获取下一个事件
		dwEventSearch := es.SDK.NET_DVR_GetNextRemoteConfig(lHandle, unsafe.Pointer(&struAcsEventCfg), uint32(unsafe.Sizeof(struAcsEventCfg)),&resultLen)
		if dwEventSearch <= -1 {
			fmt.Printf("NET_DVR_GetNextRemoteConfig接口调用失败，错误码：%d\n", es.SDK.NET_DVR_GetLastError())
			break
		}

		// 处理返回结果
		if dwEventSearch == sdk.NET_SDK_GET_NEXT_STATUS_NEED_WAIT {
			fmt.Println("配置等待....")
			time.Sleep(10 * time.Millisecond)
			continue
		} else if dwEventSearch == sdk.NET_SDK_NEXT_STATUS__FINISH {
			fmt.Println("获取事件完成")
			break
		} else if dwEventSearch == sdk.NET_SDK_GET_NEXT_STATUS_FAILED {
			fmt.Println("获取事件出现异常")
			break
		} else if dwEventSearch == sdk.NET_SDK_GET_NEXT_STATUS_SUCCESS {
			// 获取事件成功，处理事件信息
			cardNo := string(struAcsEventCfg.StructAcsEventInfo.ByCardNo[:])
			fmt.Printf("%d获取事件成功, 报警主类型：%x 报警次类型：%x 卡号：%s\n", 
				i, struAcsEventCfg.DwMajor, struAcsEventCfg.DwMinor, cardNo)
			
			// 打印刷卡时间
			fmt.Printf("刷卡时间：年：%d 月：%d 日：%d 时：%d 分：%d 秒：%d\n",
				struAcsEventCfg.StructTime.DwYear, struAcsEventCfg.StructTime.DwMonth, struAcsEventCfg.StructTime.DwDay,
				struAcsEventCfg.StructTime.DwHour, struAcsEventCfg.StructTime.DwMinute, struAcsEventCfg.StructTime.DwSecond)
			
			// 处理图片数据（如果有）
			if struAcsEventCfg.DwPicDataLen > 0 && struAcsEventCfg.PPicData != nil {
				// 这里可以添加保存图片的代码
				fmt.Printf("事件包含图片数据，长度：%d\n", struAcsEventCfg.DwPicDataLen)
			}
			
			i++
			continue
		}
	}

	fmt.Printf("查询门禁事件完成，共获取%d条事件\n", i)
	return nil
}

// sendISAPIRequest 发送ISAPI请求
func (es *EventSearch) sendISAPIRequest(lUserID int, url string, requestData []byte) ([]byte, error) {
	// 创建ISAPI输入参数
	inputParam := make([]byte, len(url)+1)
	copy(inputParam, []byte(url))

	var outputBuf [10240]byte

	// 调用远程配置
	lHandle := es.SDK.NET_DVR_StartRemoteConfig(lUserID, sdk.COMM_ISAPI_CONFIG, unsafe.Pointer(&inputParam[0]), uint32(len(inputParam)), 0, nil)
	if lHandle < 0 {
		return nil, fmt.Errorf("NET_DVR_StartRemoteConfig失败，错误码: %d", es.SDK.NET_DVR_GetLastError())
	}
	defer es.SDK.NET_DVR_StopRemoteConfig(lHandle)

	// 发送数据
	var resultLen uint32
	var result int

	for {
		if requestData == nil || len(requestData) == 0 {
			// 如果没有请求数据，只发送空数据
			emptyData := []byte{}
			result = es.SDK.NET_DVR_SendWithRecvRemoteConfig(lHandle, unsafe.Pointer(&emptyData[0]), 0, unsafe.Pointer(&outputBuf[0]), uint32(len(outputBuf)), &resultLen)
		} else {
			// 发送请求数据
			result = es.SDK.NET_DVR_SendWithRecvRemoteConfig(lHandle, unsafe.Pointer(&requestData[0]), uint32(len(requestData)), unsafe.Pointer(&outputBuf[0]), uint32(len(outputBuf)), &resultLen)
		}

		// 处理配置结果
		if result == -1 {
			return nil, fmt.Errorf("发送ISAPI请求失败，错误码: %d", es.SDK.NET_DVR_GetLastError())
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
