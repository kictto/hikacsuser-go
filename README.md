# 海康威视门禁SDK Go封装

本项目是海康威视门禁SDK的Go语言封装，支持Windows和Linux平台。

## 项目特点

- 支持Windows和Linux平台
- 不需要在可执行文件目录下复制动态链接库
- 使用CGO与海康威视SDK原生库交互

## 目录结构

```
hikacsuser-go/
├── bin/                # 编译后的可执行文件
├── cmd/                # 命令行工具
│   └── acsdemo/        # 示例程序
├── internal/           # 内部包
│   ├── models/         # 业务模型
│   └── sdk/            # SDK封装
├── lib/                # 动态库文件
│   ├── win64/          # Windows 64位库
│   ├── win32/          # Windows 32位库
│   ├── linux64/        # Linux 64位库
│   └── linux32/        # Linux 32位库
└── resources/          # 资源文件
```

## 编译说明

### Windows

```bash
# 编译项目
./build.bat

# 运行Demo
cd bin
./acsdemo.exe
```

### Linux

```bash
# 编译项目
./build.sh

# 运行Demo
cd bin
./acsdemo
```

## 使用说明

### 在其他项目中引用该库

Go中引用方式：

```go
import "github.com/yourname/hikacsuser-go/internal/sdk"

func main() {
    sdk, err := sdk.GetSDKInstance()
    if err != nil {
        panic(err)
    }
    
    // 使用SDK...
}
```

### 确保动态库正确加载

本项目设计时避免了需要复制动态库到可执行文件目录的情况，有以下几种方法可以确保动态库被正确加载：

#### 方法一：使用提供的帮助脚本

Windows:
```cmd
call path\to\hikacsuser-go\run_helper.bat
```

Linux:
```bash
source path/to/hikacsuser-go/run_helper.sh
```

#### 方法二：手动设置环境变量

Windows:
```cmd
set PATH=path\to\hikacsuser-go\lib\win64;%PATH%
```

Linux:
```bash
export LD_LIBRARY_PATH=path/to/hikacsuser-go/lib/linux64:$LD_LIBRARY_PATH
```

#### 方法三：在项目内部设置库路径

在程序启动时，调用SDK中的`setLibraryPath`方法设置库路径。此方法已在SDK内部自动调用，无需手动操作。

## 许可证

使用本项目需遵守海康威视SDK的许可协议。 