package acsapi

// 报警命令常量
const (
	// 报警命令类型
	COMM_ALARM_V30 = 0x4000 // 通用报警信息
	COMM_ALARM_ACS = 0x5002 // 门禁主机报警信息

	// 报警主类型
	ALARM_MAJOR_ALARM     = 0x1 // 报警
	ALARM_MAJOR_EXCEPTION = 0x2 // 异常
	ALARM_MAJOR_OPERATION = 0x3 // 操作
	ALARM_MAJOR_EVENT     = 0x5 // 事件

	// 报警次类型 - 报警类型(ALARM_MAJOR_ALARM)
	ALARM_MINOR_ZONE_SHORT_CIRCUIT        = 0x400 // 防区短路报警
	ALARM_MINOR_ZONE_BROKEN_CIRCUIT       = 0x401 // 防区断路报警
	ALARM_MINOR_ZONE_ABNORMAL             = 0x402 // 防区异常报警
	ALARM_MINOR_ZONE_RESTORE              = 0x403 // 防区报警恢复
	ALARM_MINOR_DEVICE_TAMPER             = 0x404 // 设备防拆报警
	ALARM_MINOR_DEVICE_TAMPER_RESTORE     = 0x405 // 设备防拆恢复
	ALARM_MINOR_READER_TAMPER             = 0x406 // 读卡器防拆报警
	ALARM_MINOR_READER_TAMPER_RESTORE     = 0x407 // 读卡器防拆恢复
	ALARM_MINOR_EVENT_INPUT               = 0x408 // 事件输入报警
	ALARM_MINOR_EVENT_INPUT_RESTORE       = 0x409 // 事件输入恢复
	ALARM_MINOR_DURESS                    = 0x40a // 胁迫报警
	ALARM_MINOR_OFFLINE_ECENT_NEARLY_FULL = 0x40b // 离线事件满90%报警
	ALARM_MINOR_CARD_AUTH_FAIL_EXCEED     = 0x40c // 卡号认证失败超次报警
	ALARM_MINOR_SD_CARD_FULL              = 0x40d // SD卡存储满报警
	ALARM_MINOR_LINKAGE_CAPTURE           = 0x40e // 联动抓拍事件报警

	// 报警次类型 - 异常类型(ALARM_MAJOR_EXCEPTION)
	EXCEPTION_MINOR_NETWORK_BROKEN      = 0x27  // 网络断开
	EXCEPTION_MINOR_RS485_ABNORMAL      = 0x3a  // RS485连接状态异常
	EXCEPTION_MINOR_RS485_RESTORE       = 0x3b  // RS485连接状态异常恢复
	EXCEPTION_MINOR_DEVICE_POWER_ON     = 0x400 // 设备上电启动
	EXCEPTION_MINOR_DEVICE_POWER_OFF    = 0x401 // 设备掉电关闭
	EXCEPTION_MINOR_WATCHDOG_RESET      = 0x402 // 看门狗复位
	EXCEPTION_MINOR_LOW_BATTERY         = 0x403 // 蓄电池电压低
	EXCEPTION_MINOR_BATTERY_RESTORE     = 0x404 // 蓄电池电压恢复正常
	EXCEPTION_MINOR_AC_OFF              = 0x405 // 交流电断电
	EXCEPTION_MINOR_AC_RESTORE          = 0x406 // 交流电恢复
	EXCEPTION_MINOR_NETWORK_RESTORE     = 0x407 // 网络恢复
	EXCEPTION_MINOR_FLASH_ABNORMAL      = 0x408 // FLASH读写异常
	EXCEPTION_MINOR_CARD_READER_OFFLINE = 0x409 // 读卡器掉线
	EXCEPTION_MINOR_CARD_READER_RESTORE = 0x40a // 读卡器掉线恢复

	// 报警次类型 - 事件类型(ALARM_MAJOR_EVENT)
	EVENT_MINOR_LEGAL_CARD_PASS                = 0x01 // 合法卡认证通过
	EVENT_MINOR_CARD_AND_PSW_PASS              = 0x02 // 刷卡加密码认证通过
	EVENT_MINOR_CARD_AND_PSW_FAIL              = 0x03 // 刷卡加密码认证失败
	EVENT_MINOR_CARD_AND_PSW_TIMEOUT           = 0x04 // 刷卡加密码认证超时
	EVENT_MINOR_CARD_NO_RIGHT                  = 0x05 // 卡无权限
	EVENT_MINOR_CARD_INVALID_PERIOD            = 0x06 // 卡不在有效期
	EVENT_MINOR_CARD_OUT_OF_DATE               = 0x07 // 卡号过期
	EVENT_MINOR_INVALID_CARD                   = 0x08 // 无效卡号
	EVENT_MINOR_ANTI_SNEAK_FAIL                = 0x09 // 反潜回认证失败
	EVENT_MINOR_INTERLOCK_DOOR_NOT_CLOSE       = 0x0a // 互锁门未关闭
	EVENT_MINOR_NOT_BELONG_MULTI_GROUP         = 0x0b // 不属于多重认证群组
	EVENT_MINOR_INVALID_MULTI_VERIFY_PERIOD    = 0x0c // 多重认证时间段不正确
	EVENT_MINOR_MULTI_VERIFY_SUPER_RIGHT_FAIL  = 0x0d // 多重认证超级权限认证失败
	EVENT_MINOR_MULTI_VERIFY_REMOTE_RIGHT_FAIL = 0x0e // 多重认证远程认证失败
	EVENT_MINOR_MULTI_VERIFY_SUCCESS           = 0x0f // 多重认证成功
	EVENT_MINOR_LEADER_CARD_OPEN_BEGIN         = 0x10 // 首卡开门开始
	EVENT_MINOR_LEADER_CARD_OPEN_END           = 0x11 // 首卡开门结束
	EVENT_MINOR_ALWAYS_OPEN_BEGIN              = 0x12 // 常开状态开始
	EVENT_MINOR_ALWAYS_OPEN_END                = 0x13 // 常开状态结束
	EVENT_MINOR_LOCK_OPEN                      = 0x14 // 门锁打开
	EVENT_MINOR_LOCK_CLOSE                     = 0x15 // 门锁关闭
	EVENT_MINOR_DOOR_BUTTON_PRESS              = 0x16 // 开门按钮打开
	EVENT_MINOR_DOOR_BUTTON_RELEASE            = 0x17 // 开门按钮释放
	EVENT_MINOR_DOOR_OPEN_NORMAL               = 0x18 // 正常开门
	EVENT_MINOR_DOOR_CLOSE_NORMAL              = 0x19 // 正常关门
	EVENT_MINOR_DOOR_OPEN_ABNORMAL             = 0x1a // 门异常打开
	EVENT_MINOR_DOOR_OPEN_TIMEOUT              = 0x1b // 门打开超时
	EVENT_MINOR_ALARMOUT_ON                    = 0x1c // 报警输出打开
	EVENT_MINOR_ALARMOUT_OFF                   = 0x1d // 报警输出关闭
	EVENT_MINOR_DOOR_CLOSE_ABNORMAL            = 0x1e // 门异常关闭
	EVENT_MINOR_LEGAL_FACE_PASS                = 0x88 // 人脸认证通过
	EVENT_MINOR_FACE_AND_CARD_PASS             = 0x89 // 人脸加卡认证通过
	EVENT_MINOR_FACE_AND_CARD_FAIL             = 0x8a // 人脸加卡认证失败
	EVENT_MINOR_FACE_AND_CARD_TIMEOUT          = 0x8b // 人脸加卡认证超时
	EVENT_MINOR_FACE_AND_PSW_PASS              = 0x8c // 人脸加密码认证通过
	EVENT_MINOR_FACE_AND_PSW_FAIL              = 0x8d // 人脸加密码认证失败
	EVENT_MINOR_FACE_AND_PSW_TIMEOUT           = 0x8e // 人脸加密码认证超时
	EVENT_MINOR_FINGERPRINT_PASS               = 0x8f // 指纹认证通过
	EVENT_MINOR_FINGERPRINT_FAIL               = 0x90 // 指纹认证失败
	EVENT_MINOR_FINGERPRINT_TIMEOUT            = 0x91 // 指纹认证超时
	EVENT_MINOR_FINGERPRINT_CARD_PASS          = 0x92 // 指纹加卡认证通过
	EVENT_MINOR_FINGERPRINT_CARD_FAIL          = 0x93 // 指纹加卡认证失败
	EVENT_MINOR_FINGERPRINT_CARD_TIMEOUT       = 0x94 // 指纹加卡认证超时
	EVENT_MINOR_FINGERPRINT_PSW_PASS           = 0x95 // 指纹加密码认证通过
	EVENT_MINOR_FINGERPRINT_PSW_FAIL           = 0x96 // 指纹加密码认证失败
	EVENT_MINOR_FINGERPRINT_PSW_TIMEOUT        = 0x97 // 指纹加密码认证超时
	EVENT_MINOR_MULTI_VERIFY_TIMEOUT           = 0x98 // 多重认证超时
	EVENT_MINOR_FACE_RECOGNIZE_FAIL            = 0x99 // 人脸识别失败
	EVENT_MINOR_FINGERPRINT_RECOGNIZE_FAIL     = 0x9a // 指纹识别失败
)
