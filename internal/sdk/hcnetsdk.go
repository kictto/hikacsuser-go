package sdk

import (
	"unsafe"
)

// SDK常量定义 - 移至 types.go
// const (
// 	// 通用常量
// 	MAX_NAMELEN               = 16  // DVR本地登录名
// 	MAX_RIGHT                 = 32  // 设备支持的权限（1-12表示本地权限，13-32表示远程权限）
// 	NAME_LEN                  = 32  // 用户名长度
// 	PASSWD_LEN                = 16  // 密码长度
// 	SERIALNO_LEN              = 48  // 序列号长度
// 	MACADDR_LEN               = 6   // mac地址长度
// 	MAX_ETHERNET              = 2   // 设备可配以太网络
// 	PATHNAME_LEN              = 128 // 路径长度
// 	MAX_TIMESEGMENT_V30       = 8   // 9000设备最大时间段数
// 	MAX_TIMESEGMENT           = 4   // 8000设备最大时间段数
// 	MAX_SHELTERNUM            = 4   // 8000设备最大遮挡区域数
// 	MAX_DAYS                  = 7   // 每周天数
// 	PHONENUMBER_LEN           = 32  // pppoe拨号号码最大长度
// 	MAX_DISKNUM_V30           = 33  // 9000设备最大硬盘数
// 	MAX_DISKNUM               = 16  // 8000设备最大硬盘数
// 	MAX_DISKNUM_V10           = 8   // 1.2版本之前版本
// 	MAX_WINDOW_V30            = 32  // 9000设备本地显示最大播放窗口数
// 	MAX_WINDOW                = 16  // 8000设备最大硬盘数
// 	MAX_VGA_V30               = 4   // 9000设备最大可接VGA数
// 	MAX_VGA                   = 1   // 8000设备最大可接VGA数
// 	MAX_USERNUM_V30           = 32  // 9000设备最大用户数
// 	MAX_USERNUM               = 16  // 8000设备最大用户数
// 	MAX_EXCEPTIONNUM_V30      = 32  // 9000设备最大异常处理数
// 	MAX_EXCEPTIONNUM          = 16  // 8000设备最大异常处理数
// 	MAX_LINK                  = 6   // 8000设备单通道最大视频流连接数
// 	MAX_ALARMIN_V30           = 160 // 9000设备最大报警输入数
// 	MAX_ALARMIN               = 16  // 8000设备最大报警输入数
// 	MAX_ALARMOUT_V30          = 96  // 9000设备最大报警输出数
// 	MAX_ALARMOUT              = 4   // 8000设备最大报警输出数
// 	MAX_CHANNUM_V30           = 64  // 9000设备最大通道数
// 	MAX_CHANNUM               = 16  // 8000设备最大通道数
// 	MAX_CARD_READER_NUM_512   = 512 // 最大读卡器数
// 	ERROR_MSG_LEN             = 32  // 下发错误信息
// 	MAX_FACE_NUM              = 2   // 最大人脸数
// 	MAX_FINGER_PRINT_LEN      = 768 // 最大指纹长度
// 	MAX_CARDNO_LEN            = 48  // 卡号最大长度
// 	NET_SDK_EMPLOYEE_NO_LEN   = 32  // 工号长度
// 	CARDNUM_LEN               = 20  // 卡号长度
// 	MAX_DOOR_NUM_256          = 256 // 最大门数256
// 	MAX_CASE_SENSOR_NUM       = 8   // 最大事件触发器数
// 	MAX_ALARMHOST_ALARMIN_NUM = 512 // 最大报警主机报警输入口数
//
// 	// 网络SDK错误码
// 	NET_DVR_NOERROR              = 0  // 没有错误
// 	NET_DVR_PASSWORD_ERROR       = 1  // 用户名密码错误
// 	NET_DVR_NOENOUGHPRI          = 2  // 权限不足
// 	NET_DVR_NOINIT               = 3  // 没有初始化
// 	NET_DVR_CHANNEL_ERROR        = 4  // 通道号错误
// 	NET_DVR_OVER_MAXLINK         = 5  // 连接到DVR的客户端个数超过最大
// 	NET_DVR_VERSIONNOMATCH       = 6  // 版本不匹配
// 	NET_DVR_NETWORK_FAIL_CONNECT = 7  // 连接服务器失败
// 	NET_DVR_NETWORK_SEND_ERROR   = 8  // 向服务器发送失败
// 	NET_DVR_NETWORK_RECV_ERROR   = 9  // 从服务器接收数据失败
// 	NET_DVR_NETWORK_RECV_TIMEOUT = 10 // 从服务器接收数据超时
// 	NET_DVR_NETWORK_ERRORDATA    = 11 // 传送的数据有误
// 	NET_DVR_ORDER_ERROR          = 12 // 调用次序错误
// 	NET_DVR_OPERNOPERMIT         = 13 // 无此权限
// 	NET_DVR_PARAMETER_ERROR      = 17 // 参数错误
//
// 	// 配置SDK初始化参数类型
// 	NET_SDK_INIT_CFG_SDK_PATH    = 2 // 设置HCNetSDK库所在目录
// 	NET_SDK_INIT_CFG_LIBEAY_PATH = 3 // 设置OpenSSL的libeay32.dll/libcrypto.so/libcrypto.dylib所在路径
// 	NET_SDK_INIT_CFG_SSLEAY_PATH = 4 // 设置OpenSSL的ssleay32.dll/libssl.so/libssl.dylib所在路径
//
// 	// 远程配置标志
// 	NET_SDK_CONFIG_STATUS_SUCCESS   = 1000 // 配置成功
// 	NET_SDK_CONFIG_STATUS_NEED_WAIT = 1001 // 配置等待
// 	NET_SDK_CONFIG_STATUS_FINISH    = 1002 // 配置完成
// 	NET_SDK_CONFIG_STATUS_FAILED    = 1003 // 配置失败
// 	NET_SDK_CONFIG_STATUS_EXCEPTION = 1004 // 配置异常
//
// 	// 获取下一个状态标志
// 	NET_SDK_GET_NEXT_STATUS_SUCCESS = 1000 // 获取成功
// 	NET_SDK_GET_NEXT_STATUS_NEED_WAIT = 1001 // 需要等待
// 	NET_SDK_NEXT_STATUS__FINISH = 1002 // 获取完成
// 	NET_SDK_GET_NEXT_STATUS_FAILED = 1003 // 获取失败
//
// 	// 设备登录模式
// 	NET_DVR_LOGIN_SUCCESS          = 1   // 登录成功
// 	NET_DVR_LOGIN_ERROR_PASSWORD   = 2   // 密码错误
// 	NET_DVR_LOGIN_ERROR_USER       = 3   // 用户名错误
// 	NET_DVR_LOGIN_ERROR_TIMEOUT    = 4   // 连接超时
// 	NET_DVR_LOGIN_ERROR_RELOGGIN   = 5   // 重复登录
// 	NET_DVR_LOGIN_ERROR_LOCKED     = 6   // 账号被锁定
// 	NET_DVR_LOGIN_ERROR_BLACKLIST  = 7   // 账号被列为黑名单
// 	NET_DVR_LOGIN_ERROR_BUSY       = 8   // 设备忙
// 	NET_DVR_LOGIN_ERROR_CONNECT    = 9   // 连接出错
// 	NET_DVR_DEV_ADDRESS_MAX_LEN    = 129 // 设备地址最大长度
// 	NET_DVR_LOGIN_USERNAME_MAX_LEN = 64  // 登录用户名最大长度
// 	NET_DVR_LOGIN_PASSWD_MAX_LEN   = 64  // 登录密码最大长度
//
// 	// ISAPI协议命令
// 	COMM_ISAPI_CONFIG        = 16010 // ISAPI协议命令
// 	NET_DVR_JSON_CONFIG      = 2550  // JSON配置命令
// 	NET_DVR_FACE_DATA_SEARCH = 2552  // 查询人脸库中的人脸数据
//
// 	// 门禁主机参数配置命令
// 	NET_DVR_GET_ACS_EVENT       = 2514  // 获取门禁事件
// 	NET_DVR_GET_ACS_CFG             = 2159 // 获取门禁主机参数
// 	NET_DVR_SET_ACS_CFG             = 2160 // 设置门禁主机参数
// 	NET_DVR_GET_ACS_WORK_STATUS_V50 = 2180 // 获取门禁主机工作状态
//
// 	// 报警主机相关常量
// 	MAX_ALARMHOST_ALARMOUT_NUM = 512 // 最大报警主机报警输出口数
//
// 	ACS_CARD_NO_LEN = 32 // 门禁卡号长度
// 	NET_SDK_MONITOR_ID_LEN = 64 // 布防点ID长度
//
// )

