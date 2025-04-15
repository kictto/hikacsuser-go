# 报警描述包 (alarmdesc)

这个包提供了海康威视门禁设备报警信息的描述转换功能，可以将报警代码转换为对应的中文描述。

## 功能

- 将报警主类型代码转换为中文描述
- 将报警次类型代码转换为中文描述
- 支持报警、异常、操作和事件四种主类型
- 支持各种次类型的详细描述

## 使用方法

```go
import "github.com/your-username/hikacsuser-go/pkg/alarmdesc"

// 获取报警主类型描述
majorTypeDesc := alarmdesc.GetAlarmMajorTypeDesc(alarmdesc.ALARM_MAJOR_ALARM) // 返回 "报警"

// 获取报警次类型描述
minorTypeDesc := alarmdesc.GetAlarmMinorTypeDesc(alarmdesc.ALARM_MAJOR_ALARM, alarmdesc.ALARM_MINOR_ZONE_SHORT_CIRCUIT) // 返回 "防区短路报警"

// 也可以直接使用次类型描述函数
alarmDesc := alarmdesc.GetAlarmMinorDesc(alarmdesc.ALARM_MINOR_ZONE_SHORT_CIRCUIT) // 返回 "防区短路报警"
exceptionDesc := alarmdesc.GetExceptionMinorDesc(alarmdesc.EXCEPTION_MINOR_NETWORK_BROKEN) // 返回 "网络断开"
operationDesc := alarmdesc.GetOperationMinorDesc(0x50) // 返回 "本地登陆"
eventDesc := alarmdesc.GetEventMinorDesc(alarmdesc.EVENT_MINOR_LEGAL_CARD_PASS) // 返回 "合法卡认证通过"
```

## 常量

包中定义了所有常用的报警主类型和次类型常量，可以直接使用。例如：

```go
// 报警主类型
ALARM_MAJOR_ALARM     = 0x1 // 报警
ALARM_MAJOR_EXCEPTION = 0x2 // 异常
ALARM_MAJOR_OPERATION = 0x3 // 操作
ALARM_MAJOR_EVENT     = 0x5 // 事件

// 报警次类型 - 报警类型(ALARM_MAJOR_ALARM)
ALARM_MINOR_ZONE_SHORT_CIRCUIT = 0x400 // 防区短路报警
// ... 更多常量
```

## 注意事项

- 对于未知的类型代码，函数会返回默认的描述文本（如"未知类型"、"未知次类型"等）
- 对于自定义类型（0x900-0x93f范围内的代码），会返回相应的自定义描述