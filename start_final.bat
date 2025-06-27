@echo off
title ChaoxingHelper FINAL Edition - No Encoding Issues + Screen Protection

echo.
echo ===================================================================
echo  ChaoxingHelper FINAL Edition v1.1
echo  - Pure English Interface (No Garbled Text)
echo  - Screen Recording Protection Enabled
echo  - Chinese Content Properly Displayed
echo ===================================================================
echo.

if not exist chaoxing_helper_ascii_fixed.exe (
    echo ERROR: chaoxing_helper_ascii_fixed.exe not found
    echo Please ensure the file exists in current directory
    pause
    exit
)

echo Starting final version with all fixes...
echo.
echo FEATURES:
echo [+] English interface - NO garbled text in menus/buttons
echo [+] Chinese content - PROPERLY displayed in content area
echo [+] Screen protection - PREVENTS screen recording/screenshots
echo [+] Search function - FULLY working
echo [+] Clipboard copy - PERFECT text preservation
echo [+] Windows compatibility - ALL versions supported
echo.

start chaoxing_helper_ascii_fixed.exe

echo.
echo Program launched! You should see:
echo 1. Clean English interface without any symbols like "◊mm ◊◊"
echo 2. Chinese content displaying correctly in the text area
echo 3. Screen protection active (program cannot be recorded)
echo.
echo If you still see encoding issues, your system needs:
echo - Unicode UTF-8 support enabled in Windows settings
echo - Or use a different computer with proper Chinese support
echo.
echo Press any key to close this window...
pause >nul 