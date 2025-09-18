package sdk

import "unsafe"

// 报警命令常量 (请根据SDK文档核对具体值)
const (
	COMM_ALARM_V30          = 0x4000 // V30报警信息
	COMM_ALARM_ACS          = 0x5002 // 门禁主机报警信息
	COMM_ID_INFO_ALARM      = 0x5200 // 门禁身份证刷卡信息
	COMM_PASSNUM_INFO_ALARM = 0x5201 // 门禁通行人数信息
	// 添加其他需要的报警命令常量...
)

// NET_DVR_TIME 时间结构体 (如果sdk/hcnetsdk_impl.go中未定义)
type NET_DVR_TIME struct {
	DwYear   uint32
	DwMonth  uint32
	DwDay    uint32
	DwHour   uint32
	DwMinute uint32
	DwSecond uint32
}

// NET_DVR_ALARMINFO_V30 报警信息结构体 (V30)
// 请根据 SDK 文档填充完整字段
type NET_DVR_ALARMINFO_V30 struct {
	DwSize               uint32
	DwAlarmType          uint32                 // 报警类型
	DwAlarmInputNumber   uint32                 // 报警输入口，触发报警的输入口号
	ByAlarmOutputNumber  [MAX_ALARMOUT_V30]byte // 触发报警的输出口
	ByAlarmRelateChannel [MAX_CHANNUM_V30]byte  // 触发报警的通道
	ByChannel            [MAX_CHANNUM_V30]byte  // dwAlarmType为0时，表示触发报警的模拟通道
	ByDiskNumber         [MAX_DISKNUM_V30]byte  // dwAlarmType为2时，表示触发报警的硬盘号
	// ... 其他字段，根据SDK文档添加
}

// NET_DVR_ACS_ALARM_INFO 门禁主机报警信息
// 请根据 SDK 文档填充完整字段
type NET_DVR_ACS_ALARM_INFO struct {
	DwSize              uint32
	DwMajor             uint32                 // 报警主类型
	DwMinor             uint32                 // 报警次类型
	StruTime            NET_DVR_TIME           // 报警时间
	SNetUser            [MAX_NAMELEN]byte      // 网络操作的用户名
	StruRemoteHostAddr  NET_DVR_IPADDR         // 远程主机地址
	StruAcsEventInfo    NET_DVR_ACS_EVENT_INFO // 详细参数
	DwPicDataLen        uint32                 // 图片数据长度
	PPicData            *byte                  // 图片数据指针
	WInductiveEventType uint16                 // 归纳事件类型
	ByPicTransType      byte                   // 图片数据传输方式: 0-binary, 1-url
	ByRes1              byte                   // 保留
	DwIOTChannelNo      uint32                 // IOT通道号

	PAcsEventInfoExtend     *byte
	ByAcsEventInfoExtend    byte
	ByTimeType              byte
	ByRes2                  byte
	ByAcsEventInfoExtendV20 byte
	PAcsEventInfoExtendV20  *byte
	ByRes                   [4]byte

	//PChIOTChannelInfo   *byte    // IOT通道信息指针
	//SzIOTChannelInfoLen uint32    // IOT通道信息长度
	//ByRes               [20]byte  // 保留
}

// NET_DVR_ACS_EVENT_INFO 门禁事件信息 (被 NET_DVR_ACS_ALARM_INFO 包含)
// 请根据 SDK 文档填充完整字段
type NET_DVR_ACS_EVENT_INFO struct {
	DwSize              uint32
	ByCardNo            [ACS_CARD_NO_LEN]byte // 卡号
	ByCardType          byte                  // 卡类型
	ByWhiteListNo       byte                  // 白名单号
	ByReportChannel     byte                  // 报告上传通道
	ByCardReaderKind    byte                  // 读卡器类型
	DwCardReaderNo      uint32                // 读卡器编号
	DwDoorNo            uint32                // 门编号
	DwVerifyNo          uint32                // 多重卡认证序号
	DwAlarmInNo         uint32                // 报警输入号
	DwAlarmOutNo        uint32                // 报警输出号
	DwCaseSensorNo      uint32                // 事件触发器编号
	DwRs485No           uint32                // RS485通道号
	DwMultiCardGroupNo  uint32                // 群组编号
	WAccessChannel      uint16                // 通道号
	ByDeviceNo          byte                  // 设备编号
	ByDistractControlNo byte                  // 分控器编号
	DwEmployeeNo        uint32                // 工号
	WLocalControllerID  uint16                // 就地控制器编号
	ByInternetAccess    byte                  // 网口ID
	ByType              byte                  // 防区类型
	// ... 其他字段
}

