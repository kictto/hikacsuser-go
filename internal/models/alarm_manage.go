package models

import (
	"fmt"
	"time"
	"unsafe"

	"github.com/clockworkchen/hikacsuser-go/internal/sdk"
)

// AlarmManage 报警布防管理
type AlarmManage struct {
	SDK         sdk.HCNetSDK
	AlarmHandle int // 报警布防句柄
}

// NewAlarmManage 创建报警布防管理实例
func NewAlarmManage(sdk sdk.HCNetSDK) *AlarmManage {
	return &AlarmManage{
		SDK:         sdk,
		AlarmHandle: -1,
	}
}

// SetupAlarm 报警布防
func (am *AlarmManage) SetupAlarm(lUserID int) error {
	if lUserID < 0 {
		return fmt.Errorf("无效的用户ID: %d", lUserID)
	}

	// 创建并初始化报警布防参数 - 使用 sdk 包中的类型
	var struSetupAlarmParam sdk.NET_DVR_SETUPALARM_PARAM
	struSetupAlarmParam.DwSize = uint32(unsafe.Sizeof(struSetupAlarmParam))
	struSetupAlarmParam.ByLevel = 1              // 布防优先级：中
	struSetupAlarmParam.ByAlarmInfoType = 1      // 使用新报警信息
	struSetupAlarmParam.ByRetAlarmTypeV40 = 0    // 使用COMM_ALARM_V30
	struSetupAlarmParam.ByFaceAlarmDetection = 1 // 启用人脸侦测报警

	// 启动报警布防 - 直接传递指针，移除 unsafe.Pointer 转换
	lHandle := am.SDK.NET_DVR_SetupAlarmChan_V41(lUserID, &struSetupAlarmParam)
	if lHandle < 0 {
		return fmt.Errorf("NET_DVR_SetupAlarmChan_V41失败，错误码: %d", am.SDK.NET_DVR_GetLastError())
	}

	// 保存报警布防句柄
	am.AlarmHandle = lHandle
	fmt.Printf("报警布防成功，句柄: %d\n", lHandle)
	return nil
}

// CloseAlarm 报警撤防
func (am *AlarmManage) CloseAlarm() error {
	// 检查是否已经布防
	if am.AlarmHandle < 0 {
		return fmt.Errorf("未进行报警布防或布防句柄无效")
	}

	// 关闭报警布防
	if !am.SDK.NET_DVR_CloseAlarmChan_V30(am.AlarmHandle) {
		return fmt.Errorf("NET_DVR_CloseAlarmChan_V30失败，错误码: %d", am.SDK.NET_DVR_GetLastError())
	}

	// 重置报警布防句柄
	am.AlarmHandle = -1
	fmt.Println("报警撤防成功")
	return nil
}

