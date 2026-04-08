@echo off
REM Lint script for go-react-ink (Windows)
REM Usage: scripts\lint.bat [--fix]

setlocal EnableDelayedExpansion

set RED=[91m
set GREEN=[92m
set YELLOW=[93m
set NC=[0m

set FIX_MODE=0
if "%1"=="--fix" set FIX_MODE=1

echo === Go-React-Ink Lint Checker ===
echo.

REM Check Go installation
echo %YELLOW%Checking Go installation...%NC%
where go >nul 2>&1
if errorlevel 1 (
    echo %RED%ERROR: Go is not installed%NC%
    echo Please install Go from https://golang.org/dl/
    exit /b 1
)
for /f "tokens=*" %%i in ('go version') do set GO_VERSION=%%i
echo %GREEN%✓ !GO_VERSION!%NC%
echo.

REM Check gofmt
echo %YELLOW%Checking gofmt...%NC%
gofmt -s -d . >nul 2>&1
if errorlevel 1 (
    echo %RED%ERROR: gofmt check failed%NC%
    exit /b 1
)

REM Get gofmt output
for /f "tokens=*" %%i in ('gofmt -s -d .') do set GOFMT_OUTPUT=%%i
if defined GOFMT_OUTPUT (
    if !FIX_MODE!==1 (
        echo %YELLOW%Fixing formatting issues...%NC%
        gofmt -s -w .
        echo %GREEN%✓ Fixed formatting issues%NC%
    ) else (
        echo %RED%ERROR: gofmt found issues:%NC%
        gofmt -s -d .
        echo.
        echo %YELLOW%Run with --fix to auto-fix these issues%NC%
        exit /b 1
    )
) else (
    echo %GREEN%✓ gofmt: No formatting issues%NC%
)
echo.

REM Check go vet
echo %YELLOW%Running go vet...%NC%
go vet ./... 2>&1
if errorlevel 1 (
    echo %RED%ERROR: go vet found issues%NC%
    exit /b 1
)
echo %GREEN%✓ go vet: No issues%NC%
echo.

REM Check golint installation
echo %YELLOW%Checking golint installation...%NC%
where golint >nul 2>&1
if errorlevel 1 (
    echo %YELLOW%golint not found, installing...%NC%
    go install golang.org/x/lint/golint@latest

    REM Check if GOPATH/bin is in PATH
    for /f "tokens=*" %%i in ('go env GOPATH') do set GOPATH=%%i
    set GOPATH_BIN=!GOPATH!\bin

    REM Try to run golint from GOPATH/bin
    if exist "!GOPATH_BIN!\golint.exe" (
        set "PATH=!PATH!;!GOPATH_BIN!"
    ) else (
        echo %RED%ERROR: Failed to install golint%NC%
        echo Please manually install: go install golang.org/x/lint/golint@latest
        exit /b 1
    )
    echo %GREEN%✓ golint installed%NC%
) else (
    echo %GREEN%✓ golint: installed%NC%
)
echo.

REM Run golint
echo %YELLOW%Running golint...%NC%
golint ./... 2>&1 | findstr /V "KeyF" | findstr /V "KeyCtrl" > lint_output.txt
set LINT_COUNT=0
for /f %%i in ('type lint_output.txt ^| find /c /v ""') do set LINT_COUNT=%%i
if !LINT_COUNT! gtr 0 (
    echo %YELLOW%golint found !LINT_COUNT! warning(s):%NC%
    type lint_output.txt
    echo.
    echo %YELLOW%Note: Some lint warnings may be acceptable. Review and fix as needed.%NC%
) else (
    echo %GREEN%✓ golint: No issues%NC%
)
del lint_output.txt >nul 2>&1
echo.

REM Final summary
echo === Summary ===
if !FIX_MODE!==1 (
    echo %GREEN%All checks completed with fixes applied%NC%
) else (
    echo %GREEN%All checks completed%NC%
)
echo.
echo To auto-fix formatting issues, run: scripts\lint.bat --fix

endlocal