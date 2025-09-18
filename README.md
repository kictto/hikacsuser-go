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

## 安装说明

### 通过Go Modules安装（推荐方式）

使用Go Modules安装本项目：

```bash
go get github.com/kictto/hikacsuser-go
```

安装后，您需要确保动态链接库能被正确加载。当使用go get或go mod tidy安装本项目时，动态链接库会一同被下载，但不会自动被添加到系统的搜索路径中。**这是由于Go不处理非Go代码的依赖关系**。

因此，您需要在使用项目前，手动设置环境变量，或在代码中动态设置库路径（详见下文使用说明）。

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
// 直接使用pkg包导出的API
import "github.com/kictto/hikacsuser-go/pkg/acsapi"

func main() {
    // 创建客户端(会自动查找和设置库路径)
    client, err := acsapi.NewACSClient()
    if err != nil {
        panic(err)
    }
    
    // 使用客户端...
    
    // 确保在程序结束时注销和清理资源
    defer func() {
        client.Logout()
        client.Cleanup()
    }()
}
```

**或者**使用底层SDK：

```go
import "github.com/kictto/hikacsuser-go/internal/sdk"

func main() {
    // 获取SDK实例(会自动查找和设置库路径)
    sdk, err := sdk.GetSDKInstance() 
    if err != nil {
        panic(err)
    }
    
    // 使用SDK...
}
```

### 确保动态库正确加载

本项目在SDK初始化时会**自动尝试**查找动态库路径，按以下顺序：
1. 首先查找可执行文件所在目录的lib子目录
2. 然后查找当前工作目录的lib子目录
3. 最后尝试使用源码相对路径查找lib目录

上述自动机制在大多数情况下都能正常工作，但如果库路径不是默认位置，可使用以下方法手动指定：

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

#### 方法三：创建库文件符号链接(仅支持Linux)

创建一个到库文件的链接，放在系统库搜索路径中：

```bash
sudo ln -s /path/to/hikacsuser-go/lib/linux64/libhcnetsdk.so /usr/local/lib/
sudo ln -s /path/to/hikacsuser-go/lib/linux64/libHCCore.so /usr/local/lib/
# 之后更新动态链接库缓存
sudo ldconfig
```

## 故障排除

如果遇到`"cannot open shared object file: No such file or directory"`错误，这表明动态库未被找到。
请检查：

1. 确认库文件是否存在于项目的lib目录中
2. 尝试使用上述方法二手动设置环境变量
3. 如果是Windows，确保系统已安装最新的VC++运行时库

## 许可证

使用本项目需遵守海康威视SDK的许可协议。 