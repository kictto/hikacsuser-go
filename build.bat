@echo off

REM 使用echo命令输出中文
echo [信息] 编译海康威视门禁系统SDK封装程序

REM 设置Go环境变量
set GOOS=windows
set CGO_ENABLED=1

REM 设置CGO链接器标志
set CGO_LDFLAGS=-LD:/workspace/hikacsuser-go/lib/win64 -lHCNetSDK -lHCCore -lPlayCtrl
set CGO_CFLAGS=-ID:/workspace/hikacsuser-go/lib/win64

REM 如果不存在bin目录，创建它
if not exist bin mkdir bin

REM 编译SDK库
cd internal\sdk
go build -o ..\..\bin\sdk.a -buildmode=archive
if %ERRORLEVEL% NEQ 0 (
    echo [错误] SDK编译失败
    exit /b %ERRORLEVEL%
)
cd ..\..\bin

REM 编译主程序
cd cmd\acsdemo
go build -o ..\..\bin\acsdemo.exe
if %ERRORLEVEL% NEQ 0 (
    echo 主程序编译失败
    exit /b %ERRORLEVEL%
)
cd ..\..

REM 复制动态库和资源文件
xcopy /Y lib\win64\*.dll bin\
xcopy /Y lib\win64\HCNetSDKCom\*.dll bin\HCNetSDKCom\

REM 创建并更新配置文件，禁用端口映射
echo ^<?xml version="1.0" encoding="GB2312"?^> > bin\HCNetSDK_Log_Switch.xml
echo ^<SdkLocal^> >> bin\HCNetSDK_Log_Switch.xml
echo     ^<SdkLog^> >> bin\HCNetSDK_Log_Switch.xml
echo         ^<logLevel^>3^</logLevel^> >> bin\HCNetSDK_Log_Switch.xml
echo         ^<logDirectory^>./sdklog/^</logDirectory^> >> bin\HCNetSDK_Log_Switch.xml
echo         ^<autoDelete^>true^</autoDelete^> >> bin\HCNetSDK_Log_Switch.xml
echo     ^</SdkLog^> >> bin\HCNetSDK_Log_Switch.xml
echo     ^<!-- 禁用端口映射 --^> >> bin\HCNetSDK_Log_Switch.xml
echo     ^<Net^> >> bin\HCNetSDK_Log_Switch.xml
echo         ^<TransUse^>false^</TransUse^> >> bin\HCNetSDK_Log_Switch.xml
echo     ^</Net^> >> bin\HCNetSDK_Log_Switch.xml
echo ^</SdkLocal^> >> bin\HCNetSDK_Log_Switch.xml

if not exist bin\resources mkdir bin\resources
if not exist bin\resources\pic mkdir bin\resources\pic
xcopy /Y resources\pic\*.jpg bin\resources\pic\

echo 编译完成，程序已保存到bin目录

REM 询问是否运行程序
set /p run=是否运行程序(y/n)?
if /i "%run%"=="y" (
    cd bin
    acsdemo.exe
)