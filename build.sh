#!/bin/bash
echo "编译海康威视门禁系统SDK封装程序"

# 设置Go环境变量
export GOOS=linux
export CGO_ENABLED=1

# 如果不存在bin目录，创建它
if [ ! -d bin ]; then
    mkdir bin
fi

# 编译SDK库
cd internal/sdk
go build -o ../../bin/sdk.a -buildmode=archive
if [ $? -ne 0 ]; then
    echo "SDK编译失败"
    exit 1
fi
cd ../..

# 编译主程序
cd cmd/acsdemo
go build -o ../../bin/acsdemo -ldflags "-linkmode external"
if [ $? -ne 0 ]; then
    echo "主程序编译失败"
    exit 1
fi
cd ../..

# 复制动态库和资源文件
cp lib/linux64/*.so bin/
if [ ! -d bin/resources ]; then
    mkdir bin/resources
fi
if [ ! -d bin/resources/pic ]; then
    mkdir bin/resources/pic
fi
cp resources/pic/*.jpg bin/resources/pic/

# 可选：创建符号链接到lib目录
# ln -s ../lib bin/lib

echo "编译完成，程序已保存到bin目录"

# 设置LD_LIBRARY_PATH以便程序能找到动态库
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:$(pwd)/lib/linux64

# 询问是否运行程序
read -p "是否运行程序(y/n)? " run
if [ "$run" = "y" ]; then
    cd bin
    ./acsdemo
fi 