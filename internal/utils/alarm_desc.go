package utils

// GetAlarmMajorTypeDesc 获取报警主类型描述
func GetAlarmMajorTypeDesc(majorType int) string {
	switch majorType {
	case 0x1:
		return "报警"
	case 0x2:
		return "异常"
	case 0x3:
		return "操作"
	case 0x5:
		return "事件"
	default:
		return "未知类型"
	}
}

// GetAlarmMinorTypeDesc 获取报警次类型描述
func GetAlarmMinorTypeDesc(majorType, minorType int) string {
	// 根据主类型选择不同的次类型映射
	switch majorType {
	case 0x1: // 报警
		return getAlarmMinorDesc(minorType)
	case 0x2: // 异常
		return getExceptionMinorDesc(minorType)
	case 0x3: // 操作
		return getOperationMinorDesc(minorType)
	case 0x5: // 事件
		return getEventMinorDesc(minorType)
	default:
		return "未知次类型"
	}
}

// getAlarmMinorDesc 获取报警类型的次类型描述
func getAlarmMinorDesc(minorType int) string {
	switch minorType {
	case 0x400:
		return "防区短路报警"
	case 0x401:
		return "防区断路报警"
	case 0x402:
		return "防区异常报警"
	case 0x403:
		return "防区报警恢复"
	case 0x404:
		return "设备防拆报警"
	case 0x405:
		return "设备防拆恢复"
	case 0x406:
		return "读卡器防拆报警"
	case 0x407:
		return "读卡器防拆恢复"
	case 0x408:
		return "事件输入报警"
	case 0x409:
		return "事件输入恢复"
	case 0x40a:
		return "胁迫报警"
	case 0x40b:
		return "离线事件满90%报警"
	case 0x40c:
		return "卡号认证失败超次报警"
	case 0x40d:
		return "SD卡存储满报警"
	case 0x40e:
		return "联动抓拍事件报警"
	case 0x40f:
		return "门控安全模块防拆报警"
	case 0x410:
		return "门控安全模块防拆恢复"
	case 0x411:
		return "POS开启"
	case 0x412:
		return "POS结束"
	case 0x413:
		return "人脸图像画质低"
	case 0x414:
		return "指纹图像画质低"
	case 0x415:
		return "消防输入短路报警"
	case 0x416:
		return "消防输入断路报警"
	case 0x417:
		return "消防输入恢复"
	case 0x418:
		return "消防按钮触发"
	case 0x419:
		return "消防按钮恢复"
	case 0x41a:
		return "维护按钮触发"
	case 0x41b:
		return "维护按钮恢复"
	case 0x41c:
		return "紧急按钮触发"
	case 0x41d:
		return "紧急按钮恢复"
	case 0x41e:
		return "分控器防拆报警"
	case 0x41f:
		return "分控器防拆报警恢复"
	case 0x422:
		return "通道控制器防拆报警"
	case 0x423:
		return "通道控制器防拆报警恢复"
	case 0x424:
		return "通道控制器消防输入报警"
	case 0x425:
		return "通道控制器消防输入报警恢复"
	case 0x442:
		return "合法事件满90%报警"
	case 0x95d:
		return "智能锁防劫持报警"
	default:
		// 处理自定义报警类型 (0x900-0x93f)
		if minorType >= 0x900 && minorType <= 0x93f {
			customNum := minorType - 0x900 + 1
			return "门禁自定义报警" + string(rune('0'+customNum/10)) + string(rune('0'+customNum%10))
		}
		return "未知报警次类型"
	}
}

