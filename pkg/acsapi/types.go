package acsapi

import (
	"time"
)

// UserInfo 用户信息结构体，用于注册用户
type UserInfo struct {
	EmployeeNo      string     // 工号，必填
	Name            string     // 姓名，默认为"测试用户"
	UserType        string     // 用户类型，默认为"normal"
	Valid           ValidInfo  // 有效期信息
	BelongGroup     string     // 所属组，默认为"1"
	DoorRight       string     // 门权限，默认为"1"
	Password        string     // 密码，默认为空
	RightPlan       []DoorPlan // 权限计划
	MaxOpenDoorTime int        // 最大开门次数，默认为0
	OpenDoorTime    int        // 开门时间，默认为0
	LocalUIRight    bool       // 本地UI权限，默认为false
	UserVerifyMode  string     // 用户验证模式，默认为"cardOrFace"
}

// ValidInfo 有效期信息
type ValidInfo struct {
	Enable    bool      // 是否启用，默认为true
	BeginTime time.Time // 开始时间，默认为2021-01-01T00:00:00
	EndTime   time.Time // 结束时间，默认为2031-01-01T00:00:00
	TimeType  string    // 时间类型，默认为"local"
}

// DoorPlan 门权限计划
type DoorPlan struct {
	DoorNo         int    // 门编号，默认为1
	PlanTemplateNo string // 计划模板编号，默认为"1"
}

// CardInfo 卡片信息结构体，用于添加卡片
type CardInfo struct {
	EmployeeNo string // 工号，必填
	CardNo     string // 卡号，必填
	CardType   string // 卡片类型，默认为"normalCard"
}

// FaceInfo 人脸信息结构体，用于添加人脸
type FaceInfo struct {
	EmployeeNo  string // 工号，必填
	FaceLibType string // 人脸库类型，默认为"blackFD"
	FDID        string // 人脸库ID，默认为"1"
	FaceURL     string // 人脸图片URL，用于URL方式添加人脸
	FaceFile    string // 人脸图片文件路径，用于二进制方式添加人脸
	FaceData    []byte // 人脸图片二进制数据，可直接传递图片数据，与FaceFile二选一
}

// LoginInfo 登录信息
type LoginInfo struct {
	DeviceIP   string // 设备IP
	DevicePort uint16 // 设备端口
	Username   string // 用户名
	Password   string // 密码
}

// DeviceInfo 设备信息
type DeviceInfo struct {
	SerialNumber      string // 设备序列号
	AlarmInPortNum    byte   // 报警输入端口数
	AlarmOutPortNum   byte   // 报警输出端口数
	DiskNum           byte   // 硬盘数
	DeviceType        byte   // 设备类型
	ChannelNum        byte   // 通道数
	StartChannel      byte   // 起始通道
	IPChannelNum      byte   // IP通道数
	AudioChannelNum   byte   // 音频通道数
	IPAlarmInPortNum  byte   // IP报警输入端口数
	IPAlarmOutPortNum byte   // IP报警输出端口数
	ZeroChannelNum    byte   // 零通道数
	PasswordLevel     byte   // 密码安全等级
	OEMCode           uint32 // OEM代码
}

// Time 时间结构体
type Time struct {
	Year   uint32 // 年
	Month  uint32 // 月
	Day    uint32 // 日
	Hour   uint32 // 时
	Minute uint32 // 分
	Second uint32 // 秒
}

// IPAddr IP地址结构体
type IPAddr struct {
	IPV4 string // IPv4地址
}

// ACSEventInfo 门禁事件信息
type ACSEventInfo struct {
	CardNo            string // 卡号
	CardType          byte   // 卡类型
	WhiteListNo       byte   // 白名单号
	ReportChannel     byte   // 报告上传通道
	CardReaderKind    byte   // 读卡器类型
	CardReaderNo      uint32 // 读卡器编号
	DoorNo            uint32 // 门编号
	VerifyNo          uint32 // 多重卡认证序号
	AlarmInNo         uint32 // 报警输入号
	AlarmOutNo        uint32 // 报警输出号
	CaseSensorNo      uint32 // 事件触发器编号
	Rs485No           uint32 // RS485通道号
	MultiCardGroupNo  uint32 // 群组编号
	AccessChannel     uint16 // 通道号
	DeviceNo          byte   // 设备编号
	DistractControlNo byte   // 分控器编号
	EmployeeNo        uint32 // 工号
	LocalControllerID uint16 // 就地控制器编号
	InternetAccess    byte   // 网口ID
	Type              byte   // 防区类型
}

// Point 坐标点结构体
type Point struct {
	X float32 // X坐标
	Y float32 // Y坐标
}

