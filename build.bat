@echo off
echo Building CxExamClient (C# Version)...
echo.

REM Check .NET environment
dotnet --version >nul 2>&1
if %errorlevel% neq 0 (
    echo Error: .NET SDK not found, please install .NET 6.0 or higher
    echo Download: https://dotnet.microsoft.com/download
    pause
    exit /b 1
)

echo Restoring NuGet packages...
dotnet restore
if %errorlevel% neq 0 (
    echo Error: NuGet package restore failed
    pause
    exit /b 1
)

echo.
echo Building project...
dotnet build --configuration Release
if %errorlevel% neq 0 (
    echo Error: Project build failed
    pause
    exit /b 1
)

echo.
echo Build successful!
echo Executable: bin\Release\net6.0-windows\CxExamClient.exe
echo.
echo Press any key to continue...
pause >nul 