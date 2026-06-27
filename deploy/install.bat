@echo off
chcp 65001 >nul 2>&1
setlocal enabledelayedexpansion

REM =============================================================================
REM OA-NSDIY 安装脚本 (Windows)
REM =============================================================================
REM 用法:
REM   install.bat                    安装最新版本
REM   install.bat -v v1.0.0          安装指定版本
REM   install.bat upgrade            升级到最新版本
REM   install.bat uninstall          卸载
REM   install.bat list-versions      列出可用版本
REM =============================================================================

set "GITHUB_OWNER=wilson-nsdiy"
set "GITHUB_REPO=nsdiy-office-system"
set "API_BASE=https://api.github.com/repos/%GITHUB_OWNER%/%GITHUB_REPO%"
set "INSTALL_DIR=C:\oa-nsdiy"
set "SERVICE_NAME=oa-nsdiy"

REM Parse arguments
set "TARGET_VERSION="
set "FORCE_YES="
set "COMMAND="
set "NEXT_ARG="

:parse_args
if "%~1"=="" goto :done_parse
if "%NEXT_ARG%"=="version" (
    set "TARGET_VERSION=%~1"
    set "NEXT_ARG="
    shift
    goto :parse_args
)
if "%~1"=="-v" (
    set "NEXT_ARG=version"
    shift
    goto :parse_args
)
if "%~1"=="--version" (
    if not "%~2"=="" (
        set "TARGET_VERSION=%~2"
        shift
    ) else (
        echo [ERROR] --version 需要指定版本号
        exit /b 1
    )
    shift
    goto :parse_args
)
if "%~1"=="-y" set "FORCE_YES=true"
if "%~1"=="--yes" set "FORCE_YES=true"
if not defined COMMAND set "COMMAND=%~1"
shift
goto :parse_args
:done_parse

REM Check admin
net session >nul 2>&1
if errorlevel 1 (
    echo [ERROR] 请以管理员权限运行此脚本
    exit /b 1
)

REM Detect platform
set "OS=windows"
set "ARCH="
if "%PROCESSOR_ARCHITECTURE%"=="AMD64" set "ARCH=amd64"
if "%PROCESSOR_ARCHITECTURE%"=="ARM64" set "ARCH=arm64"
if not defined ARCH (
    echo [ERROR] 不支持的架构: %PROCESSOR_ARCHITECTURE%
    exit /b 1
)
echo [INFO] 检测到平台: %OS%_%ARCH%

echo.
echo ==========================================
echo        OA-NSDIY 安装脚本
echo ==========================================
echo.

if "%COMMAND%"=="upgrade" goto :do_upgrade
if "%COMMAND%"=="update" goto :do_upgrade
if "%COMMAND%"=="uninstall" goto :do_uninstall
if "%COMMAND%"=="remove" goto :do_uninstall
if "%COMMAND%"=="list-versions" goto :do_list_versions
if "%COMMAND%"=="versions" goto :do_list_versions
if "%COMMAND%"=="help" goto :show_usage
if "%COMMAND%"=="--help" goto :show_usage
if "%COMMAND%"=="-h" goto :show_usage
goto :do_install

REM ============================================================
REM Get latest version
REM ============================================================
:get_latest_version
echo [INFO] 正在获取最新版本...
set "RELEASE_JSON="
for /f "delims=" %%j in ('curl -s --connect-timeout 10 --max-time 30 "%API_BASE%/releases/latest" 2^>nul') do set "RELEASE_JSON=!RELEASE_JSON!%%j"
if not defined RELEASE_JSON (
    echo [ERROR] 获取最新版本失败
    echo [INFO] 请到 Release 页面手动下载: https://github.com/%GITHUB_OWNER%/%GITHUB_REPO%/releases
    exit /b 1
)
REM Extract tag_name using findstr
for /f "tokens=2 delims=:" %%t in ('echo !RELEASE_JSON! ^| findstr /r "\"tag_name\"[^,]*"') do (
    for /f "tokens=1 delims=, " %%v in ("%%t") do set "LATEST_VERSION=%%~v"
)
if not defined LATEST_VERSION (
    echo [ERROR] 获取最新版本失败
    exit /b 1
)
echo [INFO] 最新版本: %LATEST_VERSION%
exit /b 0