// ACSEventInfoExtend 门禁事件扩展信息
type ACSEventInfoExtend struct {
	FrontSerialNo       uint32 // 事件流水号
	UserType            byte   // 人员类型：0-无效，1-普通人（主人），2-来宾（访客），3-禁止名单人，4-管理员
	CurrentVerifyMode   byte   // 读卡器当前验证方式
	CurrentEvent        byte   // 是否为实时事件：0-无效，1-是（实时事件），2-否（离线事件）
	PurePwdVerifyEnable byte   // 设备是否支持纯密码验证：0-不支持，1-支持
	EmployeeNo          string // 工号，人员ID
	AttendanceStatus    byte   // 考勤状态：0-未定义,1-上班，2-下班，3-开始休息，4-结束休息，5-开始加班，6-结束加班
	StatusValue         byte   // 考勤状态值
	UUID                string // UUID，该字段仅对接萤石平台设备才会使用
	DeviceName          string // 设备序列号
}

// ACSEventInfoExtendV20 门禁事件扩展信息V20
type ACSEventInfoExtendV20 struct {
	RemoteCheck          byte    // 是否需要远程核验（0-无效，1-不需要（默认），2-需要）
	ThermometryUnit      byte    // 测温单位：0-摄氏度（默认），1-华氏度，2-开尔文
	IsAbnomalTemperature byte    // 人脸抓拍测温是否温度异常：1-是，0-否
	CurrTemperature      float32 // 人脸温度，精确到小数点后一位
	RegionCoordinates    Point   // 人脸温度坐标
	QRCodeInfo           string  // 二维码信息
	VisibleLightData     []byte  // 热成像相机可见光图片
	ThermalData          []byte  // 热成像图片
	AttendanceLabel      string  // 考勤自定义标签
	XCoordinate          uint16  // x坐标，基于左上角，图片的归一化坐标，范围0-1000
	YCoordinate          uint16  // y坐标，基于左上角，图片的归一化坐标，范围0-1000
	Width                uint16  // 人脸宽度，范围0-1000
	Height               uint16  // 人脸高度，范围0-1000
	HealthCode           byte    // 健康码状态, 0-无效, 1-未申领, 2-未查询, 3-绿码, 4-黄码, 5-红码, 6-无此人员, 7-健康码信息接口异常（查询失败）, 8-查询超时超时
	NADCode              byte    // 核酸检测, 0-无效, 1-未查询到核酸检测, 2-核酸阴性（未过期）, 3-核酸阴性（已过期）, 4-核酸检测无效（已过期）
	TravelCode           byte    // 行程码, 0-无效, 1-14天内一直在低风险地区, 2-14天内离开过低风险地区, 3-14天内到达过中风险地区, 4-其他
	VaccineStatus        byte    // 疫苗状态, 0-无效, 1-未接种疫苗, 2-接种疫苗中, 3-完成疫苗接种
}

// ACSAlarmInfo 门禁主机报警信息
type ACSAlarmInfo struct {
	Size                  uint32                 // 结构体大小
	Major                 uint32                 // 报警主类型
	Minor                 uint32                 // 报警次类型
	Time                  Time                   // 报警时间
	NetUser               string                 // 网络操作的用户名
	RemoteHostAddr        IPAddr                 // 远程主机地址
	AcsEventInfo          ACSEventInfo           // 详细参数
	PicDataLen            uint32                 // 图片数据长度
	PicData               []byte                 // 图片数据
	PicUri                string                 // 图片URI（当PicTransType=1时有效）
	InductiveEventType    uint16                 // 归纳事件类型
	PicTransType          byte                   // 图片数据传输方式: 0-binary, 1-url
	IOTChannelNo          uint32                 // IOT通道号
	AcsEventInfoExtend    *ACSEventInfoExtend    // 门禁事件扩展信息
	TimeType              byte                   // 时间类型：0-设备本地时间，1-UTC时间
	AcsEventInfoExtendV20 *ACSEventInfoExtendV20 // 门禁事件扩展信息V20
}

// AlarmInfoV30 通用报警信息结构体
type AlarmInfoV30 struct {
	Size               uint32 // 结构体大小
	AlarmType          uint32 // 报警类型
	AlarmInputNumber   uint32 // 报警输入口，触发报警的输入口号
	AlarmOutputNumber  []byte // 触发报警的输出口
	AlarmRelateChannel []byte // 触发报警的通道
	Channel            []byte // dwAlarmType为0时，表示触发报警的模拟通道
	DiskNumber         []byte // dwAlarmType为2时，表示触发报警的硬盘号
}