// NET_DVR_IPADDR IP地址结构 (被 NET_DVR_ACS_ALARM_INFO 包含)
type NET_DVR_IPADDR struct {
	SIpV4 [16]byte
	ByRes [128]byte
}

// 定义常量 (根据需要从C头文件或文档中获取)
const (
	MAX_ALARMOUT_V30 = 4
	MAX_CHANNUM_V30  = 16
	MAX_DISKNUM_V30  = 16
	MAX_NAMELEN      = 16
	ACS_CARD_NO_LEN  = 32
)

// // // 其他需要的结构体和常量定义...

// // 注意: 上述结构体定义是示例性的，你需要根据你的海康SDK文档精确定义字段和类型。
// // 特别是字节数组的大小 (如 MAX_ALARMOUT_V30 等) 和具体的数据类型。

// 查找图片需要的常量
const (
	NET_DVR_FILE_SUCCESS   = 1000 // 获取文件成功
	NET_DVR_FILE_NOFIND    = 1001 // 未查找到文件
	NET_DVR_ISFINDING      = 1002 // 正在查找请等待
	NET_DVR_NOMOREFILE     = 1003 // 没有更多的文件
	NET_DVR_FILE_EXCEPTION = 1004 // 查找文件时异常
)

// NET_DVR_FIND_PICTURE_PARAM 查找图片参数结构体
type NET_DVR_FIND_PICTURE_PARAM struct {
	DwSize        uint32                // 结构体大小
	LChannel      int32                 // 通道号
	ByFileType    byte                  // 图片查找类型
	ByNeedCard    byte                  // 是否需要卡号
	ByProvince    byte                  // 省份索引值
	ByEventType   byte                  // 事件类型
	SCardNum      [CARDNUM_LEN_V30]byte // 卡号
	StruStartTime NET_DVR_TIME          // 查找图片的开始时间
	StruStopTime  NET_DVR_TIME          // 查找图片的结束时间
	ByRes         [40]byte              // 保留字节
}

// NET_DVR_FIND_PICTURE 查找图片结果结构体
type NET_DVR_FIND_PICTURE struct {
	SFileName           [PICTURE_NAME_LEN]byte // 图片文件名
	StruTime            NET_DVR_TIME           // 图片的时间
	DwFileSize          uint32                 // 文件大小
	SCardNum            [CARDNUM_LEN_V30]byte  // 卡号
	ByPlateColor        byte                   // 车牌颜色
	ByVehicleType       byte                   // 车辆类型
	ByFileType          byte                   // 文件类型
	ByRecogResult       byte                   // 识别结果
	SLicense            [MAX_LICENSE_LEN]byte  // 车牌号码
	ByEventSearchStatus byte                   // 连续图片表示同一个事件的图片是否全部查找完成
	ByRes               [75]byte               // 保留字节
}

// NET_DVR_GETPIC_PARAM 获取图片参数结构体
type NET_DVR_GETPIC_PARAM struct {
	DwSize       uint32   // 结构体大小
	DwSignalType uint32   // 图片类型，0-无效，1-JPG，2-BMP
	ByPictype    byte     // 图片传输方式，0-二进制传输，1-文件传输
	ByRes1       [3]byte  // 保留
	DwPicSize    uint32   // 图片大小
	ByRes2       [32]byte // 保留
	PicName      *byte    // 图片名称
}

// 图片查找相关常量
const (
	CARDNUM_LEN_V30  = 40 // 卡号长度
	PICTURE_NAME_LEN = 64 // 图片名称长度
	MAX_LICENSE_LEN  = 16 // 车牌号长度
)

// NET_VCA_POINT 坐标点结构体
type NET_VCA_POINT struct {
	FX float32 // X坐标, 0.000~1
	FY float32 // Y坐标, 0.000~1
}

// 定义常量
const (
	//NET_SDK_EMPLOYEE_NO_LEN = 32 // 员工号长度
	NET_SDK_UUID_LEN = 36 // UUID长度
	NET_DEV_NAME_LEN = 32 // 设备名称长度
)

// NET_DVR_ACS_EVENT_INFO_EXTEND 门禁事件扩展信息
type NET_DVR_ACS_EVENT_INFO_EXTEND struct {
	DwFrontSerialNo       uint32                        // 事件流水号
	ByUserType            byte                          // 人员类型：0-无效，1-普通人（主人），2-来宾（访客），3-禁止名单人，4-管理员
	ByCurrentVerifyMode   byte                          // 读卡器当前验证方式
	ByCurrentEvent        byte                          // 是否为实时事件：0-无效，1-是（实时事件），2-否（离线事件）
	ByPurePwdVerifyEnable byte                          // 设备是否支持纯密码验证：0-不支持，1-支持
	ByEmployeeNo          [NET_SDK_EMPLOYEE_NO_LEN]byte // 工号，人员ID
	ByAttendanceStatus    byte                          // 考勤状态
	ByStatusValue         byte                          // 考勤状态值
	ByRes2                [2]byte                       // 保留
	ByUUID                [NET_SDK_UUID_LEN]byte        // UUID
	ByDeviceName          [NET_DEV_NAME_LEN]byte        // 设备序列号
	ByRes                 [24]byte                      // 保留
}

