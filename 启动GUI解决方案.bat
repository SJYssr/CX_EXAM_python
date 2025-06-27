@echo off
chcp 65001 >nul
cls
echo ====================================================
echo              超星学习通辅助工具
echo                GUI解决方案选择
echo ====================================================
echo.
echo 🎯 我为您准备了两个GUI解决方案：
echo.
echo 【1】 Fyne GUI + 字体修复版
echo      - 文件：chaoxing_gui_fixed.exe
echo      - 特点：原生GUI界面，应用字体修复
echo      - 适合：希望使用桌面应用的用户
echo.
echo 【2】 Webview 终极版
echo      - 文件：chaoxing_webview_ultimate.exe  
echo      - 特点：浏览器界面，中文显示100%正常
echo      - 适合：要求中文显示完美的用户
echo.
echo 【3】 控制台版本（备选）
echo      - 文件：chaoxing_working_console.exe
echo      - 特点：简单稳定，中文显示正常
echo      - 适合：不需要GUI的用户
echo.
echo ====================================================

:menu
echo.
set /p choice="请选择版本 (1/2/3)："

if "%choice%"=="1" goto gui_fixed
if "%choice%"=="2" goto webview
if "%choice%"=="3" goto console
echo 无效选择，请重新输入！
goto menu

:gui_fixed
echo.
echo 正在启动 Fyne GUI + 字体修复版...
if exist "chaoxing_gui_fixed.exe" (
    start "" "chaoxing_gui_fixed.exe"
    echo ✅ 程序已启动！
    echo 💡 如果仍有乱码，请尝试Webview版本
) else (
    echo ❌ 找不到 chaoxing_gui_fixed.exe
)
goto end

:webview
echo.
echo 正在启动 Webview 终极版...
if exist "chaoxing_webview_ultimate.exe" (
    start "" "chaoxing_webview_ultimate.exe"
    echo ✅ 程序已启动！
    echo 🌐 浏览器将自动打开，中文显示100%正常！
) else (
    echo ❌ 找不到 chaoxing_webview_ultimate.exe
)
goto end

:console
echo.
echo 正在启动控制台版本...
if exist "chaoxing_working_console.exe" (
    chaoxing_working_console.exe
) else (
    echo ❌ 找不到 chaoxing_working_console.exe
)
goto end

:end
echo.
echo 如需帮助，请查看项目文档
pause 