// getExceptionMinorDesc 获取异常类型的次类型描述
func getExceptionMinorDesc(minorType int) string {
	switch minorType {
	case 0x27:
		return "网络断开"
	case 0x3a:
		return "RS485连接状态异常"
	case 0x3b:
		return "RS485连接状态异常恢复"
	case 0x400:
		return "设备上电启动"
	case 0x401:
		return "设备掉电关闭"
	case 0x402:
		return "看门狗复位"
	case 0x403:
		return "蓄电池电压低"
	case 0x404:
		return "蓄电池电压恢复正常"
	case 0x405:
		return "交流电断电"
	case 0x406:
		return "交流电恢复"
	case 0x407:
		return "网络恢复"
	case 0x408:
		return "FLASH读写异常"
	case 0x409:
		return "读卡器掉线"
	case 0x40a:
		return "读卡器掉线恢复"
	case 0x40b:
		return "指示灯关闭"
	case 0x40c:
		return "指示灯恢复"
	case 0x40d:
		return "通道控制器掉线"
	case 0x40e:
		return "通道控制器恢复"
	case 0x40f:
		return "门控安全模块掉线"
	case 0x410:
		return "门控安全模块掉线恢复"
	case 0x411:
		return "电池电压低（仅人脸设备使用）"
	case 0x412:
		return "电池电压恢复正常（仅人脸设备使用）"
	case 0x413:
		return "就地控制器网络断开"
	case 0x414:
		return "就地控制器网络恢复"
	case 0x415:
		return "主控RS485环路节点断开"
	case 0x416:
		return "主控RS485环路节点恢复"
	case 0x417:
		return "就地控制器掉线"
	case 0x418:
		return "就地控制器掉线恢复"
	case 0x419:
		return "就地下行RS485环路断开"
	case 0x41a:
		return "就地下行RS485环路恢复"
	case 0x41b:
		return "分控器在线"
	case 0x41c:
		return "分控器离线"
	case 0x41d:
		return "身份证阅读器未连接（智能专用）"
	case 0x41e:
		return "身份证阅读器连接恢复（智能专用）"
	case 0x41f:
		return "指纹模组未连接（智能专用）"
	case 0x420:
		return "指纹模组连接恢复（智能专用）"
	case 0x421:
		return "摄像头未连接"
	case 0x422:
		return "摄像头连接恢复"
	case 0x423:
		return "COM口未连接"
	case 0x424:
		return "COM口连接恢复"
	case 0x425:
		return "设备未授权"
	case 0x426:
		return "人证设备在线"
	case 0x427:
		return "人证设备离线"
	case 0x428:
		return "本地登录锁定"
	case 0x429:
		return "本地登录解锁"
	case 0x42a:
		return "与反潜回服务器通信断开"
	case 0x42b:
		return "与反潜回服务器通信恢复"
	case 0x42c:
		return "电机或传感器异常"
	case 0x42d:
		return "CAN总线异常"
	case 0x42e:
		return "CAN总线恢复"
	case 0x42f:
		return "闸机腔体温度超限"
	case 0x430:
		return "红外对射异常"
	case 0x431:
		return "红外对射恢复"
	case 0x432:
		return "灯板通信异常"
	case 0x433:
		return "灯板通信恢复"
	case 0x434:
		return "红外转接板通信异常"
	case 0x435:
		return "红外转接板通信恢复"
	default:
		// 处理自定义异常类型 (0x900-0x93f)
		if minorType >= 0x900 && minorType <= 0x93f {
			customNum := minorType - 0x900 + 1
			return "门禁自定义异常" + string(rune('0'+customNum/10)) + string(rune('0'+customNum%10))
		}
		return "未知异常次类型"
	}
}

// getOperationMinorDesc 获取操作类型的次类型描述
func getOperationMinorDesc(minorType int) string {
	switch minorType {
	case 0x50:
		return "本地登陆"
	case 0x51:
		return "本地注销登陆"
	case 0x5a:
		return "本地升级"
	case 0x70:
		return "远程登录"
	case 0x71:
		return "远程注销登陆"
	case 0x79:
		return "远程布防"
	case 0x7a:
		return "远程撤防"
	case 0x7b:
		return "远程重启"
	case 0x7e:
		return "远程升级"
	case 0x86:
		return "远程导出配置文件"
	case 0x87:
		return "远程导入配置文件"
	case 0xd6:
		return "远程手动开启报警输出"
	case 0xd7:
		return "远程手动关闭报警输出"
	case 0x400:
		return "远程开门"
	case 0x401:
		return "远程关门（对于梯控，表示受控）"
	case 0x402:
		return "远程常开（对于梯控，表示自由）"
	case 0x403:
		return "远程常关（对于梯控，表示禁用）"
	case 0x404:
		return "远程手动校时"
	case 0x405:
		return "NTP自动校时"
	case 0x406:
		return "远程清空卡号"
	case 0x407:
		return "远程恢复默认参数"
	case 0x408:
		return "防区布防"
	case 0x409:
		return "防区撤防"
	case 0x40a:
		return "本地恢复默认参数"
	case 0x40b:
		return "远程抓拍"
	case 0x40c:
		return "修改网络中心参数配置"
	case 0x40d:
		return "修改GPRS中心参数配置"
	case 0x40e:
		return "修改中心组参数配置"
	case 0x40f:
		return "解除码输入"
	case 0x410:
		return "自动重新编号"
	case 0x411:
		return "自动补充编号"
	case 0x412:
		return "导入普通配置文件"
	case 0x413:
		return "导出普通配置文件"
	case 0x414:
		return "导入卡权限参数"
	case 0x415:
		return "导出卡权限参数"
	case 0x416:
		return "本地U盘升级"
	case 0x417:
		return "访客呼梯"
	case 0x418:
		return "住户呼梯"
	case 0x419:
		return "远程实时布防"
	case 0x41a:
		return "远程实时撤防"
	case 0x41b:
		return "遥控器未对码操作失败"
	case 0x41c:
		return "遥控器关门"
	case 0x41d:
		return "遥控器开门"
	case 0x41e:
		return "遥控器常开门"
	default:
		// 处理自定义操作类型 (0x900-0x93f)
		if minorType >= 0x900 && minorType <= 0x93f {
			customNum := minorType - 0x900 + 1
			return "门禁自定义操作" + string(rune('0'+customNum/10)) + string(rune('0'+customNum%10))
		}
		return "未知操作次类型"
	}
}

