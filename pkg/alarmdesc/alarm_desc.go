package alarmdesc

// GetAlarmMajorTypeDesc 获取报警主类型描述
// 根据报警主类型代码返回对应的中文描述
import "github.com/clockworkchen/hikacsuser-go/internal/utils"

func GetAlarmMajorTypeDesc(majorType int) string {
	return utils.GetAlarmMajorTypeDesc(majorType)
}

// GetAlarmMinorTypeDesc 获取报警次类型描述
// 根据报警主类型和次类型代码返回对应的中文描述
func GetAlarmMinorTypeDesc(majorType, minorType int) string {
	return utils.GetAlarmMinorTypeDesc(majorType, minorType)
}
