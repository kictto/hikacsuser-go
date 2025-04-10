@echo off
REM 这个脚本帮助设置动态库路径，使依赖这个项目的其他项目不需要拷贝动态库
REM 用法: call run_helper.bat

REM 获取脚本所在目录的绝对路径
set "SCRIPT_DIR=%~dp0"

REM 根据系统架构选择合适的库路径
if "%PROCESSOR_ARCHITECTURE%"=="AMD64" (
    set "LIB_PATH=%SCRIPT_DIR%lib\win64"
) else (
    set "LIB_PATH=%SCRIPT_DIR%lib\win32"
)

REM 设置环境变量
set "PATH=%LIB_PATH%;%PATH%"

echo 已设置PATH环境变量，包含: %LIB_PATH%
echo 现在可以运行依赖海康威视SDK的程序了 