// BYTE_ARRAY 类型定义，用于接口传递字节数组
type BYTE_ARRAY struct {
	Size    uint32
	Buffer  []byte
	Pointer unsafe.Pointer
}

// NET_DVR_USER_LOGIN_INFO 登录信息结构体
type NET_DVR_USER_LOGIN_INFO struct {
	SDeviceAddress [NET_DVR_DEV_ADDRESS_MAX_LEN]byte    // 设备地址，IP或者普通域名
	ByUseTransport byte                                 // 传输模式，与Java版本保持一致
	WPort          uint16                               // 设备端口号
	SUserName      [NET_DVR_LOGIN_USERNAME_MAX_LEN]byte // 登录用户名
	SPassword      [NET_DVR_LOGIN_PASSWD_MAX_LEN]byte   // 登录密码
	BUseAsynLogin  bool                                 // 是否异步登录
	Reserved       [128]byte                            // 保留
	ByLoginMode    byte                                 // 登录模式
	ByHttps        byte                                 // HTTPS登录标志
	// 以下字段新增，与Java版本保持一致
	ByProxyType  byte      // 0:不使用代理，1：使用标准代理，2：使用EHome代理
	ByUseUTCTime byte      // 时间格式 0-不进行转换，默认,1-转换成UTC，2-转换成本地时间
	IProxyID     int32     // 代理服务器序号
	ByVerifyMode byte      // 认证方式，0-不认证，1-双向认证，2-单向认证
	ByRes2       [119]byte // 保留
}

