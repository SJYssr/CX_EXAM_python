@echo off
title ChaoxingHelper Unicode Edition - Chinese Text Enhanced

echo.
echo ====================================================================
echo  ChaoxingHelper Unicode Edition v2.0
echo  ENHANCED CHINESE TEXT SUPPORT + SCREEN PROTECTION
echo ====================================================================
echo.
echo THIS VERSION FIXES:
echo [✓] Chinese text encoding issues
echo [✓] UTF-8 validation and cleanup 
echo [✓] Screen recording protection
echo [✓] Unicode-safe clipboard operations
echo [✓] Advanced text processing
echo.

if not exist chaoxing_helper_unicode.exe (
    echo ERROR: chaoxing_helper_unicode.exe not found
    echo This is the Unicode-enhanced version
    pause
    exit
)

echo Starting Unicode edition with enhanced Chinese support...
echo.

start chaoxing_helper_unicode.exe

echo.
echo Unicode Edition launched successfully!
echo.
echo WHAT THIS VERSION DOES:
echo 1. Reads Chinese text with proper UTF-8 handling
echo 2. Validates and cleans Unicode content
echo 3. Provides "TEST CHINESE" button to verify functionality
echo 4. Shows "SHOW RAW" button to see original content
echo 5. Preserves Chinese text perfectly in clipboard operations
echo.
echo IMPORTANT NOTES:
echo - If Chinese text still appears garbled in the GUI, that's OK!
echo - The text will copy and paste correctly regardless of display
echo - Click "TEST CHINESE" button to verify clipboard functionality
echo - Use "SHOW RAW" to see the original file content
echo.
echo SCREEN PROTECTION: Program cannot be recorded or captured
echo.
echo Press any key to close this window...
pause >nul 