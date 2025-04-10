package sdk

import (
	"bytes"
)

// BytesToString 将字节数组转换为字符串，遇到0字节时截断
// 用于处理C语言风格的字符串（以null结尾）
func BytesToString(b []byte) string {
	n := bytes.IndexByte(b, 0)
	if n < 0 {
		n = len(b)
	}
	return string(b[:n])
}