// NET_DVR_DEVICEINFO_V30 设备信息
type NET_DVR_DEVICEINFO_V30 struct {
	SSerialNumber       [SERIALNO_LEN]byte // 序列号
	ByAlarmInPortNum    byte               // 报警输入个数
	ByAlarmOutPortNum   byte               // 报警输出个数
	ByDiskNum           byte               // 硬盘个数
	ByDVRType           byte               // 设备类型
	ByChanNum           byte               // 模拟通道个数
	ByStartChan         byte               // 起始通道号
	ByIPChanNum         byte               // 数字通道个数
	ByAudioChanNum      byte               // 音频通道个数
	ByIPAlarmInPortNum  byte               // IP报警输入个数
	ByIPAlarmOutPortNum byte               // IP报警输出个数
	ByZeroChanNum       byte               // 零通道编码个数
	ByMainProto         byte               // 主码流传输协议类型
	BySubProto          byte               // 子码流传输协议类型
	BySupport           byte               // 能力集扩展
	BySupport1          byte               // 能力集扩展1
	BySupport2          byte               // 能力集扩展2
	BySupport3          byte               // 能力集扩展3
	ByRes               [22]byte           // 保留
}

// NET_DVR_DEVICEINFO_V40 设备信息V40结构
type NET_DVR_DEVICEINFO_V40 struct {
	StructSize             uint32                 // 结构体大小
	DeviceInfo             NET_DVR_DEVICEINFO_V30 // 设备信息V30结构
	BySupportLock          byte                   // 是否支持锁定功能
	ByRetryLoginTime       byte                   // 重试登录次数
	ByPasswordLevel        byte                   // 密码安全等级
	ByRes1                 byte                   // 保留字节
	DwSurplusLockTime      uint32                 // 剩余锁定时间
	ByCharEncodeType       byte                   // 字符编码类型
	BySupportDev5          byte                   // 是否支持v50版本参数
	BySupport              byte                   // 能力扩展
	ByLoginMode            byte                   // 登录模式
	DwOEMCode              uint32                 // OEM码
	IResidualValidity      int32                  // 密码剩余有效期
	ByResidualValidity     byte                   // 剩余有效期是否有效
	BySingleStartDTalkChan byte                   // 独立音轨接入起始通道
	BySingleDTalkChanNums  byte                   // 独立音轨接入通道数量
	ByPassWordResetLevel   byte                   // 密码重置等级
	BySupportStreamEncrypt byte                   // 是否支持码流加密
	ByMarketType           byte                   // 市场类型
	ByRes2                 [238]byte              // 保留字节
}