REM ============================================================
REM List versions
REM ============================================================
:do_list_versions
echo [INFO] 正在获取可用版本...
for /f "delims=" %%j in ('curl -s --connect-timeout 10 --max-time 30 "%API_BASE%/releases" 2^>nul') do set "VERSIONS_JSON=!VERSIONS_JSON!%%j"
if not defined VERSIONS_JSON (
    echo [ERROR] 获取版本列表失败
    exit /b 1
)
echo.
echo 可用版本列表:
echo ----------------------------------------
REM Simple extraction - show tag_names
echo !VERSIONS_JSON! | findstr /r /o "\"tag_name\""
echo ----------------------------------------
echo.
exit /b 0

REM ============================================================
REM Download and extract
REM ============================================================
:download_and_extract
set "VERSION_NUM=%LATEST_VERSION:v=%"
set "ARCHIVE_NAME=oa-nsdiy_%VERSION_NUM%_%OS%_%ARCH%.zip"
set "TEMP_DIR=%TEMP%\oa-nsdiy-install-%RANDOM%"
mkdir "%TEMP_DIR%"

echo [INFO] 正在下载 %ARCHIVE_NAME%...

REM Try GitHub download URL
set "DOWNLOAD_URL=https://github.com/%GITHUB_OWNER%/%GITHUB_REPO%/releases/download/%LATEST_VERSION%/%ARCHIVE_NAME%"

curl -fsSL "%DOWNLOAD_URL%" -o "%TEMP_DIR%\%ARCHIVE_NAME%"
if errorlevel 1 (
    echo [ERROR] 下载失败
    echo [INFO] 请到 Release 页面手动下载: https://github.com/%GITHUB_OWNER%/%GITHUB_REPO%/releases
    rmdir /s /q "%TEMP_DIR%" 2>nul
    exit /b 1
)

echo [INFO] 正在解压...
powershell -NoProfile -Command "Expand-Archive -Path '%TEMP_DIR%\%ARCHIVE_NAME%' -DestinationPath '%TEMP_DIR%\extracted' -Force"
if errorlevel 1 (
    echo [ERROR] 解压失败
    rmdir /s /q "%TEMP_DIR%" 2>nul
    exit /b 1
)

mkdir "%INSTALL_DIR%" 2>nul

REM Copy binary
if exist "%TEMP_DIR%\extracted\oa-nsdiy.exe" (
    copy /y "%TEMP_DIR%\extracted\oa-nsdiy.exe" "%INSTALL_DIR%\oa-nsdiy.exe" >nul
) else if exist "%TEMP_DIR%\extracted\server.exe" (
    copy /y "%TEMP_DIR%\extracted\server.exe" "%INSTALL_DIR%\oa-nsdiy.exe" >nul
) else (
    echo [ERROR] 压缩包中未找到二进制文件
    rmdir /s /q "%TEMP_DIR%" 2>nul
    exit /b 1
)

REM Write version file
echo %LATEST_VERSION%> "%INSTALL_DIR%\.version"

REM Copy .env.example if present in archive
if exist "%TEMP_DIR%\extracted\.env.example" (
    copy /y "%TEMP_DIR%\extracted\.env.example" "%INSTALL_DIR%\.env.example" >nul
)

rmdir /s /q "%TEMP_DIR%" 2>nul

echo [成功] 二进制文件已安装到 %INSTALL_DIR%\oa-nsdiy.exe (%LATEST_VERSION%)
exit /b 0

REM ============================================================
REM Generate .env
REM ============================================================
:generate_env
set "ENV_FILE=%INSTALL_DIR%\.env"
if exist "%ENV_FILE%" (
    echo [INFO] .env 已存在，保留现有配置
    exit /b 0
)

echo [INFO] 正在生成 .env 配置...

REM Generate random JWT secret using PowerShell
for /f "delims=" %%s in ('powershell -NoProfile -Command "[System.BitConverter]::ToString((1..32 ^| ForEach-Object { Get-Random -Maximum 256 }) -as [byte[]]).Replace('-','').ToLower()"') do set "JWT_SECRET=%%s"

