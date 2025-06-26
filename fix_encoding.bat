@echo off
title 字符编码修复工具
chcp 65001 >nul

echo ================================================
echo           字符编码修复工具
echo ================================================
echo.

echo 正在设置字符编码...
chcp 65001 >nul

echo 检查系统信息...
systeminfo | findstr /C:"系统区域设置" >nul
if errorlevel 1 (
    echo 检测到非中文系统，正在应用兼容性设置...
) else (
    echo 检测到中文系统
)

echo.
echo 请按照以下步骤操作：
echo.
echo 1. 按下 Win + R 键
echo 2. 输入: intl.cpl
echo 3. 按回车键打开区域设置
echo 4. 点击"管理"选项卡
echo 5. 点击"更改系统区域设置"
echo 6. 选择"中文(简体，中国)"
echo 7. 勾选"Beta版：使用Unicode UTF-8提供全球语言支持"
echo 8. 重启计算机
echo.

echo 或者，直接运行简化版程序：
if exist cx_exam_simple.exe (
    echo 找到简化版程序，正在启动...
    start cx_exam_simple.exe
) else (
    echo 正在创建简化版程序...
    call create_simple.bat
)

echo.
echo 按任意键继续...
pause >nul 