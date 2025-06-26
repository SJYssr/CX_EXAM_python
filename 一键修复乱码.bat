@echo off
:: 设置UTF-8编码
chcp 65001 >nul
title 一键修复字符编码问题

echo.
echo ====================================================
echo                一键修复字符编码问题
echo ====================================================
echo.
echo 正在检测系统编码...

:: 检测是否为中文系统
for /f "tokens=*" %%i in ('systeminfo ^| findstr /C:"系统区域设置"') do (
    echo 检测到中文系统
    goto :chinese_system
)

for /f "tokens=*" %%i in ('systeminfo ^| findstr /C:"System Locale"') do (
    echo 检测到英文系统，需要设置中文支持
    goto :english_system
)

:chinese_system
echo 您的系统已经是中文系统，直接启动简化版程序...
goto :start_simple

:english_system
echo.
echo 检测到您的系统不是中文系统，请按以下步骤操作：
echo.
echo 1. 按 Win + R 键
echo 2. 输入: intl.cpl
echo 3. 按回车打开区域设置
echo 4. 点击"管理"/"Administrative"选项卡
echo 5. 点击"更改系统区域设置"/"Change system locale"
echo 6. 选择"中文(简体，中国)"/"Chinese (Simplified, China)"
echo 7. 勾选"Beta版：使用Unicode UTF-8提供全球语言支持"
echo 8. 重启计算机
echo.
echo 现在为您启动英文界面的简化版程序...
echo.

:start_simple
if exist cx_exam_simple.exe (
    echo 找到简化版程序，正在启动...
    start cx_exam_simple.exe
    echo.
    echo 简化版程序已启动！
    echo 这个版本使用英文界面，避免了字符编码问题。
    echo 功能完全相同：搜索、查看题库、自动输入等。
) else (
    echo 错误：找不到简化版程序 cx_exam_simple.exe
    echo 请确保文件存在于当前目录
    goto :manual_fix
)

echo.
echo 如果简化版程序也无法正常显示，请尝试以下操作：
echo 1. 右键点击程序图标
echo 2. 选择"属性" -> "兼容性"
echo 3. 勾选"以兼容模式运行这个程序"
echo 4. 选择 Windows 8 或 Windows 10
echo.
goto :end

:manual_fix
echo.
echo 手动修复方法：
echo 1. Win + I 打开设置
echo 2. 时间和语言 -> 区域
echo 3. 国家或地区设置为"中国"
echo 4. 区域格式设置为"中文(简体，中国)"
echo 5. 语言 -> Windows显示语言设置为"中文(简体，中国)"
echo.

:end
echo 按任意键关闭此窗口...
pause >nul 