// NET_DVR_ACS_EVENT_INFO_EXTEND_V20 门禁事件扩展信息V20
type NET_DVR_ACS_EVENT_INFO_EXTEND_V20 struct {
	ByRemoteCheck          byte          // 是否需要远程核验
	ByThermometryUnit      byte          // 测温单位
	ByIsAbnomalTemperature byte          // 人脸抓拍测温是否温度异常
	ByRes2                 byte          // 保留
	FCurrTemperature       float32       // 人脸温度
	StruRegionCoordinates  NET_VCA_POINT // 人脸温度坐标
	DwQRCodeInfoLen        uint32        // 二维码信息长度
	DwVisibleLightDataLen  uint32        // 热成像相机可见光图片长度
	DwThermalDataLen       uint32        // 热成像图片长度
	PQRCodeInfo            *byte         // 二维码信息指针
	PVisibleLightData      *byte         // 热成像相机可见光图片指针
	PThermalData           *byte         // 热成像图片指针
	ByAttendanceLabel      [64]byte      // 考勤自定义标签
	WXCoordinate           uint16        // x坐标
	WYCoordinate           uint16        // y坐标
	WWidth                 uint16        // 人脸宽度
	WHeight                uint16        // 人脸高度
	ByHealthCode           byte          // 健康码状态
	ByNADCode              byte          // 核酸检测
	ByTravelCode           byte          // 行程码
	ByVaccineStatus        byte          // 疫苗状态
	ByRes                  [948]byte     // 保留
}

// NET_DVR_TIME_V30
type NET_DVR_TIME_V30 struct {
	DwYear   uint32
	DwMonth  uint32
	DwDay    uint32
	DwHour   uint32
	DwMinute uint32
	DwSecond uint32
}

// NET_DVR_ID_CARD_INFO
type NET_DVR_ID_CARD_INFO struct {
	ByName  [128]byte
	ByIDNum [32]byte
	ByRes   [16]byte
}

// NET_DVR_ID_CARD_INFO_EXTEND
type NET_DVR_ID_CARD_INFO_EXTEND struct {
	ByRemoteCheck          byte
	ByThermometryUnit      byte
	ByIsAbnomalTemperature byte
	ByRes2                 byte
	FCurrTemperature       float32
	StruRegionCoordinates  NET_VCA_POINT
	DwQRCodeInfoLen        uint32
	DwVisibleLightDataLen  uint32
	DwThermalDataLen       uint32
	PQRCodeInfo            unsafe.Pointer
	PVisibleLightData      unsafe.Pointer
	PThermalData           unsafe.Pointer
	WXCoordinate           uint16
	WYCoordinate           uint16
	WWidth                 uint16
	WHeight                uint16
	ByHealthCode           byte
	ByNADCode              byte
	ByTravelCode           byte
	ByVaccineStatus        byte
	ByRes                  [1012]byte
}

// NET_DVR_ID_CARD_INFO_ALARM
type NET_DVR_ID_CARD_INFO_ALARM struct {
	DwSize                  uint32
	StruIDCardCfg           NET_DVR_ID_CARD_INFO
	DwMajor                 uint32
	DwMinor                 uint32
	StruSwipeTime           NET_DVR_TIME_V30
	ByNetUser               [MAX_NAMELEN]byte
	StruRemoteHostAddr      NET_DVR_IPADDR
	DwCardReaderNo          uint32
	DwDoorNo                uint32
	DwPicDataLen            uint32
	PPicData                unsafe.Pointer
	ByCardType              byte
	ByDeviceNo              byte
	ByMask                  byte
	ByRes2                  byte
	DwFingerPrintDataLen    uint32
	PFingerPrintData        unsafe.Pointer
	DwCapturePicDataLen     uint32
	PCapturePicData         unsafe.Pointer
	DwCertificatePicDataLen uint32
	PCertificatePicData     unsafe.Pointer
	ByCardReaderKind        byte
	ByRes3                  [2]byte
	ByIDCardInfoExtend      byte
	PIDCardInfoExtend       unsafe.Pointer
	DwSerialNo              uint32
	ByRes                   [168]byte
}
