@echo off
chcp 65001 >nul
echo ====================================================
echo              超星学习通辅助工具
echo                完美解决方案 v6.0
echo ====================================================
echo.
echo 🎯 特性：
echo   ✅ 完美解决中文显示问题
echo   ✅ 剪贴板功能100%正常
echo   ✅ Fyne GUI界面美观现代
echo   ✅ UTF-8编码完全支持
echo   ✅ 智能文件编码检测
echo.
echo 正在启动程序...
echo.

if exist "chaoxing_perfect_v6.exe" (
    start "" "chaoxing_perfect_v6.exe"
    echo ✅ 程序已启动！
) else (
    echo ❌ 错误：找不到 chaoxing_perfect_v6.exe
    echo 请确保程序文件存在
    pause
)

echo.
echo 💡 使用说明：
echo   1. 在搜索框输入关键词
echo   2. 点击"本地搜索"查找题目
echo   3. 点击"复制到剪贴板"复制结果
echo   4. 在考试页面按 Ctrl+V 粘贴
echo.
echo 如果需要帮助，请查看 README.md
pause 