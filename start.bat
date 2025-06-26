@echo off
chcp 65001 >nul
title 超星学习通辅助工具

echo.
echo ===============================================
echo  🚀 超星学习通辅助工具 - Fyne GUI版
echo ===============================================
echo.
echo 🔧 设置字符编码为UTF-8...
echo 📍 当前目录: %CD%
echo.

if not exist cx_exam_fyne_fixed.exe (
    echo ❌ 错误：找不到 cx_exam_fyne_fixed.exe
    echo 💡 请确保程序已正确编译
    pause
    exit
)

echo ✅ 找到程序文件
echo 🎯 启动程序...
echo.

start cx_exam_fyne_fixed.exe

echo ✅ 程序已启动！
echo 💡 如果程序没有显示或出现乱码，请检查：
echo    1. Windows系统语言设置
echo    2. 区域格式设置为中文
echo    3. 显示语言为简体中文
echo.
echo 🔍 按任意键关闭此窗口...
pause >nul 