package hclib

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// GetModuleLibPath 尝试查找当前模块在GOPATH/pkg/mod中的lib目录绝对路径
// 返回适合当前操作系统的库路径 (e.g., .../lib/win64 or .../lib/linux64)
func GetModuleLibPath() (string, error) {
	// 获取调用者（应该是sdk包里的某个文件）的文件路径
	_, callerFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("无法获取调用者信息")
	}

	// 预期结构: GOPATH/pkg/mod/github.com/clockworkchen/hikacsuser-go@vX.Y.Z/internal/sdk/somefile.go
	// 我们需要向上找到模块根目录下的 lib 目录
	// 从 sdk 包向上两级是模块根目录
	moduleRoot := filepath.Dir(filepath.Dir(callerFile))
	libDir := filepath.Join(moduleRoot, "lib")

	var osArchDir string
	if runtime.GOOS == "windows" {
		if runtime.GOARCH == "amd64" {
			osArchDir = "win64"
		} else {
			osArchDir = "win32"
		}
	} else if runtime.GOOS == "linux" {
		if runtime.GOARCH == "amd64" {
			osArchDir = "linux64"
		} else {
			osArchDir = "linux32"
		}
	} else {
		return "", fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
	}

	finalLibPath := filepath.Join(libDir, osArchDir)

	// 验证路径是否存在
	if _, err := os.Stat(finalLibPath); os.IsNotExist(err) {
		// 如果基于 runtime.Caller 的路径不存在，尝试备用方法（可能在 GOPATH 外？）
		// 这里可以添加更多复杂的查找逻辑，但 runtime.Caller 在模块缓存中通常是可靠的
		return "", fmt.Errorf("无法在模块缓存中找到库路径: %s (%v)", finalLibPath, err)
	}

	return finalLibPath, nil
}
