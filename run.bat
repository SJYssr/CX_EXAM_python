@echo off
chcp 65001 >nul
cls
echo ====================================
echo    超星考试客户端工具 (C# 版本)
echo ====================================
echo.

REM 检查.NET环境
echo [1/3] 检查运行环境...
dotnet --version >nul 2>&1
if %errorlevel% neq 0 (
    echo 错误: 未找到 .NET Runtime，请安装 .NET 6.0 或更高版本
    echo 下载地址: https://dotnet.microsoft.com/download
    echo.
    pause
    exit /b 1
)
echo ✓ .NET 环境正常

REM 检查可执行文件
echo [2/3] 检查程序文件...
if not exist "bin\Release\net6.0-windows\CxExamClient.exe" (
    echo ⚠ 可执行文件不存在，开始编译...
    echo.
    call .\build.bat
    if %errorlevel% neq 0 (
        echo ❌ 编译失败
        pause
        exit /b 1
    )
    echo ✓ 编译完成
) else (
    echo ✓ 程序文件存在
)

REM 启动程序
echo [3/3] 启动程序...
echo.
echo 正在启动超星考试客户端工具...
echo ====================================
echo.

"bin\Release\net6.0-windows\CxExamClient.exe"

REM 检查退出状态
if %errorlevel% neq 0 (
    echo.
    echo ❌ 程序异常退出，错误代码: %errorlevel%
    echo 请检查错误日志或重新编译程序
) else (
    echo.
    echo ✓ 程序正常退出
)

echo.
echo 按任意键关闭窗口...
pause >nul 