package acsapi

import "time"

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
