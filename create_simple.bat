@echo off
chcp 65001 >nul
echo Creating Simple Version...

echo Compiling simple GUI version...
go build -ldflags="-s -w -H windowsgui" -o cx_exam_simple.exe -tags=simple .

if exist cx_exam_simple.exe (
    echo Simple version created successfully!
    echo Starting simple version...
    start cx_exam_simple.exe
) else (
    echo Failed to create simple version
    echo Trying console version...
    go build -o cx_exam_console_temp.exe console_version.go
    if exist cx_exam_console_temp.exe (
        start cx_exam_console_temp.exe
    )
)

pause 