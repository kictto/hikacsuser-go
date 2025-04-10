# 海康威视门禁系统API包

## 简介

本包提供了海康威视门禁系统的公共API接口，将内部实现封装为可供其他项目依赖使用的公共方法。支持批量和单个数据操作，解决了`Use of the internal package is not allowed`的问题。

## 功能特性

- 提供完整的门禁系统操作API
- 支持单个和批量数据操作
- 简化了设备登录和资源管理
- 提供了错误处理机制
- 完全兼容原有的功能实现

## 使用方法

### 导入包

```go
import "github.com/clockworkchen/hikacsuser-go/pkg/acsapi"
```

### 创建客户端并登录

```go
// 创建门禁系统客户端
client, err := acsapi.NewACSClient()
if err != nil {
    fmt.Printf("创建客户端失败: %v\n", err)
    return
}

// 登录设备
loginInfo, deviceInfo, err := client.Login("192.168.1.5", 8000, "admin", "password")
if err != nil {
    fmt.Printf("登录设备失败: %v\n", err)
    return
}

// 可以使用返回的登录信息和设备信息
fmt.Printf("设备序列号: %s\n", deviceInfo.SerialNumber)
fmt.Printf("设备类型: %d\n", deviceInfo.DeviceType)

// 确保在程序结束时注销设备并清理资源
defer func() {
    client.Logout()
    client.Cleanup()
}()
```

### 用户管理

```go
// 添加单个用户
err = client.AddUser("12345")

// 批量添加用户
errors := client.AddUsers([]string{"23456", "34567"})

// 查询用户信息
err = client.SearchUser()

// 删除单个用户
err = client.DeleteUser("12345")

// 批量删除用户
errors = client.DeleteUsers([]string{"23456", "34567"})

// 删除所有用户
err = client.DeleteAllUsers()
```

### 卡片管理

```go
// 添加单个卡片
err = client.AddCard("12345", "12345")

// 批量添加卡片
errors = client.AddCards(
    []string{"23456", "34567"}, 
    []string{"23456", "34567"},
)

// 查询卡片信息
err = client.SearchCard("12345")

// 删除卡片
err = client.DeleteCard("12345")

// 批量删除卡片
errors = client.DeleteCards([]string{"23456", "34567"})
```

### 人脸管理

```go
// 通过二进制方式添加人脸
err = client.AddFaceByBinary("12345")

// 通过二进制方式批量添加人脸
errors = client.AddFacesByBinary([]string{"23456", "34567"})

// 通过URL方式添加人脸
err = client.AddFaceByUrl("12345")

// 通过URL方式批量添加人脸
errors = client.AddFacesByUrl([]string{"23456", "34567"})

// 查询人脸信息
err = client.SearchFace("12345")

// 删除人脸
err = client.DeleteFace("12345")

// 批量删除人脸
errors = client.DeleteFaces([]string{"23456", "34567"})

// 采集人脸
err = client.CaptureFace()
```

### 门禁管理

```go
// 获取门禁参数
err = client.GetACSConfig()

// 获取门禁状态
err = client.GetACSStatus()

// 远程控门
err = client.RemoteControlGate()

// 查询所有门禁历史事件
err = client.SearchAllEvents()

// 设置卡片模板
err = client.SetCardTemplate(2)
```

### 布控管理

```go
// 定义布控回调函数
callback := func(alarmType int, alarmInfo interface{}) error {
    fmt.Printf("收到报警事件，类型: %d\n", alarmType)
    // 处理报警信息
    return nil
}

// 设置布控（不强制替换现有布控）
alarmHandle, err := client.SetupAlarm(callback, false)
if err != nil {
    fmt.Printf("设置布控失败: %v\n", err)
    return
}

// 获取布控状态
isSetup, handle := client.GetAlarmStatus()
fmt.Printf("布控状态: %v, 句柄: %d\n", isSetup, handle)

// 查询报警事件
err = client.SearchAlarmEvent()
if err != nil {
    fmt.Printf("查询报警事件失败: %v\n", err)
}

// 关闭布控
err = client.CloseAlarm()
if err != nil {
    fmt.Printf("关闭布控失败: %v\n", err)
}
```

## 错误处理

单个操作方法返回`error`类型，批量操作方法返回`[]error`类型，可以通过检查返回值来判断操作是否成功：

```go
// 单个操作错误处理
err = client.AddUser("12345")
if err != nil {
    fmt.Printf("添加用户失败: %v\n", err)
}

// 批量操作错误处理
errors := client.AddUsers([]string{"23456", "34567"})
if len(errors) > 0 {
    for _, err := range errors {
        fmt.Printf("批量添加用户错误: %v\n", err)
    }
}
```

## 完整示例

请参考`pkg/acsapi/example/main.go`文件，其中包含了完整的使用示例。