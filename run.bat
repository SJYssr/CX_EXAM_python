@echo off
chcp 65001 >nul
echo 启动超星学习通辅助工具...

if exist chaoxing-helper.exe (
    echo ✅ 找到程序文件，正在启动...
    chaoxing-helper.exe
) else (
    echo ❌ 未找到程序文件 chaoxing-helper.exe
    echo 请先运行 build.bat 编译程序
    pause
) 