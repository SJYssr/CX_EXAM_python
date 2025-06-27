@echo off
chcp 65001 >nul
cls
echo ====================================================
echo              超星学习通辅助工具
echo              真正有效版 - 控制台版本
echo ====================================================
echo.
echo 🎯 特性：
echo   ✅ 控制台界面，中文显示100%%正常
echo   ✅ 剪贴板功能完全正常
echo   ✅ 简单易用，绝对有效
echo   ✅ 无GUI问题，稳定可靠
echo.
echo 📖 使用方法：
echo   1. 输入关键词搜索题目
echo   2. 结果自动复制到剪贴板
echo   3. 在考试页面按 Ctrl+V 粘贴
echo   4. 输入 'test' 测试功能
echo   5. 输入 'exit' 退出程序
echo.
echo 正在启动程序...
echo.

if exist "chaoxing_working_console.exe" (
    chaoxing_working_console.exe
) else (
    echo ❌ 错误：找不到 chaoxing_working_console.exe
    echo 请确保程序文件存在
    pause
)

echo.
echo 程序已退出
pause 