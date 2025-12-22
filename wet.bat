@echo off
setlocal enabledelayedexpansion

:: Set default flags
set DO_RUN=0
set DO_BUILD=0
set DO_TEST=0

:: Parse arguments
set PASS_ARGS=
:parse_loop
if "%~1"=="" goto end_parse
if "%~1"=="--build" (
    set DO_BUILD=1
    shift
    goto parse_loop
)
if "%~1"=="--test" (
    set DO_TEST=1
    shift
    goto parse_loop
)
if "%~1"=="--run" (
    set DO_RUN=1
    shift
    goto parse_loop
)
if "%~1"=="--all" (
    set DO_BUILD=1
    set DO_TEST=1
    set DO_RUN=1
    shift
    goto parse_loop
)
if "%~1"=="--" (
    shift
    goto collect_args
)
shift
goto parse_loop

:: Collect remaining arguments after --
:collect_args
if "%~1"=="" goto end_parse
set PASS_ARGS=!PASS_ARGS! %1
shift
goto collect_args

:end_parse

:: Building step
if %DO_BUILD%==1 (
    call scripts\build.bat
    if errorlevel 1 exit /b 1
)

:: Testing step
if %DO_TEST%==1 (
    call scripts\test.bat
    if errorlevel 1 exit /b 1
)

:: Running step
if %DO_RUN%==1 (
    call scripts\run.bat %PASS_ARGS%
    if errorlevel 1 exit /b 1
)

endlocal
