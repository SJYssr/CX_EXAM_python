@echo off
title ChaoxingHelper ASCII Edition

echo.
echo ================================================
echo  ChaoxingHelper ASCII Edition v1.0
echo ================================================
echo.
echo Starting pure English interface program...
echo This version avoids ALL character encoding issues.
echo.

if not exist chaoxing_helper_ascii.exe (
    echo ERROR: chaoxing_helper_ascii.exe not found
    echo Please ensure the file exists in current directory
    pause
    exit
)

echo Program found! Starting...
echo.

start chaoxing_helper_ascii.exe

echo Program launched successfully!
echo.
echo FEATURES:
echo - Pure English interface
echo - Search question database
echo - Copy answers to clipboard
echo - Compatible with ALL Windows versions
echo - No encoding issues
echo.
echo If you see any problems:
echo 1. Close and restart the program
echo 2. Check if tiku.txt file exists
echo 3. Run as administrator
echo.
echo Press any key to close this window...
pause >nul 