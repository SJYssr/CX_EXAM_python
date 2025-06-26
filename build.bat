@echo off
chcp 65001 >nul
echo 超星学习通辅助工具 - Go版本编译脚本
echo =====================================

:: 检查Go环境
echo [1/5] 检查Go环境...
go version >nul 2>&1
if errorlevel 1 (
    echo ❌ 错误：未找到Go环境，请先安装Go 1.21+
    echo 下载地址：https://golang.org/dl/
    pause
    exit /b 1
)

for /f "tokens=3" %%a in ('go version') do set GO_VERSION=%%a
echo ✅ Go版本：%GO_VERSION%

:: 设置Go代理（中国用户）
echo [2/5] 设置Go代理...
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GOSUMDB=sum.golang.google.cn
echo ✅ 已设置中国Go代理

:: 初始化模块
echo [3/5] 初始化Go模块...
if not exist go.mod (
    go mod init chaoxing-helper-go
)

:: 下载依赖
echo [4/5] 下载依赖包...
go mod tidy
if errorlevel 1 (
    echo ❌ 错误：依赖下载失败
    echo 请检查网络连接或尝试手动执行：go mod tidy
    pause
    exit /b 1
)
echo ✅ 依赖下载完成

:: 编译程序
echo [5/5] 编译程序...
echo 正在编译，请稍候...

:: 设置编译参数
set CGO_ENABLED=1
set GOOS=windows
set GOARCH=amd64

:: 编译优化参数
set LDFLAGS=-w -s -H windowsgui
set GCFLAGS=-trimpath

go build -ldflags "%LDFLAGS%" -gcflags "%GCFLAGS%" -o chaoxing-helper.exe
if errorlevel 1 (
    echo ❌ 编译失败
    echo 请检查代码是否有语法错误
    pause
    exit /b 1
)

echo ✅ 编译成功！
echo.
echo 📁 生成文件：chaoxing-helper.exe
echo 📊 文件大小：
for %%i in (chaoxing-helper.exe) do echo    %%~zi 字节

echo.
echo 🚀 运行程序：
echo    方式1：双击 chaoxing-helper.exe
echo    方式2：执行 chaoxing-helper.exe
echo    方式3：执行 run.bat

echo.
echo 📋 使用说明：
echo    - F3: 切换窗口大小
echo    - F4: 刷新输入法保护
echo    - ESC: 退出程序
echo    - 右键: 改变透明度

:: 询问是否立即运行
echo.
set /p choice="是否立即运行程序？(Y/n): "
if /i "%choice%"=="n" goto :end
if /i "%choice%"=="no" goto :end

echo.
echo 🎯 启动程序...
start chaoxing-helper.exe

:end
echo.
echo ✨ 编译完成，按任意键退出...
pause >nul 