// NET_DVR_LOCAL_SDK_PATH SDK本地路径
type NET_DVR_LOCAL_SDK_PATH struct {
	SPath [256]byte // SDK路径
}

// RemoteConfigCallback 远程配置回调函数类型
type RemoteConfigCallback func(dwType uint32, lpBuffer unsafe.Pointer, dwBufLen uint32, pUserData unsafe.Pointer)

// MSGCallBack_V31 报警回调函数签名 V31/V50
type MSGCallBack_V31 func(lCommand int, pAlarmer *NET_DVR_ALARMER, pAlarmInfo unsafe.Pointer, dwBufLen uint32, pUserData unsafe.Pointer) bool

// NET_DVR_ALARMER 报警设备信息结构体
type NET_DVR_ALARMER struct {
	ByUserIDValid     byte               // userid是否有效 0-无效，1-有效
	BySerialValid     byte               // 序列号是否有效 0-无效，1-有效
	ByVersionValid    byte               // 版本号是否有效 0-无效，1-有效
	ByDeviceNameValid byte               // 设备名字是否有效 0-无效，1-有效
	ByMacAddrValid    byte               // MAC地址是否有效 0-无效，1-有效
	ByLinkPortValid   byte               // login端口是否有效 0-无效，1-有效
	ByDeviceIPValid   byte               // 设备IP是否有效 0-无效，1-有效
	BySocketIPValid   byte               // socket ip是否有效 0-无效，1-有效
	LUserID           int32              // NET_DVR_Login()返回值, 布防时有效
	SSerialNumber     [SERIALNO_LEN]byte // 序列号
	DwDeviceVersion   uint32             // 版本信息 高16位表示主版本，低16位表示次版本
	SDeviceName       [NAME_LEN]byte     // 设备名字
	ByMacAddr         [MACADDR_LEN]byte  // MAC地址
	WLinkPort         uint16             // link port
	SDeviceIP         [128]byte          // IP地址
	SSocketIP         [128]byte          // 报警主动上传时的socket IP地址
	ByIpProtocol      byte               // Ip协议 0-IPV4, 1-IPV6
	ByRes2            [11]byte
}

// NET_DVR_SETUPALARM_PARAM 布防参数结构体 V41
type NET_DVR_SETUPALARM_PARAM struct {
	DwSize               uint32 // 结构体大小
	ByLevel              byte   // 布防优先级，0-一等级（高），1-二等级（中），2-三等级（低）
	ByAlarmInfoType      byte   // 报警信息上传方式：0-老报警信息（NET_DVR_ALARMINFO），1-新报警信息(NET_DVR_ALARMINFO_V30)
	ByRetAlarmTypeV40    byte   // V40报警信息类型
	ByRetDevInfoVersion  byte   // V40报警信息对应设备信息版本号
	ByRetVQDAlarmType    byte   // VQD报警上传类型（用于报警类型区分）
	ByFaceAlarmDetection byte   // 人脸报警信息类型
	BySupport            byte
	ByBrokenNetHttp      byte
	WSeverityFilter      uint16 // 严重程度，用于SMART IPC
	BySnapTimes          byte   // 设备联动抓图次数
	BySnapSeq            byte   // 设备联动抓图序号
	ByRelRecordChan      byte
	ByRes1               [12]byte
	ByChannel            byte // 触发报警的通道号
	ByRes                [35]byte
}

