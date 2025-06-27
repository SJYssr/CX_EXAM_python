@echo off
chcp 65001 >nul
echo 正在发布超星考试客户端工具 (C# 版本)...
echo.

dotnet --version >nul 2>&1
if %errorlevel% neq 0 (
    echo 错误: 未找到.NET SDK
    pause
    exit /b 1
)

echo 正在发布为单文件可执行程序...
dotnet publish -c Release -r win-x64 --self-contained true -p:PublishSingleFile=true -p:PublishTrimmed=true
if %errorlevel% neq 0 (
    echo 错误: 发布失败
    pause
    exit /b 1
)

echo.
echo 发布成功！
echo 单文件可执行程序位置: bin\Release\net6.0-windows\win-x64\publish\CxExamClient.exe
echo 文件大小约: 
dir "bin\Release\net6.0-windows\win-x64\publish\CxExamClient.exe" | findstr CxExamClient.exe
echo.
echo 注意: 请确保目标机器安装了.NET 6.0 Runtime或使用--self-contained参数
echo.
pause 