(
echo # Generated by install.bat. Edit before production use.
echo # See deploy/.env.example in the repo for full documentation.
echo.
echo # [必须修改] JWT 密钥
echo JWT_SECRET=%JWT_SECRET%
echo.
echo # 服务监听
echo SERVER_HOST=0.0.0.0
echo SERVER_PORT=3001
echo SERVER_MODE=release
echo.
echo # 数据库: sqlite (默认，文件在 data/) 或 postgres
echo DATABASE_DRIVER=sqlite
echo DATABASE_SOURCE=
echo.
echo # 日志
echo LOG_LEVEL=info
echo LOG_FORMAT=json
echo LOG_OUTPUT_TO_STDOUT=true
echo LOG_OUTPUT_TO_FILE=true
) > "%ENV_FILE%"

echo [成功] .env 已生成（含随机 JWT_SECRET）
exit /b 0

REM ============================================================
REM Install
REM ============================================================
:do_install
if defined TARGET_VERSION if exist "%INSTALL_DIR%\oa-nsdiy.exe" goto :install_version

if defined TARGET_VERSION (
    set "LATEST_VERSION=%TARGET_VERSION%"
    call :validate_version
    if errorlevel 1 exit /b 1
    call :fetch_release_by_version
) else (
    call :get_latest_version
    if errorlevel 1 exit /b 1
)

call :download_and_extract
if errorlevel 1 exit /b 1

mkdir "%INSTALL_DIR%\data" 2>nul

call :generate_env
if errorlevel 1 exit /b 1

call :print_completion
exit /b 0

REM ============================================================
REM Install specific version
REM ============================================================
:install_version
set "LATEST_VERSION=%TARGET_VERSION%"
call :validate_version
if errorlevel 1 exit /b 1

echo [INFO] 正在安装指定版本: %LATEST_VERSION%

REM Read current version
set "CURRENT_VERSION=unknown"
if exist "%INSTALL_DIR%\.version" set /p CURRENT_VERSION=<"%INSTALL_DIR%\.version"

echo [INFO] 当前版本: %CURRENT_VERSION%

if "%CURRENT_VERSION%"=="%LATEST_VERSION%" (
    echo [WARN] 已经是该版本，无需操作
    exit /b 0
)

REM Stop service if running
call :stop_service

REM Backup
if not "%CURRENT_VERSION%"=="unknown" (
    copy /y "%INSTALL_DIR%\oa-nsdiy.exe" "%INSTALL_DIR%\oa-nsdiy.backup.%CURRENT_VERSION%.exe" >nul 2>&1
    echo [INFO] 备份已创建: oa-nsdiy.backup.%CURRENT_VERSION%.exe
)

call :fetch_release_by_version
call :download_and_extract
if errorlevel 1 exit /b 1

call :start_service

echo.
echo ==========================================
echo   指定版本安装完成！
echo ==========================================
echo.
echo   当前版本: %LATEST_VERSION%
echo.
exit /b 0

REM ============================================================
REM Upgrade
REM ============================================================
:do_upgrade
if not exist "%INSTALL_DIR%\oa-nsdiy.exe" (
    echo [ERROR] OA-NSDIY 尚未安装，请先执行全新安装
    exit /b 1
)

echo [INFO] 正在升级 OA-NSDIY...

set "CURRENT_VERSION=unknown"
if exist "%INSTALL_DIR%\.version" set /p CURRENT_VERSION=<"%INSTALL_DIR%\.version"
echo [INFO] 当前版本: %CURRENT_VERSION%

call :stop_service

copy /y "%INSTALL_DIR%\oa-nsdiy.exe" "%INSTALL_DIR%\oa-nsdiy.backup.exe" >nul
echo [INFO] 备份已创建: %INSTALL_DIR%\oa-nsdiy.backup.exe

call :get_latest_version
if errorlevel 1 exit /b 1

call :download_and_extract
if errorlevel 1 exit /b 1

call :start_service
echo [成功] 升级完成！
exit /b 0

REM ============================================================
REM Uninstall
REM ============================================================
:do_uninstall
echo [WARN] 这将从系统中移除 OA-NSDIY。

if not "%FORCE_YES%"=="true" (
    set /p CONFIRM="确定要继续吗？(y/N) "
    if /i not "!CONFIRM!"=="y" (
        echo [INFO] 卸载已取消
        exit /b 0
    )
)

call :stop_service

echo [INFO] 正在移除文件...
if exist "%INSTALL_DIR%\oa-nsdiy.exe" del /q "%INSTALL_DIR%\oa-nsdiy.exe"
if exist "%INSTALL_DIR%\.version" del /q "%INSTALL_DIR%\.version"
if exist "%INSTALL_DIR%\oa-nsdiy.backup*.exe" del /q "%INSTALL_DIR%\oa-nsdiy.backup*.exe"