// SearchAlarmEvent 查询报警事件
func (am *AlarmManage) SearchAlarmEvent(lUserID int) error {
	if lUserID < 0 {
		return fmt.Errorf("无效的用户ID")
	}

	// 创建查询条件结构体
	// 定义NET_DVR_TIME结构体
	type NET_DVR_TIME struct {
		DwYear   uint32
		DwMonth  uint32
		DwDay    uint32
		DwHour   uint32
		DwMinute uint32
		DwSecond uint32
	}

	// 创建查询条件结构体
	type NET_DVR_ALARM_SEARCH_COND struct {
		DwSize          uint32
		StructStartTime NET_DVR_TIME // 开始时间
		StructEndTime   NET_DVR_TIME // 结束时间
		ByAlarmType     byte         // 报警类型，0-全部
		ByRes           [3]byte      // 保留
		DwMaxResults    uint32       // 最大返回条数
		ByRes1          [128]byte    // 保留
	}

	// 创建并初始化查询条件
	var struAlarmSearchCond NET_DVR_ALARM_SEARCH_COND
	struAlarmSearchCond.DwSize = uint32(unsafe.Sizeof(struAlarmSearchCond))

	// 设置查询所有类型的报警
	struAlarmSearchCond.ByAlarmType = 0 // 查询所有类型报警

	// 设置开始时间 - 7天前
	now := time.Now()
	startTime := now.AddDate(0, 0, -7)
	struAlarmSearchCond.StructStartTime.DwYear = uint32(startTime.Year())
	struAlarmSearchCond.StructStartTime.DwMonth = uint32(startTime.Month())
	struAlarmSearchCond.StructStartTime.DwDay = uint32(startTime.Day())
	struAlarmSearchCond.StructStartTime.DwHour = 0
	struAlarmSearchCond.StructStartTime.DwMinute = 0
	struAlarmSearchCond.StructStartTime.DwSecond = 0

	// 设置结束时间 - 当前时间
	struAlarmSearchCond.StructEndTime.DwYear = uint32(now.Year())
	struAlarmSearchCond.StructEndTime.DwMonth = uint32(now.Month())
	struAlarmSearchCond.StructEndTime.DwDay = uint32(now.Day())
	struAlarmSearchCond.StructEndTime.DwHour = uint32(now.Hour())
	struAlarmSearchCond.StructEndTime.DwMinute = uint32(now.Minute())
	struAlarmSearchCond.StructEndTime.DwSecond = uint32(now.Second())

	// 设置最大返回条数
	struAlarmSearchCond.DwMaxResults = 50

	// 启动远程配置
	lHandle := am.SDK.NET_DVR_StartRemoteConfig(lUserID, sdk.NET_DVR_GET_ACS_EVENT, unsafe.Pointer(&struAlarmSearchCond), uint32(unsafe.Sizeof(struAlarmSearchCond)), 0, nil)
	if lHandle < 0 {
		return fmt.Errorf("NET_DVR_StartRemoteConfig失败，错误码: %d", am.SDK.NET_DVR_GetLastError())
	}
	defer am.SDK.NET_DVR_StopRemoteConfig(lHandle)

	// 创建接收报警事件的结构体
	type NET_DVR_ALARM_EVENT_INFO struct {
		DwSize         uint32
		ByAlarmType    byte         // 报警类型
		ByRes          [3]byte      // 保留
		UnionAlarmInfo [128]byte    // 报警信息联合体，不同报警类型对应不同的结构体
		StructTime     NET_DVR_TIME // 报警时间
		ByRes1         [64]byte     // 保留
	}

	// 初始化接收报警事件的结构体
	var struAlarmEventInfo NET_DVR_ALARM_EVENT_INFO
	struAlarmEventInfo.DwSize = uint32(unsafe.Sizeof(struAlarmEventInfo))

	// 循环获取报警事件
	var i int = 0
	for {
		var resultLen uint32
		// 获取下一个报警事件
		dwAlarmSearch := am.SDK.NET_DVR_GetNextRemoteConfig(lHandle, unsafe.Pointer(&struAlarmEventInfo), uint32(unsafe.Sizeof(struAlarmEventInfo)), &resultLen)
		if dwAlarmSearch <= -1 {
			fmt.Printf("NET_DVR_GetNextRemoteConfig接口调用失败，错误码：%d\n", am.SDK.NET_DVR_GetLastError())
			break
		}

		// 处理返回结果
		if dwAlarmSearch == sdk.NET_SDK_GET_NEXT_STATUS_NEED_WAIT {
			fmt.Println("配置等待....")
			time.Sleep(10 * time.Millisecond)
			continue
		} else if dwAlarmSearch == sdk.NET_SDK_NEXT_STATUS__FINISH {
			fmt.Println("获取报警事件完成")
			break
		} else if dwAlarmSearch == sdk.NET_SDK_GET_NEXT_STATUS_FAILED {
			fmt.Println("获取报警事件出现异常")
			break
		} else if dwAlarmSearch == sdk.NET_SDK_GET_NEXT_STATUS_SUCCESS {
			// 获取报警事件成功，处理事件信息
			fmt.Printf("%d获取报警事件成功, 报警类型：%d\n", i, struAlarmEventInfo.ByAlarmType)

			// 打印报警时间
			fmt.Printf("报警时间：年：%d 月：%d 日：%d 时：%d 分：%d 秒：%d\n",
				struAlarmEventInfo.StructTime.DwYear, struAlarmEventInfo.StructTime.DwMonth, struAlarmEventInfo.StructTime.DwDay,
				struAlarmEventInfo.StructTime.DwHour, struAlarmEventInfo.StructTime.DwMinute, struAlarmEventInfo.StructTime.DwSecond)

			i++
			continue
		}
	}

	fmt.Printf("查询报警事件完成，共获取%d条事件\n", i)
	return nil
}
