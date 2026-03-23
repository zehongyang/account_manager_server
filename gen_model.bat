@echo off
setlocal enabledelayedexpansion
for /f "tokens=2 delims=:." %%a in ('chcp') do set "_OLD_CP=%%a"
set "_OLD_CP=%_OLD_CP: =%"
chcp 65001 >nul

:: 检查参数
if "%~1"=="" (
    echo [ERROR] 请传入需要生成的表名.
    echo 用法: %0 ^<table_name^> [models_dir]
    echo 示例: %0 users ./models/user
    exit /b 1
)

set TABLE_NAME=%~1
set OUT_DIR=%~2

:: 如果未传入第二个参数，默认使用 ./models
if "%OUT_DIR%"=="" set OUT_DIR=./models

set REVERSE_YML=reverse.yml
set TEMP_YML=reverse_temp.yml

echo [INFO] 开始生成模型...
echo [INFO] 表名: %TABLE_NAME%
echo [INFO] 输出目录: %OUT_DIR%

:: 使用 PowerShell 读取原始 yml，修改 output_dir 并在 language: golang 后面追加 include_tables 来指定单表
powershell -Command "$content = Get-Content '%REVERSE_YML%' -Encoding UTF8; $content = $content -replace 'output_dir:.*', 'output_dir: %OUT_DIR%'; $content = $content -replace 'language:\s*golang', \"language: golang`n    include_tables:`n      - '%TABLE_NAME%'\"; Set-Content '%TEMP_YML%' -Value $content -Encoding UTF8"

if errorlevel 1 (
    echo [ERROR] 生成临时配置文件失败
    exit /b 1
)

:: 执行 reverse 工具
:: 使用 go run 确保无论有没有全局安装 xorm 都可以运行
go run xorm.io/reverse@latest -f %TEMP_YML%

if errorlevel 1 (
    echo [ERROR] 模型生成失败
    :: 为了排查问题，失败时暂不删除临时配置
    exit /b 1
)

powershell -Command "$raw = Join-Path '%OUT_DIR%' 'models.go'; $final = Join-Path '%OUT_DIR%' '%TABLE_NAME%.go'; if (Test-Path $raw) { Move-Item -Path $raw -Destination $final -Force }; if (Test-Path $final) { $lines = Get-Content $final -Encoding UTF8; $seen = @{}; $out = New-Object System.Collections.Generic.List[string]; foreach($line in $lines){ if($line -match '^\s*[A-Za-z_][A-Za-z0-9_]*\s+.*`xorm:\"'){ if($seen.ContainsKey($line)){ continue }; $seen[$line] = $true }; $out.Add($line) }; Set-Content $final -Value $out -Encoding UTF8 }"

:: 清理临时文件
del %TEMP_YML%

echo [SUCCESS] 模型生成成功! 保存路径: %OUT_DIR%
if defined _OLD_CP chcp %_OLD_CP% >nul