echo [WARN] 数据目录未被移除: %INSTALL_DIR%\data
echo [WARN] 如不再需要，请手动删除: rmdir /s /q "%INSTALL_DIR%"

echo [成功] OA-NSDIY 已卸载
exit /b 0

REM ============================================================
REM Service helpers
REM ============================================================
:stop_service
sc query "%SERVICE_NAME%" >nul 2>&1
if not errorlevel 1 (
    echo [INFO] 正在停止服务...
    net stop "%SERVICE_NAME%" >nul 2>&1
    sc delete "%SERVICE_NAME%" >nul 2>&1
)
exit /b 0

:start_service
echo [INFO] 正在启动服务...
pushd "%INSTALL_DIR%"
start "" "%INSTALL_DIR%\oa-nsdiy.exe"
popd
echo [成功] 服务已启动
echo [INFO] 如需注册为 Windows 服务，请使用 nssm: nssm install %SERVICE_NAME% "%INSTALL_DIR%\oa-nsdiy.exe"
exit /b 0

REM ============================================================
REM Validate version
REM ============================================================
:validate_version
if not defined LATEST_VERSION (
    echo [ERROR] 指定要安装的版本号 (例如: v1.0.0)
    exit /b 1
)
REM Ensure v-prefix
if "%LATEST_VERSION:~0,1%" neq "v" set "LATEST_VERSION=v%LATEST_VERSION%"
echo [INFO] 正在验证版本 %LATEST_VERSION%...
curl -s -o nul -w "%%{http_code}" --connect-timeout 10 --max-time 30 "%API_BASE%/releases/tags/%LATEST_VERSION%" 2>nul | findstr "200" >nul
if errorlevel 1 (
    echo [ERROR] 指定版本不存在: %LATEST_VERSION%
    exit /b 1
)
exit /b 0

REM ============================================================
REM Fetch release by version
REM ============================================================
:fetch_release_by_version
for /f "delims=" %%j in ('curl -s --connect-timeout 10 --max-time 30 "%API_BASE%/releases/tags/%LATEST_VERSION%" 2^>nul') do set "RELEASE_JSON=!RELEASE_JSON!%%j"
exit /b 0

REM ============================================================
REM Completion message
REM ============================================================
:print_completion
echo.
echo ==========================================
echo   OA-NSDIY 安装完成！
echo ==========================================
echo.
echo   安装目录: %INSTALL_DIR%
echo   服务地址: 0.0.0.0:3001
echo.
echo ==========================================
echo   后续步骤
echo ==========================================
echo.
echo   1. 编辑配置文件（修改 JWT_SECRET 等必填项）:
echo      notepad %INSTALL_DIR%\.env
echo.
echo   2. 启动服务:
echo      cd %INSTALL_DIR% ^&^& oa-nsdiy.exe
echo.
echo   3. 如需注册为 Windows 服务，使用 nssm:
echo      nssm install %SERVICE_NAME% "%INSTALL_DIR%\oa-nsdiy.exe"
echo.
echo   4. 访问 Web:
echo      http://YOUR_SERVER_IP:3001
echo.
echo ==========================================
echo   常用命令
echo ==========================================
echo.
echo   查看状态: sc query %SERVICE_NAME%
echo   停止服务: net stop %SERVICE_NAME%
echo   重启服务: net stop %SERVICE_NAME% ^&^& net start %SERVICE_NAME%
echo.
exit /b 0

REM ============================================================
REM Usage
REM ============================================================
:show_usage
echo 用法: %~nx0 [命令] [选项]
echo.
echo 命令:
echo   (无参数)             安装最新版本
echo   install              安装 OA-NSDIY
echo   upgrade              升级到最新版本
echo   uninstall            卸载 OA-NSDIY
echo   list-versions        列出可用版本
echo.
echo 选项:
echo   -v, --version 版本    指定要安装的版本号 (例如: v1.0.0)
echo   -y, --yes            跳过确认提示 (用于卸载)
echo.
echo 示例:
echo   %~nx0                        安装最新版本
echo   %~nx0 install -v v1.0.0      安装指定版本
echo   %~nx0 upgrade                升级到最新
echo   %~nx0 uninstall              卸载
echo   %~nx0 list-versions          列出可用版本
echo.
exit /b 0