// getEventMinorDesc 获取事件类型的次类型描述
func getEventMinorDesc(minorType int) string {
	switch minorType {
	case 0x01:
		return "合法卡认证通过"
	case 0x02:
		return "刷卡加密码认证通过"
	case 0x03:
		return "刷卡加密码认证失败"
	case 0x04:
		return "数卡加密码认证超时"
	case 0x05:
		return "刷卡加密码超次"
	case 0x06:
		return "未分配权限"
	case 0x07:
		return "无效时段"
	case 0x08:
		return "卡号过期"
	case 0x09:
		return "无此卡号"
	case 0x0a:
		return "反潜回认证失败"
	case 0x0b:
		return "互锁门未关闭"
	case 0x0c:
		return "卡不属于多重认证群组"
	case 0x0d:
		return "卡不在多重认证时间段内"
	case 0x0e:
		return "多重认证模式超级权限认证失败"
	case 0x0f:
		return "多重认证模式远程认证失败"
	case 0x10:
		return "多重认证成功"
	case 0x11:
		return "首卡开门开始"
	case 0x12:
		return "首卡开门结束"
	case 0x13:
		return "常开状态开始"
	case 0x14:
		return "常开状态结束"
	case 0x15:
		return "门锁打开"
	case 0x16:
		return "门锁关闭"
	case 0x17:
		return "开门按钮打开"
	case 0x18:
		return "开门按钮放开"
	case 0x19:
		return "正常开门（门磁）"
	case 0x1a:
		return "正常关门（门磁）"
	case 0x1b:
		return "门异常打开（门磁）"
	case 0x1c:
		return "门打开超时（门磁）"
	case 0x1d:
		return "报警输出打开"
	case 0x1e:
		return "报警输出关闭"
	case 0x1f:
		return "常关状态开始"
	case 0x20:
		return "常关状态结束"
	case 0x21:
		return "多重多重认证需要远程开门"
	case 0x22:
		return "多重认证超级密码认证成功事件"
	case 0x23:
		return "多重认证重复认证事件"
	case 0x24:
		return "多重认证超时"
	case 0x25:
		return "门铃响"
	case 0x26:
		return "指纹比对通过"
	case 0x27:
		return "指纹比对失败"
	case 0x28:
		return "刷卡加指纹认证通过"
	case 0x29:
		return "刷卡加指纹认证失败"
	case 0x2a:
		return "刷卡加指纹认证超时"
	case 0x2b:
		return "刷卡加指纹加密码认证通过"
	case 0x2c:
		return "刷卡加指纹加密码认证失败"
	case 0x2d:
		return "刷卡加指纹加密码认证超时"
	case 0x2e:
		return "指纹加密码认证通过"
	case 0x2f:
		return "指纹加密码认证失败"
	case 0x30:
		return "指纹加密码认证超时"
	case 0x31:
		return "指纹不存在"
	case 0x32:
		return "刷卡平台认证"
	case 0x33:
		return "呼叫中心事件"
	case 0x34:
		return "消防继电器导通触发门常开"
	case 0x35:
		return "消防继电器恢复门恢复正常"
	case 0x36:
		return "人脸加指纹认证通过"
	case 0x37:
		return "人脸加指纹认证失败"
	case 0x38:
		return "人脸加指纹认证超时"
	case 0x39:
		return "人脸加密码认证通过"
	case 0x3a:
		return "人脸加密码认证失败"
	case 0x3b:
		return "人脸加密码认证超时"
	case 0x3c:
		return "人脸加刷卡认证通过"
	case 0x3d:
		return "人脸加刷卡认证失败"
	case 0x3e:
		return "人脸加刷卡认证超时"
	case 0x3f:
		return "人脸加密码加指纹认证通过"
	case 0x40:
		return "人脸加密码加指纹认证失败"
	case 0x41:
		return "人脸加密码加指纹认证超时"
	case 0x42:
		return "人脸加刷卡加指纹认证通过"
	case 0x43:
		return "人脸加刷卡加指纹认证失败"
	case 0x44:
		return "人脸加刷卡加指纹认证超时"
	case 0x45:
		return "工号加指纹认证通过"
	case 0x46:
		return "工号加指纹认证失败"
	case 0x47:
		return "工号加指纹认证超时"
	case 0x48:
		return "工号加指纹加密码认证通过"
	case 0x49:
		return "工号加指纹加密码认证失败"
	case 0x4a:
		return "工号加指纹加密码认证超时"
	case 0x4b:
		return "人脸认证通过"
	case 0x4c:
		return "人脸认证失败"
	case 0x4d:
		return "工号加人脸认证通过"
	case 0x4e:
		return "工号加人脸认证失败"
	case 0x4f:
		return "工号加人脸认证超时"
	case 0x50:
		return "人脸抓拍失败"
	case 0x51:
		return "首卡授权开始"
	case 0x52:
		return "首卡授权结束"
	case 0x53:
		return "门锁输入短路报警"
	case 0x54:
		return "门锁输入断路报警"
	case 0x55:
		return "门锁输入异常报警"
	case 0x56:
		return "门磁输入短路报警"
	case 0x57:
		return "门磁输入断路报警"
	case 0x58:
		return "门磁输入异常报警"
	case 0x59:
		return "开门按钮输入短路报警"
	case 0x5a:
		return "开门按钮输入断路报警"
	case 0x5b:
		return "开门按钮输入异常报警"
	case 0x5c:
		return "门锁异常打开"
	case 0x5d:
		return "门锁打开超时"
	case 0x5e:
		return "首卡未授权开门失败"
	case 0x5f:
		return "呼梯继电器断开"
	case 0x60:
		return "呼梯继电器闭合"
	case 0x61:
		return "自动按键继电器断开"
	case 0x62:
		return "自动按键继电器闭合"
	case 0x63:
		return "按键梯控继电器断开"
	case 0x64:
		return "按键梯控继电器闭合"
	case 0x65:
		return "工号加密码认证通过"
	case 0x66:
		return "工号加密码认证失败"
	case 0x67:
		return "工号加密码认证超时"
	case 0x68:
		return "真人检测失败"
	case 0x69:
		return "人证比对通过"
	case 0x70:
		return "人证比对失败"
	case 0x71:
		return "非授权名单事件"
	case 0x72:
		return "合法短信"
	case 0x73:
		return "非法短信"
	case 0x74:
		return "MAC侦测"
	case 0x75:
		return "门状态常闭或休眠状态认证失败"
	case 0x76:
		return "认证计划休眠模式认证失败"
	case 0x77:
		return "卡加密校验失败"
	case 0x78:
		return "反潜回服务器应答失败"
	case 0x85:
		return "尾随通行"
	case 0x86:
		return "反向闯入"
	case 0x87:
		return "外力冲撞"
	default:
		// 尝试从Part2中获取描述
		if desc := getEventMinorDescPart2(minorType); desc != "未知事件次类型" {
			return desc
		}
		
		// 处理自定义事件类型 (0x500-0x53f)
		if minorType >= 0x500 && minorType <= 0x53f {
			customNum := minorType - 0x500 + 1
			return "门禁自定义事件" + string(rune('0'+customNum/10)) + string(rune('0'+customNum%10))
		}
		return "未知事件次类型"
	}
}