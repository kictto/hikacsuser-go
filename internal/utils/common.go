package utils

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
	"unsafe"
)

// IsWindows 检查系统是否是Windows
func IsWindows() bool {
	return runtime.GOOS == "windows"
}

// IsLinux 检查系统是否是Linux
func IsLinux() bool {
	return runtime.GOOS == "linux"
}

// GetSDKPath 根据不同的操作系统获取SDK路径
func GetSDKPath() string {
	execPath, err := os.Executable()
	if err != nil {
		return ""
	}
	rootDir := filepath.Dir(execPath)
	libDir := filepath.Join(rootDir, "lib")

	if IsWindows() {
		if runtime.GOARCH == "amd64" {
			return filepath.Join(libDir, "win64")
		} else {
			return filepath.Join(libDir, "win32")
		}
	} else if IsLinux() {
		if runtime.GOARCH == "amd64" {
			return filepath.Join(libDir, "linux64")
		} else {
			return filepath.Join(libDir, "linux32")
		}
	}
	return ""
}

// GetDLLPath 获取动态库路径
func GetDLLPath() string {
	sdkPath := GetSDKPath()
	if IsWindows() {
		return filepath.Join(sdkPath, "HCNetSDK.dll")
	} else if IsLinux() {
		return filepath.Join(sdkPath, "libhcnetsdk.so")
	}
	return ""
}

// ByteToString 字节数组转字符串，遇到0截断
func ByteToString(bytes []byte) string {
	var i int
	for i = 0; i < len(bytes); i++ {
		if bytes[i] == 0 {
			break
		}
	}
	return string(bytes[:i])
}

// CopyStringToByteArray 复制字符串到字节数组
func CopyStringToByteArray(s string, dst []byte) {
	copy(dst, []byte(s))
	// 确保字符串以0结尾
	if len(s) < len(dst) {
		dst[len(s)] = 0
	}
}

// UTF8ToGBK UTF8编码转换为GBK编码
func UTF8ToGBK(utf8Bytes []byte) ([]byte, error) {
	// 在实际使用中，需要引入第三方库实现UTF8到GBK的转换
	// 或者使用内置的转换库
	// 这里仅作为示例
	return utf8Bytes, nil
}

// UTF8ToGBKStr UTF8编码转换为GBK字符串
func UTF8ToGBKStr(utf8Bytes []byte) (string, error) {
	gbkBytes, err := UTF8ToGBK(utf8Bytes)
	if err != nil {
		return "", err
	}
	return string(gbkBytes), nil
}

// LoadPicture 加载图片文件到字节数组
func LoadPicture(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}

// SavePicture 保存字节数组到图片文件
func SavePicture(data []byte, filePath string) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return ioutil.WriteFile(filePath, data, 0644)
}

// WriteBufferToPointer 将字节数组写入指针
func WriteBufferToPointer(data []byte, ptr unsafe.Pointer) {
	if len(data) == 0 || ptr == nil {
		return
	}

	// 使用反射或者unsafe包转换指针
	// 注意：这在Go中是不安全的操作，仅在需要时使用
	p := (*[1 << 30]byte)(ptr)
	for i, b := range data {
		p[i] = b
	}
}

// HexDump 将字节数组转换为十六进制显示
func HexDump(data []byte) string {
	return hex.Dump(data)
}

// FormatTime 格式化时间
func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// GetCurrentTime 获取当前时间字符串
func GetCurrentTime() string {
	return FormatTime(time.Now())
}

// GetResourcePath 获取资源文件路径
func GetResourcePath(relativePath string) string {
	execPath, err := os.Executable()
	if err != nil {
		return ""
	}
	rootDir := filepath.Dir(execPath)
	return filepath.Join(rootDir, "resources", relativePath)
}

// CreateDirectory 创建目录
func CreateDirectory(path string) error {
	return os.MkdirAll(path, 0755)
}

// RemoveFile 删除文件
func RemoveFile(path string) error {
	return os.Remove(path)
}

// FileExists 检查文件是否存在
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// GetFileSize 获取文件大小
func GetFileSize(path string) (int64, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return fi.Size(), nil
}

// CopyFile 复制文件
func CopyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// BytesCombine 合并多个字节数组
func BytesCombine(pBytes ...[]byte) []byte {
	var buffer bytes.Buffer
	for _, b := range pBytes {
		buffer.Write(b)
	}
	return buffer.Bytes()
}

// Print 打印消息到控制台
func Print(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

// Println 打印消息到控制台并换行
func Println(args ...interface{}) {
	fmt.Println(args...)
}

// Printf 格式化打印消息
func Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

// PrintError 打印错误信息
func PrintError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
}

// ConvertToByteArray 将各种类型转换为固定大小的字节数组
func ConvertToByteArray(value interface{}, size int) []byte {
	result := make([]byte, size)

	switch v := value.(type) {
	case string:
		copy(result, []byte(v))
	case []byte:
		copy(result, v)
	case byte:
		if size > 0 {
			result[0] = v
		}
	case int:
		str := fmt.Sprintf("%d", v)
		copy(result, []byte(str))
	// 可以添加更多类型的处理
	default:
		str := fmt.Sprintf("%v", v)
		copy(result, []byte(str))
	}

	// 确保结尾是0
	if len(result) > 0 {
		for i := len(result) - 1; i >= 0; i-- {
			if result[i] != 0 {
				if i+1 < len(result) {
					result[i+1] = 0
				}
				break
			}
		}
	}

	return result
}

// IsEmpty 检查字符串是否为空
func IsEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

// IsNotEmpty 检查字符串是否非空
func IsNotEmpty(s string) bool {
	return !IsEmpty(s)
}
