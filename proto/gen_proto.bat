@echo off
setlocal

:: 检查是否输入了文件名
if "%~1"=="" (
    echo Usage: gen_proto.bat your_file.proto
    pause
    exit /b
)

:: 执行编译命令
:: --go_out: 生成基础结构体
protoc --go_out=. %1

echo [Success] Proto file %1 has been compiled.
