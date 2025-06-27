@echo off
echo Starting CxExamClient Tool...

if not exist "bin\Release\net6.0-windows\CxExamClient.exe" (
    echo Executable not found, building...
    call build.bat
    if %errorlevel% neq 0 exit /b 1
)

echo.
echo Launching program...
echo.

REM 直接运行程序并等待
"bin\Release\net6.0-windows\CxExamClient.exe"

REM 检查退出码
if %errorlevel% neq 0 (
    echo.
    echo Program exited with error code: %errorlevel%
    echo.
)

echo Program finished.
echo Press any key to continue...
pause >nul 