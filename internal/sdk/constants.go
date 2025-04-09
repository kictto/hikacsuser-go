package sdk

// SDK常量定义
const (
	// 通用常量
	//MAX_NAMELEN               = 16  // DVR本地登录名
	MAX_RIGHT                 = 32  // 设备支持的权限（1-12表示本地权限，13-32表示远程权限）
	NAME_LEN                  = 32  // 用户名长度
	PASSWD_LEN                = 16  // 密码长度
	SERIALNO_LEN              = 48  // 序列号长度
	MACADDR_LEN               = 6   // mac地址长度
	MAX_ETHERNET              = 2   // 设备可配以太网络
	PATHNAME_LEN              = 128 // 路径长度
	MAX_TIMESEGMENT_V30       = 8   // 9000设备最大时间段数
	MAX_TIMESEGMENT           = 4   // 8000设备最大时间段数
	MAX_SHELTERNUM            = 4   // 8000设备最大遮挡区域数
	MAX_DAYS                  = 7   // 每周天数
	PHONENUMBER_LEN           = 32  // pppoe拨号号码最大长度
	//MAX_DISKNUM_V30           = 33  // 9000设备最大硬盘数
	MAX_DISKNUM               = 16  // 8000设备最大硬盘数
	MAX_DISKNUM_V10           = 8   // 1.2版本之前版本
	MAX_WINDOW_V30            = 32  // 9000设备本地显示最大播放窗口数
	MAX_WINDOW                = 16  // 8000设备最大硬盘数
	MAX_VGA_V30               = 4   // 9000设备最大可接VGA数
	MAX_VGA                   = 1   // 8000设备最大可接VGA数
	MAX_USERNUM_V30           = 32  // 9000设备最大用户数
	MAX_USERNUM               = 16  // 8000设备最大用户数
	MAX_EXCEPTIONNUM_V30      = 32  // 9000设备最大异常处理数
	MAX_EXCEPTIONNUM          = 16  // 8000设备最大异常处理数
	MAX_LINK                  = 6   // 8000设备单通道最大视频流连接数
	MAX_ALARMIN_V30           = 160 // 9000设备最大报警输入数
	MAX_ALARMIN               = 16  // 8000设备最大报警输入数
	//MAX_ALARMOUT_V30          = 96  // 9000设备最大报警输出数
	MAX_ALARMOUT              = 4   // 8000设备最大报警输出数
	//MAX_CHANNUM_V30           = 64  // 9000设备最大通道数
	MAX_CHANNUM               = 16  // 8000设备最大通道数
	MAX_CARD_READER_NUM_512   = 512 // 最大读卡器数
	ERROR_MSG_LEN             = 32  // 下发错误信息
	MAX_FACE_NUM              = 2   // 最大人脸数
	MAX_FINGER_PRINT_LEN      = 768 // 最大指纹长度
	MAX_CARDNO_LEN            = 48  // 卡号最大长度
	NET_SDK_EMPLOYEE_NO_LEN   = 32  // 工号长度
	CARDNUM_LEN               = 20  // 卡号长度
	MAX_DOOR_NUM_256          = 256 // 最大门数256
	MAX_CASE_SENSOR_NUM       = 8   // 最大事件触发器数
	MAX_ALARMHOST_ALARMIN_NUM = 512 // 最大报警主机报警输入口数

	// 网络SDK错误码
	NET_DVR_NOERROR              = 0  // 没有错误
	NET_DVR_PASSWORD_ERROR       = 1  // 用户名密码错误
	NET_DVR_NOENOUGHPRI          = 2  // 权限不足
	NET_DVR_NOINIT               = 3  // 没有初始化
	NET_DVR_CHANNEL_ERROR        = 4  // 通道号错误
	NET_DVR_OVER_MAXLINK         = 5  // 连接到DVR的客户端个数超过最大
	NET_DVR_VERSIONNOMATCH       = 6  // 版本不匹配
	NET_DVR_NETWORK_FAIL_CONNECT = 7  // 连接服务器失败
	NET_DVR_NETWORK_SEND_ERROR   = 8  // 向服务器发送失败
	NET_DVR_NETWORK_RECV_ERROR   = 9  // 从服务器接收数据失败
	NET_DVR_NETWORK_RECV_TIMEOUT = 10 // 从服务器接收数据超时
	NET_DVR_NETWORK_ERRORDATA    = 11 // 传送的数据有误
	NET_DVR_ORDER_ERROR          = 12 // 调用次序错误
	NET_DVR_OPERNOPERMIT         = 13 // 无此权限
	NET_DVR_PARAMETER_ERROR      = 17 // 参数错误

	// 配置SDK初始化参数类型
	NET_SDK_INIT_CFG_SDK_PATH    = 2 // 设置HCNetSDK库所在目录
	NET_SDK_INIT_CFG_LIBEAY_PATH = 3 // 设置OpenSSL的libeay32.dll/libcrypto.so/libcrypto.dylib所在路径
	NET_SDK_INIT_CFG_SSLEAY_PATH = 4 // 设置OpenSSL的ssleay32.dll/libssl.so/libssl.dylib所在路径

	// 远程配置标志
	NET_SDK_CONFIG_STATUS_SUCCESS   = 1000 // 配置成功
	NET_SDK_CONFIG_STATUS_NEED_WAIT = 1001 // 配置等待
	NET_SDK_CONFIG_STATUS_FINISH    = 1002 // 配置完成
	NET_SDK_CONFIG_STATUS_FAILED    = 1003 // 配置失败
	NET_SDK_CONFIG_STATUS_EXCEPTION = 1004 // 配置异常
	
	// 获取下一个状态标志
	NET_SDK_GET_NEXT_STATUS_SUCCESS = 1000 // 获取成功
	NET_SDK_GET_NEXT_STATUS_NEED_WAIT = 1001 // 需要等待
	NET_SDK_NEXT_STATUS__FINISH = 1002 // 获取完成
	NET_SDK_GET_NEXT_STATUS_FAILED = 1003 // 获取失败

	// 设备登录模式
	NET_DVR_LOGIN_SUCCESS          = 1   // 登录成功
	NET_DVR_LOGIN_ERROR_PASSWORD   = 2   // 密码错误
	NET_DVR_LOGIN_ERROR_USER       = 3   // 用户名错误
	NET_DVR_LOGIN_ERROR_TIMEOUT    = 4   // 连接超时
	NET_DVR_LOGIN_ERROR_RELOGGIN   = 5   // 重复登录
	NET_DVR_LOGIN_ERROR_LOCKED     = 6   // 账号被锁定
	NET_DVR_LOGIN_ERROR_BLACKLIST  = 7   // 账号被列为黑名单
	NET_DVR_LOGIN_ERROR_BUSY       = 8   // 设备忙
	NET_DVR_LOGIN_ERROR_CONNECT    = 9   // 连接出错
	NET_DVR_DEV_ADDRESS_MAX_LEN    = 129 // 设备地址最大长度
	NET_DVR_LOGIN_USERNAME_MAX_LEN = 64  // 登录用户名最大长度
	NET_DVR_LOGIN_PASSWD_MAX_LEN   = 64  // 登录密码最大长度

	// ISAPI协议命令
	COMM_ISAPI_CONFIG        = 16010 // ISAPI协议命令
	NET_DVR_JSON_CONFIG      = 2550  // JSON配置命令
	NET_DVR_FACE_DATA_SEARCH = 2552  // 查询人脸库中的人脸数据

	// 门禁主机参数配置命令
	NET_DVR_GET_ACS_EVENT       = 2514  // 获取门禁事件
	NET_DVR_GET_ACS_CFG             = 2159 // 获取门禁主机参数
	NET_DVR_SET_ACS_CFG             = 2160 // 设置门禁主机参数
	NET_DVR_GET_ACS_WORK_STATUS_V50 = 2180 // 获取门禁主机工作状态

	// 报警主机相关常量
	MAX_ALARMHOST_ALARMOUT_NUM = 512 // 最大报警主机报警输出口数

	//ACS_CARD_NO_LEN = 32 // 门禁卡号长度
	NET_SDK_MONITOR_ID_LEN = 64 // 布防点ID长度
)