// HCNetSDK 定义了SDK的主要函数接口
type HCNetSDK interface {
	// 初始化和清理
	NET_DVR_Init() bool
	NET_DVR_Cleanup() bool
	NET_DVR_SetLogToFile(LogLevel int, LogDir string, bAutoDel bool) bool
	NET_DVR_SetSDKInitCfg(enumType int, lpInBuff unsafe.Pointer) bool

	// 登录注销
	NET_DVR_Login_V40(pLoginInfo *NET_DVR_USER_LOGIN_INFO, lpDeviceInfo *NET_DVR_DEVICEINFO_V40) int
	NET_DVR_Logout(lUserID int) bool

	// 参数配置
	NET_DVR_GetDVRConfig(lUserID int, dwCommand uint32, lChannel int,
		lpOutBuffer unsafe.Pointer, dwOutBufferSize uint32, lpBytesReturned *uint32) bool
	NET_DVR_SetDVRConfig(lUserID int, dwCommand uint32, lChannel int,
		lpInBuffer unsafe.Pointer, dwInBufferSize uint32) bool

	// 远程配置
	NET_DVR_StartRemoteConfig(lUserID int, dwCommand uint32,
		lpInBuffer unsafe.Pointer, dwInBufferSize uint32,
		fRemoteConfigCallback uintptr, pUserData unsafe.Pointer) int64
	NET_DVR_StopRemoteConfig(lHandle int64) bool
	NET_DVR_SendWithRecvRemoteConfig(lHandle int64, lpInBuff unsafe.Pointer, dwInBuffSize uint32,
		lpOutBuff unsafe.Pointer, dwOutBuffSize uint32, dwOutDataLen *uint32) int
	NET_DVR_GetNextRemoteConfig(lHandle int64, lpOutBuff unsafe.Pointer, dwOutBuffSize uint32, lpOutDataLen *uint32) int

	// 错误处理
	NET_DVR_GetLastError() uint32
	NET_DVR_GetErrorMsg(pErrorNo *int32) string

	// 门禁控制
	NET_DVR_ControlGateway(lUserID int, lGatewayIndex int, dwStaic uint32) bool

	// 报警布防
	NET_DVR_SetupAlarmChan_V41(lUserID int, lpSetupParam *NET_DVR_SETUPALARM_PARAM) int
	NET_DVR_CloseAlarmChan_V30(lAlarmHandle int) bool
	NET_DVR_SetDVRMessageCallBack_V50(iIndex int, fMessageCallBack MSGCallBack_V31, pUser unsafe.Pointer) bool

	// 图片查找和获取
	NET_DVR_FindPicture(lUserID int, pFindParam *NET_DVR_FIND_PICTURE_PARAM) int
	NET_DVR_FindNextPicture(lFindHandle int, lpFindData *NET_DVR_FIND_PICTURE) int
	NET_DVR_CloseFindPicture(lFindHandle int) bool
	NET_DVR_GetPicture_V50(lUserID int, pPicParam *NET_DVR_FIND_PICTURE, pParam *NET_DVR_GETPIC_PARAM) bool
}

// HCNetSDKImpl SDK实现
// type HCNetSDKImpl struct {
// 	// SDK库实例
// }

// // LoadLibrary 加载SDK库
// func LoadLibrary(dllPath string) (*HCNetSDKImpl, error) {
// 	// 实现动态库加载逻辑
// 	return &HCNetSDKImpl{}, nil
// }

// // GetSDKInstance 获取SDK实例
// func GetSDKInstance() (HCNetSDK, error) {
// 	dllPath := utils.GetDLLPath()
// 	instance, err := LoadLibrary(dllPath)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return instance, nil
// }
