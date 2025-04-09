package utils

// getEventMinorDesc 获取事件类型的次类型描述（续）
func getEventMinorDescPart2(minorType int) string {
	switch minorType {
	case 0x88:
		return "翻越"
	case 0x89:
		return "通行超时"
	case 0x8a:
		return "误闯报警"
	case 0x8b:
		return "闸机自由通行时未认证通过"
	case 0x8c:
		return "摆臂被阻挡"
	case 0x8d:
		return "摆臂阻挡消除"
	case 0x8e:
		return "设备升级本地人脸建模失败"
	case 0x8f:
		return "逗留事件"
	case 0x97:
		return "密码不匹配"
	case 0x98:
		return "工号不存在"
	case 0x99:
		return "组合认证通过"
	case 0x9a:
		return "组合认证超时"
	case 0x9b:
		return "认证方式不匹配"
	case 0x609:
		return "智能锁多重开门"
	default:
		// 处理自定义事件类型 (0x500-0x53f)
		if minorType >= 0x500 && minorType <= 0x53f {
			customNum := minorType - 0x500 + 1
			return "门禁自定义事件" + string(rune('0'+customNum/10)) + string(rune('0'+customNum%10))
		}
		return "未知事件次类型"
	}
}

// 完善getEventMinorDesc函数，调用getEventMinorDescPart2
func init() {
	// 这个init函数会在包被导入时自动执行
	// 用于确保getEventMinorDesc函数能够正确调用getEventMinorDescPart2
}