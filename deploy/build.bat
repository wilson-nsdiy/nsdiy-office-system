@echo off
chcp 65001 >nul 2>&1
setlocal enabledelayedexpansion

REM =============================================================================
REM OA-NSDIY 生产构建脚本 (Windows)
REM =============================================================================
REM 一键构建生产环境可执行文件（前端嵌入后端）
REM
REM 用法:
REM   deploy\build.bat              构建当前平台
REM   deploy\build.bat linux         交叉编译 Linux amd64
REM   deploy\build.bat darwin        交叉编译 macOS arm64
REM =============================================================================

set "SCRIPT_DIR=%~dp0"
set "PROJECT_ROOT=%SCRIPT_DIR%.."

REM Parse target platform
set "TARGET_OS=%~1"
for /f "delims=" %%v in ('git -C "%PROJECT_ROOT%" describe --tags --always --dirty 2^>nul') do set "VERSION=%%v"
if not defined VERSION set "VERSION=dev"

if "%TARGET_OS%"=="linux" (
    set "GOOS=linux"
    set "GOARCH=amd64"
    set "OUTPUT_NAME=oa-nsdiy-linux-amd64"
) else if "%TARGET_OS%"=="darwin" (
    set "GOOS=darwin"
    set "GOARCH=arm64"
    set "OUTPUT_NAME=oa-nsdiy-darwin-arm64"
) else if "%TARGET_OS%"=="" (
    for /f "delims=" %%o in ('go env GOOS') do set "GOOS=%%o"
    for /f "delims=" %%a in ('go env GOARCH') do set "GOARCH=%%a"
    set "OUTPUT_NAME=oa-nsdiy"
) else (
    echo [ERROR] 不支持的目标平台: %TARGET_OS% ^(仅支持 linux/darwin^)
    exit /b 1
)

set "OUTPUT_DIR=%PROJECT_ROOT%\backend\bin"

echo.
echo ==========================================
echo   OA-NSDIY 生产构建
echo ==========================================
echo.

REM 1. Build frontend
echo [INFO] 构建前端...
pushd "%PROJECT_ROOT%\frontend"
call npm install --silent
if errorlevel 1 (
    echo [ERROR] npm install 失败
    popd
    exit /b 1
)
call npm run build
if errorlevel 1 (
    echo [ERROR] 前端构建失败
    popd
    exit /b 1
)
popd
echo [INFO] 前端构建完成

REM 2. Copy frontend to embed directory
echo [INFO] 复制前端文件到后端...
set "DIST_DIR=%PROJECT_ROOT%\backend\internal\web\dist"
if exist "%DIST_DIR%\*" del /q "%DIST_DIR%\*" 2>nul
if not exist "%DIST_DIR%" mkdir "%DIST_DIR%"
xcopy /e /y /q "%PROJECT_ROOT%\frontend\out\*" "%DIST_DIR%\" >nul
echo [INFO] 前端文件复制完成

REM 3. Build backend
echo [INFO] 构建后端 (GOOS=%GOOS% GOARCH=%GOARCH%)...
pushd "%PROJECT_ROOT%\backend"
set CGO_ENABLED=0
go build -tags embed -ldflags="-s -w -X main.Version=%VERSION%" -trimpath -o "%OUTPUT_DIR%\%OUTPUT_NAME%.exe" ./cmd/server
if errorlevel 1 (
    echo [ERROR] 后端构建失败
    popd
    exit /b 1
)
popd

REM 4. Done
echo.
echo ==========================================
echo   构建成功!
echo ==========================================
echo.
echo   文件: %OUTPUT_DIR%\%OUTPUT_NAME%.exe
echo   版本: %VERSION%
echo.
echo 部署步骤:
echo   1. 上传 %OUTPUT_NAME%.exe 到目标服务器
echo   2. Windows: 使用 install.bat 安装脚本
echo   3. 或手动部署: copy deploy\.env.example deploy\.env ^&^& %OUTPUT_NAME%.exe
echo.

endlocal
