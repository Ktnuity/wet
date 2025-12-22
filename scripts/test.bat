@echo off

echo [=] Building

call "%~dp0build.bat"
if errorlevel 1 exit /b 1

echo [=] Testing

"%~dp0..\bin\test.exe"
if errorlevel 1 exit /b 1

