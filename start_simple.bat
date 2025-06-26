@echo off
title Chaoxing Helper Tool

echo.
echo ===============================================
echo  Chaoxing Helper Tool - Fyne GUI Version
echo ===============================================
echo.

if not exist cx_exam_fyne_fixed.exe (
    echo ERROR: cx_exam_fyne_fixed.exe not found
    echo Please ensure the program is compiled correctly
    pause
    exit
)

echo Starting program...
echo.

start cx_exam_fyne_fixed.exe

echo Program started successfully!
echo.
echo If you see garbled text or encoding issues:
echo 1. Check Windows language settings
echo 2. Set region format to Chinese
echo 3. Set display language to Simplified Chinese
echo 4. Enable Unicode UTF-8 support in Control Panel
echo.
echo Press any key to close this window...
pause >nul 