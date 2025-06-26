@echo off
chcp 65001 >nul
echo 🚀 编译超星学习通辅助工具 - Fyne GUI版本...

echo 📦 清理旧文件...
if exist cx_exam_fyne.exe del cx_exam_fyne.exe

echo 🔧 设置编译环境...
set CGO_ENABLED=1
set GOOS=windows
set GOARCH=amd64

echo 📚 下载依赖...
go mod tidy

echo 🔨 开始编译...
go build -ldflags="-s -w -H windowsgui" -o cx_exam_fyne.exe .

if exist cx_exam_fyne.exe (
    echo ✅ 编译成功！
    echo 📁 可执行文件: cx_exam_fyne.exe
    echo 📊 文件大小:
    dir cx_exam_fyne.exe | findstr cx_exam_fyne.exe
    echo.
    echo 🎯 运行程序...
    echo ⚠️  如果出现中文乱码，请检查系统区域设置
    start cx_exam_fyne.exe
) else (
    echo ❌ 编译失败！
    pause
)

pause 