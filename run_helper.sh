#!/bin/bash
# 这个脚本帮助设置动态库路径，使依赖这个项目的其他项目不需要拷贝动态库
# 用法: source ./run_helper.sh

# 获取脚本所在目录的绝对路径
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

# 根据系统架构选择合适的库路径
if [ "$(uname -m)" == "x86_64" ]; then
    LIB_PATH="$SCRIPT_DIR/lib/linux64"
else
    LIB_PATH="$SCRIPT_DIR/lib/linux32"
fi

# 设置环境变量
export LD_LIBRARY_PATH="$LIB_PATH:$LD_LIBRARY_PATH"

echo "已设置LD_LIBRARY_PATH环境变量，包含: $LIB_PATH"
echo "现在可以运行依赖海康威视SDK的程序了" 