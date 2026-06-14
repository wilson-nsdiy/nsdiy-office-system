#!/bin/bash
#
# OA-NSDIY Installation Script
# OA-NSDIY 安装脚本
# Usage: curl -fsSL https://gitee.com/zhouws-chn/oa-nsdiy/raw/master/deploy/install.sh | sudo bash
#

set -e

# Bash 4+ is required for associative arrays used by the localized message table.
if [ -z "${BASH_VERSION:-}" ]; then
    echo "Error: This installer must be run with Bash 4.0 or later." >&2
    echo "Please install Bash 4+ and run it with that interpreter." >&2
    exit 1
fi

BASH_MAJOR_VERSION="${BASH_VERSION%%.*}"
if [ "$BASH_MAJOR_VERSION" -lt 4 ]; then
    echo "Error: Bash 4.0 or later is required. Current version: $BASH_VERSION" >&2
    exit 1
fi

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Configuration
GITEE_OWNER="zhouws-chn"
GITEE_REPO="oa-nsdiy"
INSTALL_DIR="/opt/oa-nsdiy"
SERVICE_NAME="oa-nsdiy"
SERVICE_USER="oa-nsdiy"
CONFIG_DIR="/etc/oa-nsdiy"

# Server configuration (will be set by user)
SERVER_HOST="0.0.0.0"
SERVER_PORT="3001"

# Language (default: zh = Chinese)
LANG_CHOICE="zh"

# ============================================================
# Language strings / 语言字符串
# ============================================================

declare -A MSG_ZH=(
    ["info"]="信息"
    ["success"]="成功"
    ["warning"]="警告"
    ["error"]="错误"

    ["select_lang"]="请选择语言 / Select language"
    ["lang_zh"]="中文"
    ["lang_en"]="English"
    ["enter_choice"]="请输入选择 (默认: 1)"

    ["install_title"]="OA-NSDIY 安装脚本"
    ["run_as_root"]="请使用 root 权限运行 (使用 sudo)"
    ["detected_platform"]="检测到平台"
    ["unsupported_arch"]="不支持的架构"
    ["unsupported_os"]="不支持的操作系统"
    ["missing_deps"]="缺少依赖"
    ["install_deps_first"]="请先安装以下依赖"
    ["fetching_version"]="正在获取最新版本..."
    ["latest_version"]="最新版本"
    ["failed_get_version"]="获取最新版本失败"
    ["downloading"]="正在下载"
    ["download_failed"]="下载失败"
    ["verifying_checksum"]="正在校验文件..."
    ["checksum_verified"]="校验通过"
    ["checksum_failed"]="校验失败"
    ["checksum_not_found"]="无法验证校验和（checksums.txt 未找到）"
    ["extracting"]="正在解压..."
    ["binary_installed"]="二进制文件已安装到"
    ["user_exists"]="用户已存在"
    ["creating_user"]="正在创建系统用户"
    ["user_created"]="用户已创建"
    ["setting_up_dirs"]="正在设置目录..."
    ["dirs_configured"]="目录配置完成"
    ["installing_service"]="正在安装 systemd 服务..."
    ["service_installed"]="systemd 服务已安装"

    ["install_complete"]="OA-NSDIY 安装完成！"
    ["install_dir"]="安装目录"
    ["next_steps"]="后续步骤"
    ["step1_edit_env"]="编辑配置文件（修改 JWT_SECRET 等必填项）"
    ["step2_start_service"]="启动服务"
    ["step3_enable_autostart"]="设置开机自启"
    ["step4_access"]="访问 Web"
    ["useful_commands"]="常用命令"
    ["cmd_status"]="查看状态"
    ["cmd_logs"]="查看日志"
    ["cmd_restart"]="重启服务"
    ["cmd_stop"]="停止服务"

    ["generating_env"]="正在生成 .env 配置..."
    ["env_generated"]=".env 已生成（含随机 JWT_SECRET）"
    ["edit_env_hint"]="请编辑 %s 修改配置，尤其是："

    ["upgrading"]="正在升级 OA-NSDIY..."
    ["current_version"]="当前版本"
    ["stopping_service"]="正在停止服务..."
    ["backup_created"]="备份已创建"
    ["starting_service"]="正在启动服务..."
    ["upgrade_complete"]="升级完成！"

    ["installing_version"]="正在安装指定版本"
    ["version_not_found"]="指定版本不存在"
    ["same_version"]="已经是该版本，无需操作"
    ["rollback_complete"]="版本回退完成！"
    ["install_version_complete"]="指定版本安装完成！"
    ["validating_version"]="正在验证版本..."
    ["available_versions"]="可用版本列表"
    ["fetching_versions"]="正在获取可用版本..."
    ["not_installed"]="OA-NSDIY 尚未安装，请先执行全新安装"
    ["fresh_install_hint"]="用法"

    ["uninstall_confirm"]="这将从系统中移除 OA-NSDIY。"
    ["are_you_sure"]="确定要继续吗？(y/N)"
    ["uninstall_cancelled"]="卸载已取消"
    ["removing_files"]="正在移除文件..."
    ["removing_install_dir"]="正在移除安装目录..."
    ["removing_user"]="正在移除用户..."
    ["config_not_removed"]="配置目录未被移除"
    ["remove_manually"]="如不再需要，请手动删除"
    ["purge_prompt"]="是否同时删除数据目录？这将清除所有数据 [y/N]: "
    ["removing_config_dir"]="正在移除数据目录..."
    ["uninstall_complete"]="OA-NSDIY 已卸载"

    ["usage"]="用法"
    ["cmd_none"]="(无参数)"
    ["cmd_install"]="安装 OA-NSDIY"
    ["cmd_upgrade"]="升级到最新版本"
    ["cmd_uninstall"]="卸载 OA-NSDIY"
    ["cmd_install_version"]="安装/回退到指定版本"
    ["cmd_list_versions"]="列出可用版本"
    ["opt_version"]="指定要安装的版本号 (例如: v1.0.0)"

    ["server_config_title"]="服务器配置"
    ["server_config_desc"]="配置 OA-NSDIY 服务监听地址"
    ["server_host_prompt"]="服务器监听地址"
    ["server_host_hint"]="0.0.0.0 表示监听所有网卡，127.0.0.1 仅本地访问"
    ["server_port_prompt"]="服务器端口"
    ["server_port_hint"]="建议使用 1024-65535 之间的端口"
    ["server_config_summary"]="服务器配置"
    ["invalid_port"]="无效端口号，请输入 1-65535 之间的数字"

    ["starting_service"]="正在启动服务..."
    ["service_started"]="服务已启动"
    ["service_start_failed"]="服务启动失败，请检查日志"
    ["enabling_autostart"]="正在设置开机自启..."
    ["autostart_enabled"]="开机自启已启用"
    ["getting_public_ip"]="正在获取公网 IP..."
    ["public_ip_failed"]="无法获取公网 IP，使用本地 IP"

    ["manual_download_hint"]="如自动下载失败，请到 Release 页面手动下载"
)

declare -A MSG_EN=(
    ["info"]="INFO"
    ["success"]="SUCCESS"
    ["warning"]="WARNING"
    ["error"]="ERROR"

    ["select_lang"]="请选择语言 / Select language"
    ["lang_zh"]="中文"
    ["lang_en"]="English"
    ["enter_choice"]="Enter your choice (default: 1)"

    ["install_title"]="OA-NSDIY Installation Script"
    ["run_as_root"]="Please run as root (use sudo)"
    ["detected_platform"]="Detected platform"
    ["unsupported_arch"]="Unsupported architecture"
    ["unsupported_os"]="Unsupported OS"
    ["missing_deps"]="Missing dependencies"
    ["install_deps_first"]="Please install them first"
    ["fetching_version"]="Fetching latest version..."
    ["latest_version"]="Latest version"
    ["failed_get_version"]="Failed to get latest version"
    ["downloading"]="Downloading"
    ["download_failed"]="Download failed"
    ["verifying_checksum"]="Verifying checksum..."
    ["checksum_verified"]="Checksum verified"
    ["checksum_failed"]="Checksum verification failed"
    ["checksum_not_found"]="Could not verify checksum (checksums.txt not found)"
    ["extracting"]="Extracting..."
    ["binary_installed"]="Binary installed to"
    ["user_exists"]="User already exists"
    ["creating_user"]="Creating system user"
    ["user_created"]="User created"
    ["setting_up_dirs"]="Setting up directories..."
    ["dirs_configured"]="Directories configured"
    ["installing_service"]="Installing systemd service..."
    ["service_installed"]="Systemd service installed"

    ["install_complete"]="OA-NSDIY installation completed!"
    ["install_dir"]="Installation directory"
    ["next_steps"]="NEXT STEPS"
    ["step1_edit_env"]="Edit the config file (set JWT_SECRET and other required fields)"
    ["step2_start_service"]="Start the service"
    ["step3_enable_autostart"]="Enable auto-start on boot"
    ["step4_access"]="Access the web UI"
    ["useful_commands"]="USEFUL COMMANDS"
    ["cmd_status"]="Check status"
    ["cmd_logs"]="View logs"
    ["cmd_restart"]="Restart"
    ["cmd_stop"]="Stop"

    ["generating_env"]="Generating .env configuration..."
    ["env_generated"]=".env generated (with random JWT_SECRET)"
    ["edit_env_hint"]="Edit %s to change settings, especially:"

    ["upgrading"]="Upgrading OA-NSDIY..."
    ["current_version"]="Current version"
    ["stopping_service"]="Stopping service..."
    ["backup_created"]="Backup created"
    ["starting_service"]="Starting service..."
    ["upgrade_complete"]="Upgrade completed!"

    ["installing_version"]="Installing specified version"
    ["version_not_found"]="Specified version not found"
    ["same_version"]="Already at this version, no action needed"
    ["rollback_complete"]="Version rollback completed!"
    ["install_version_complete"]="Specified version installed!"
    ["validating_version"]="Validating version..."
    ["available_versions"]="Available versions"
    ["fetching_versions"]="Fetching available versions..."
    ["not_installed"]="OA-NSDIY is not installed. Please run a fresh install first"
    ["fresh_install_hint"]="Usage"

    ["uninstall_confirm"]="This will remove OA-NSDIY from your system."
    ["are_you_sure"]="Are you sure? (y/N)"
    ["uninstall_cancelled"]="Uninstall cancelled"
    ["removing_files"]="Removing files..."
    ["removing_install_dir"]="Removing installation directory..."
    ["removing_user"]="Removing user..."
    ["config_not_removed"]="Config directory was NOT removed."
    ["remove_manually"]="Remove it manually if you no longer need it."
    ["purge_prompt"]="Also remove data directory? This will delete ALL data [y/N]: "
    ["removing_config_dir"]="Removing data directory..."
    ["uninstall_complete"]="OA-NSDIY has been uninstalled"

    ["usage"]="Usage"
    ["cmd_none"]="(none)"
    ["cmd_install"]="Install OA-NSDIY"
    ["cmd_upgrade"]="Upgrade to the latest version"
    ["cmd_uninstall"]="Remove OA-NSDIY"
    ["cmd_install_version"]="Install/rollback to a specific version"
    ["cmd_list_versions"]="List available versions"
    ["opt_version"]="Specify version to install (e.g., v1.0.0)"

    ["server_config_title"]="Server Configuration"
    ["server_config_desc"]="Configure OA-NSDIY server listen address"
    ["server_host_prompt"]="Server listen address"
    ["server_host_hint"]="0.0.0.0 listens on all interfaces, 127.0.0.1 for local only"
    ["server_port_prompt"]="Server port"
    ["server_port_hint"]="Recommended range: 1024-65535"
    ["server_config_summary"]="Server configuration"
    ["invalid_port"]="Invalid port number, please enter a number between 1-65535"

    ["starting_service"]="Starting service..."
    ["service_started"]="Service started"
    ["service_start_failed"]="Service failed to start, please check logs"
    ["enabling_autostart"]="Enabling auto-start on boot..."
    ["autostart_enabled"]="Auto-start enabled"
    ["getting_public_ip"]="Getting public IP..."
    ["public_ip_failed"]="Failed to get public IP, using local IP"

    ["manual_download_hint"]="If automatic download fails, please download manually from the Release page"
)

# Get message based on current language
msg() {
    local key="$1"
    if [ "$LANG_CHOICE" = "en" ]; then
        echo "${MSG_EN[$key]}"
    else
        echo "${MSG_ZH[$key]}"
    fi
}

print_info()    { echo -e "${BLUE}[$(msg 'info')]${NC} $1"; }
print_success() { echo -e "${GREEN}[$(msg 'success')]${NC} $1"; }
print_warning() { echo -e "${YELLOW}[$(msg 'warning')]${NC} $1"; }
print_error()   { echo -e "${RED}[$(msg 'error')]${NC} $1"; }

# Check if running interactively (works even when piped via curl | bash)
is_interactive() {
    [ -e /dev/tty ] && [ -r /dev/tty ] && [ -w /dev/tty ]
}

# ============================================================
# Language selection
# ============================================================
select_language() {
    if ! is_interactive; then
        LANG_CHOICE="zh"
        return
    fi

    echo ""
    echo -e "${CYAN}=============================================="
    echo "  $(msg 'select_lang')"
    echo "==============================================${NC}"
    echo ""
    echo "  1) $(msg 'lang_zh') (默认/default)"
    echo "  2) $(msg 'lang_en')"
    echo ""

    read -p "$(msg 'enter_choice'): " lang_input < /dev/tty
    case "$lang_input" in
        2|en|EN|english|English) LANG_CHOICE="en" ;;
        *) LANG_CHOICE="zh" ;;
    esac
    echo ""
}

# ============================================================
# Server configuration
# ============================================================
validate_port() {
    local port="$1"
    [[ "$port" =~ ^[0-9]+$ ]] && [ "$port" -ge 1 ] && [ "$port" -le 65535 ]
}

configure_server() {
    if ! is_interactive; then
        print_info "$(msg 'server_config_summary'): ${SERVER_HOST}:${SERVER_PORT} (default)"
        return
    fi

    echo ""
    echo -e "${CYAN}=============================================="
    echo "  $(msg 'server_config_title')"
    echo "==============================================${NC}"
    echo ""
    echo -e "${BLUE}$(msg 'server_config_desc')${NC}"
    echo ""

    echo -e "${YELLOW}$(msg 'server_host_hint')${NC}"
    read -p "$(msg 'server_host_prompt') [${SERVER_HOST}]: " input_host < /dev/tty
    [ -n "$input_host" ] && SERVER_HOST="$input_host"

    echo ""
    echo -e "${YELLOW}$(msg 'server_port_hint')${NC}"
    while true; do
        read -p "$(msg 'server_port_prompt') [${SERVER_PORT}]: " input_port < /dev/tty
        if [ -z "$input_port" ]; then
            break
        elif validate_port "$input_port"; then
            SERVER_PORT="$input_port"
            break
        else
            print_error "$(msg 'invalid_port')"
        fi
    done

    echo ""
    print_info "$(msg 'server_config_summary'): ${SERVER_HOST}:${SERVER_PORT}"
    echo ""
}

# ============================================================
# Root / platform / dependency checks
# ============================================================
check_root() {
    if [ "$(id -u)" -ne 0 ]; then
        print_error "$(msg 'run_as_root')"
        exit 1
    fi
}

detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)

    case "$ARCH" in
        x86_64) ARCH="amd64" ;;
        aarch64|arm64) ARCH="arm64" ;;
        *) print_error "$(msg 'unsupported_arch'): $ARCH"; exit 1 ;;
    esac

    case "$OS" in
        linux) OS="linux" ;;
        *) print_error "$(msg 'unsupported_os'): $OS"; exit 1 ;;
    esac

    print_info "$(msg 'detected_platform'): ${OS}_${ARCH}"
}

check_dependencies() {
    local missing=()
    command -v curl &>/dev/null || missing+=("curl")
    command -v tar   &>/dev/null || missing+=("tar")
    [ ${#missing[@]} -gt 0 ] && {
        print_error "$(msg 'missing_deps'): ${missing[*]}"
        print_info "$(msg 'install_deps_first')"
        exit 1
    }
}

# ============================================================
# Gitee Release helpers
# ============================================================
# Gitee API base (no trailing slash)
API_BASE="https://gitee.com/api/v5/repos/${GITEE_OWNER}/${GITEE_REPO}"

# Fetch latest release version tag (and cache the raw JSON)
# Sets globals: LATEST_VERSION, RELEASE_ID, RELEASE_JSON
get_latest_version() {
    print_info "$(msg 'fetching_version')"
    RELEASE_JSON=$(curl -s --connect-timeout 10 --max-time 30 "${API_BASE}/releases/latest" 2>/dev/null)

    if [ -z "$RELEASE_JSON" ]; then
        print_error "$(msg 'failed_get_version')"
        print_info "$(msg 'manual_download_hint'): https://gitee.com/${GITEE_OWNER}/${GITEE_REPO}/releases"
        exit 1
    fi

    # tag_name is the version (e.g. v1.0.0). Use grep+sed to avoid jq dependency.
    LATEST_VERSION=$(echo "$RELEASE_JSON" | grep -o '"tag_name"[[:space:]]*:[[:space:]]*"[^"]*"' | head -1 | sed -E 's/.*"([^"]+)"$/\1/')
    RELEASE_ID=$(echo "$RELEASE_JSON" | grep -o '"id"[[:space:]]*:[[:space:]]*[0-9]*' | head -1 | sed -E 's/.*:([0-9]*)/\1/')

    if [ -z "$LATEST_VERSION" ]; then
        print_error "$(msg 'failed_get_version')"
        print_info "$(msg 'manual_download_hint'): https://gitee.com/${GITEE_OWNER}/${GITEE_REPO}/releases"
        exit 1
    fi

    print_info "$(msg 'latest_version'): $LATEST_VERSION"
}

# List available versions (from Gitee releases)
list_versions() {
    print_info "$(msg 'fetching_versions')"
    local versions
    versions=$(curl -s --connect-timeout 10 --max-time 30 "${API_BASE}/releases" 2>/dev/null | grep -o '"tag_name"[[:space:]]*:[[:space:]]*"[^"]*"' | sed -E 's/.*"([^"]+)"$/\1/' | head -20)

    if [ -z "$versions" ]; then
        print_error "$(msg 'failed_get_version')"
        print_info "$(msg 'manual_download_hint'): https://gitee.com/${GITEE_OWNER}/${GITEE_REPO}/releases"
        exit 1
    fi

    echo ""
    echo "$(msg 'available_versions'):"
    echo "----------------------------------------"
    echo "$versions" | while read -r version; do echo "  $version"; done
    echo "----------------------------------------"
    echo ""
}

# Validate a version tag exists on Gitee. Outputs normalized version (v-prefixed) to stdout.
validate_version() {
    local version="$1"
    [ -z "$version" ] && { print_error "$(msg 'opt_version')" >&2; exit 1; }

    # Ensure v-prefix
    [[ "$version" =~ ^v ]] || version="v$version"

    print_info "$(msg 'validating_version') $version" >&2

    local http_code
    http_code=$(curl -s -o /dev/null -w "%{http_code}" --connect-timeout 10 --max-time 30 "${API_BASE}/releases/tags/${version}" 2>/dev/null)
    [ -z "$http_code" ] || ! [[ "$http_code" =~ ^[0-9]+$ ]] && {
        print_error "Network error: Failed to connect to Gitee API" >&2
        exit 1
    }
    [ "$http_code" != "200" ] && {
        print_error "$(msg 'version_not_found'): $version" >&2
        echo "" >&2
        list_versions >&2
        exit 1
    }
    echo "$version"
}

# Fetch the release JSON for a specific version tag into RELEASE_JSON / RELEASE_ID.
fetch_release_by_version() {
    local version="$1"
    RELEASE_JSON=$(curl -s --connect-timeout 10 --max-time 30 "${API_BASE}/releases/tags/${version}" 2>/dev/null)
    RELEASE_ID=$(echo "$RELEASE_JSON" | grep -o '"id"[[:space:]]*:[[:space:]]*[0-9]*' | head -1 | sed -E 's/.*:([0-9]*)/\1/')
}

# Resolve the download URL for a given asset file name from the cached RELEASE_JSON.
# Gitee attaches files with an internal attach_file_id; the download endpoint is
#   /v5/repos/{owner}/{repo}/releases/{release_id}/attach_files/{attach_file_id}/download
# Strategy:
#   1) Parse assets array for matching file name -> get attach_file_id
#   2) If id missing, fall back to GitHub-style URL (some Gitee versions support it)
# Outputs URL to stdout; returns 1 on failure.
get_asset_url() {
    local asset_name="$1"
    local json="$RELEASE_JSON"

    # Try to find the attach_file_id for the matching asset name.
    # Gitee asset objects look like: {"name":"x.tar.gz", "id":12345, ...}
    # We extract id that follows (after some chars) the matching name line.
    local attach_id
    attach_id=$(echo "$json" | tr ',' '\n' | grep -A0 "\"name\"[[:space:]]*:[[:space:]]*\"${asset_name}\"" >/dev/null 2>&1 && \
        echo "$json" | python3 -c "
import sys, json, re
try:
    data = json.load(sys.stdin)
    assets = data.get('assets', []) if isinstance(data, dict) else []
    for a in assets:
        if isinstance(a, dict) and a.get('name') == '${asset_name}':
            print(a.get('id', '')); break
except Exception:
    pass
" 2>/dev/null)

    if [ -n "$attach_id" ] && [ -n "$RELEASE_ID" ]; then
        echo "${API_BASE}/releases/${RELEASE_ID}/attach_files/${attach_id}/download"
        return 0
    fi

    # Fallback 1: GitHub-style clean URL (works on some Gitee versions)
    if [ -n "$LATEST_VERSION" ]; then
        echo "https://gitee.com/${GITEE_OWNER}/${GITEE_REPO}/releases/download/${LATEST_VERSION}/${asset_name}"
        return 0
    fi

    return 1
}

# ============================================================
# Version detection of installed binary
# ============================================================
get_current_version() {
    if [ -f "$INSTALL_DIR/.version" ]; then
        cat "$INSTALL_DIR/.version" 2>/dev/null || echo "unknown"
    elif [ -f "$INSTALL_DIR/oa-nsdiy" ]; then
        "$INSTALL_DIR/oa-nsdiy" --version 2>/dev/null | grep -oE 'v?[0-9]+\.[0-9]+\.[0-9]+' | head -1 || echo "unknown"
    else
        echo "not_installed"
    fi
}

# ============================================================
# Download & extract
# ============================================================
download_and_extract() {
    local version_num=${LATEST_VERSION#v}
    local archive_name="oa-nsdiy_${version_num}_${OS}_${ARCH}.tar.gz"
    local checksum_asset="checksums.txt"

    print_info "$(msg 'downloading') ${archive_name}..."

    TEMP_DIR=$(mktemp -d)
    trap "rm -rf $TEMP_DIR" EXIT

    # Resolve archive URL (and fall back to direct GitHub-style if needed)
    local download_url
    download_url=$(get_asset_url "$archive_name")
    if [ -z "$download_url" ]; then
        # Last resort: assume GitHub-style URL
        download_url="https://gitee.com/${GITEE_OWNER}/${GITEE_REPO}/releases/download/${LATEST_VERSION}/${archive_name}"
    fi

    if ! curl -fsSL "$download_url" -o "$TEMP_DIR/$archive_name"; then
        print_error "$(msg 'download_failed')"
        print_info "$(msg 'manual_download_hint'): https://gitee.com/${GITEE_OWNER}/${GITEE_REPO}/releases"
        exit 1
    fi

    # Download & verify checksum
    print_info "$(msg 'verifying_checksum')"
    local checksum_url
    checksum_url=$(get_asset_url "$checksum_asset")
    if [ -n "$checksum_url" ] && curl -fsSL "$checksum_url" -o "$TEMP_DIR/checksums.txt" 2>/dev/null; then
        local expected_checksum actual_checksum
        expected_checksum=$(grep "$archive_name" "$TEMP_DIR/checksums.txt" | awk '{print $1}')
        actual_checksum=$(sha256sum "$TEMP_DIR/$archive_name" | awk '{print $1}')
        if [ -n "$expected_checksum" ] && [ "$expected_checksum" != "$actual_checksum" ]; then
            print_error "$(msg 'checksum_failed')"
            print_error "Expected: $expected_checksum"
            print_error "Actual:   $actual_checksum"
            exit 1
        fi
        [ -n "$expected_checksum" ] && print_success "$(msg 'checksum_verified')" || print_warning "$(msg 'checksum_not_found')"
    else
        print_warning "$(msg 'checksum_not_found')"
    fi

    # Extract
    print_info "$(msg 'extracting')"
    tar -xzf "$TEMP_DIR/$archive_name" -C "$TEMP_DIR"

    mkdir -p "$INSTALL_DIR"

    # Copy binary (try common names)
    if [ -f "$TEMP_DIR/oa-nsdiy" ]; then
        cp "$TEMP_DIR/oa-nsdiy" "$INSTALL_DIR/oa-nsdiy"
    elif [ -f "$TEMP_DIR/server" ]; then
        cp "$TEMP_DIR/server" "$INSTALL_DIR/oa-nsdiy"
    else
        print_error "Binary not found in archive"
        exit 1
    fi
    chmod +x "$INSTALL_DIR/oa-nsdiy"

    # Persist version file (for upgrade/rollback detection)
    echo "$LATEST_VERSION" > "$INSTALL_DIR/.version"

    print_success "$(msg 'binary_installed') $INSTALL_DIR/oa-nsdiy ($LATEST_VERSION)"
}

# ============================================================
# User / directories / .env
# ============================================================
create_user() {
    if id "$SERVICE_USER" &>/dev/null; then
        print_info "$(msg 'user_exists'): $SERVICE_USER"
        local current_shell
        current_shell=$(getent passwd "$SERVICE_USER" 2>/dev/null | cut -d: -f7)
        if [ "$current_shell" = "/bin/false" ] || [ "$current_shell" = "/sbin/nologin" ]; then
            if usermod -s /bin/sh "$SERVICE_USER" 2>/dev/null; then
                print_success "User shell updated to /bin/sh"
            fi
        fi
    else
        print_info "$(msg 'creating_user') $SERVICE_USER..."
        useradd -r -s /bin/sh -d "$INSTALL_DIR" "$SERVICE_USER"
        print_success "$(msg 'user_created')"
    fi
}

setup_directories() {
    print_info "$(msg 'setting_up_dirs')"
    mkdir -p "$INSTALL_DIR" "$INSTALL_DIR/data" "$CONFIG_DIR"
    chown -R "$SERVICE_USER:$SERVICE_USER" "$INSTALL_DIR" "$CONFIG_DIR"
    print_success "$(msg 'dirs_configured')"
}

# Generate .env at $INSTALL_DIR/.env with a random JWT_SECRET if absent.
generate_env() {
    local env_file="$INSTALL_DIR/.env"
    [ -f "$env_file" ] && { print_info ".env already exists, keeping existing"; return; }

    print_info "$(msg 'generating_env')"
    local jwt_secret
    jwt_secret=$(openssl rand -hex 32 2>/dev/null || head -c 32 /dev/urandom | od -An -tx1 | tr -d ' \n')

    cat > "$env_file" << EOF
# Generated by install.sh. Edit before production use.
# See deploy/.env.example in the repo for full documentation.

# [必须修改] JWT 密钥
JWT_SECRET=${jwt_secret}

# 服务监听
SERVER_HOST=${SERVER_HOST}
SERVER_PORT=${SERVER_PORT}
SERVER_MODE=release

# 数据库: sqlite (默认，文件在 data/) 或 postgres
DATABASE_DRIVER=sqlite
DATABASE_SOURCE=

# 日志
LOG_LEVEL=info
LOG_FORMAT=json
LOG_OUTPUT_TO_STDOUT=true
LOG_OUTPUT_TO_FILE=true
EOF
    chmod 600 "$env_file"
    chown "$SERVICE_USER:$SERVICE_USER" "$env_file"
    print_success "$(msg 'env_generated')"
}

# ============================================================
# systemd service
# ============================================================
install_service() {
    print_info "$(msg 'installing_service')"

    cat > /etc/systemd/system/${SERVICE_NAME}.service << EOF
[Unit]
Description=OA-NSDIY - Studio OA Management System
Documentation=https://gitee.com/${GITEE_OWNER}/${GITEE_REPO}
After=network.target

[Service]
Type=simple
User=${SERVICE_USER}
Group=${SERVICE_USER}
WorkingDirectory=${INSTALL_DIR}
ExecStart=${INSTALL_DIR}/${SERVICE_NAME}
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal
SyslogIdentifier=${SERVICE_NAME}

# Read configuration from .env
EnvironmentFile=${INSTALL_DIR}/.env

# Security hardening
NoNewPrivileges=true
ProtectSystem=strict
ProtectHome=true
PrivateTmp=true
ReadWritePaths=${INSTALL_DIR}

[Install]
WantedBy=multi-user.target
EOF

    systemctl daemon-reload
    print_success "$(msg 'service_installed')"
}

# ============================================================
# Service management
# ============================================================
start_service() {
    print_info "$(msg 'starting_service')"
    if systemctl start "$SERVICE_NAME"; then
        print_success "$(msg 'service_started')"
    else
        print_error "$(msg 'service_start_failed')"
        print_info "sudo journalctl -u ${SERVICE_NAME} -n 50"
        return 1
    fi
}

enable_autostart() {
    print_info "$(msg 'enabling_autostart')"
    systemctl enable "$SERVICE_NAME" 2>/dev/null && print_success "$(msg 'autostart_enabled')" || print_warning "Failed to enable auto-start"
}

get_public_ip() {
    print_info "$(msg 'getting_public_ip')"
    local response
    response=$(curl -s --connect-timeout 5 --max-time 10 "https://ipinfo.io/json" 2>/dev/null)
    if [ -n "$response" ]; then
        PUBLIC_IP=$(echo "$response" | grep -o '"ip": *"[^"]*"' | sed 's/"ip": *"\([^"]*\)"/\1/')
        [ -n "$PUBLIC_IP" ] && { print_success "Public IP: $PUBLIC_IP"; return 0; }
    fi
    print_warning "$(msg 'public_ip_failed')"
    PUBLIC_IP=$(hostname -I 2>/dev/null | awk '{print $1}' || echo "YOUR_SERVER_IP")
}

# ============================================================
# Completion message
# ============================================================
print_completion() {
    local display_host="${PUBLIC_IP:-YOUR_SERVER_IP}"
    [ "$SERVER_HOST" = "127.0.0.1" ] && display_host="127.0.0.1"

    echo ""
    echo "=============================================="
    print_success "$(msg 'install_complete')"
    echo "=============================================="
    echo ""
    echo "$(msg 'install_dir'): $INSTALL_DIR"
    echo "$(msg 'server_config_summary'): ${SERVER_HOST}:${SERVER_PORT}"
    echo ""
    echo "=============================================="
    echo "  $(msg 'next_steps')"
    echo "=============================================="
    echo ""
    echo "  1. $(msg 'step1_edit_env'):"
    echo "     sudo nano ${INSTALL_DIR}/.env"
    echo "     # $(msg 'edit_env_hint' "$INSTALL_DIR/.env")"
    echo "     #   - JWT_SECRET  ($(msg 'step1_edit_env'))"
    echo "     #   - CORS_ALLOW_ORIGINS"
    echo ""
    echo "  2. $(msg 'step2_start_service'):"
    echo "     sudo systemctl restart ${SERVICE_NAME}"
    echo ""
    echo "  3. $(msg 'step3_enable_autostart'):"
    echo "     sudo systemctl enable ${SERVICE_NAME}"
    echo ""
    echo "  4. $(msg 'step4_access'):"
    echo "     http://${display_host}:${SERVER_PORT}"
    echo "     ($(msg 'step4_access') admin / admin123)"
    echo ""
    echo "=============================================="
    echo "  $(msg 'useful_commands')"
    echo "=============================================="
    echo ""
    echo "  $(msg 'cmd_status'):   sudo systemctl status ${SERVICE_NAME}"
    echo "  $(msg 'cmd_logs'):     sudo journalctl -u ${SERVICE_NAME} -f"
    echo "  $(msg 'cmd_restart'):  sudo systemctl restart ${SERVICE_NAME}"
    echo "  $(msg 'cmd_stop'):     sudo systemctl stop ${SERVICE_NAME}"
    echo ""
    echo "=============================================="
}

# ============================================================
# Upgrade / install-version / uninstall
# ============================================================
upgrade() {
    [ ! -f "$INSTALL_DIR/oa-nsdiy" ] && {
        print_error "$(msg 'not_installed')"
        print_info "$(msg 'fresh_install_hint'): $0 install"
        exit 1
    }

    print_info "$(msg 'upgrading')"

    local current_version
    current_version=$(get_current_version)
    print_info "$(msg 'current_version'): $current_version"

    systemctl is-active --quiet "$SERVICE_NAME" && {
        print_info "$(msg 'stopping_service')"
        systemctl stop "$SERVICE_NAME"
    }

    cp "$INSTALL_DIR/oa-nsdiy" "$INSTALL_DIR/oa-nsdiy.backup"
    print_info "$(msg 'backup_created'): $INSTALL_DIR/oa-nsdiy.backup"

    get_latest_version
    download_and_extract
    chown "$SERVICE_USER:$SERVICE_USER" "$INSTALL_DIR/oa-nsdiy"

    print_info "$(msg 'starting_service')"
    systemctl start "$SERVICE_NAME"
    print_success "$(msg 'upgrade_complete')"
}

install_version() {
    local target_version="$1"

    [ ! -f "$INSTALL_DIR/oa-nsdiy" ] && {
        print_error "$(msg 'not_installed')"
        print_info "$(msg 'fresh_install_hint'): $0 install -v $target_version"
        exit 1
    }

    target_version=$(validate_version "$target_version")
    print_info "$(msg 'installing_version'): $target_version"

    local current_version
    current_version=$(get_current_version)
    print_info "$(msg 'current_version'): $current_version"

    if [ "$current_version" = "$target_version" ] || [ "$current_version" = "${target_version#v}" ]; then
        print_warning "$(msg 'same_version')"
        exit 0
    fi

    systemctl is-active --quiet "$SERVICE_NAME" && {
        print_info "$(msg 'stopping_service')"
        systemctl stop "$SERVICE_NAME"
    }

    local backup_name
    if [ "$current_version" != "unknown" ] && [ "$current_version" != "not_installed" ]; then
        backup_name="oa-nsdiy.backup.${current_version}"
    else
        backup_name="oa-nsdiy.backup.$(date +%Y%m%d%H%M%S)"
    fi
    cp "$INSTALL_DIR/oa-nsdiy" "$INSTALL_DIR/$backup_name"
    print_info "$(msg 'backup_created'): $INSTALL_DIR/$backup_name"

    LATEST_VERSION="$target_version"
    fetch_release_by_version "$target_version"
    download_and_extract
    chown "$SERVICE_USER:$SERVICE_USER" "$INSTALL_DIR/oa-nsdiy"

    print_info "$(msg 'starting_service')"
    systemctl start "$SERVICE_NAME" && print_success "$(msg 'service_started')" || {
        print_error "$(msg 'service_start_failed')"
        print_info "sudo journalctl -u ${SERVICE_NAME} -n 50"
    }

    local new_version
    new_version=$(get_current_version)
    echo ""
    echo "=============================================="
    print_success "$(msg 'install_version_complete')"
    echo "=============================================="
    echo ""
    echo "  $(msg 'current_version'): $new_version"
    echo ""
}

uninstall() {
    print_warning "$(msg 'uninstall_confirm')"

    if ! is_interactive; then
        if [ "${FORCE_YES:-}" != "true" ]; then
            print_error "Non-interactive mode detected. Use 'curl ... | bash -s -- uninstall -y' to confirm."
            exit 1
        fi
    else
        read -p "$(msg 'are_you_sure') " -n 1 -r < /dev/tty
        echo
        [[ $REPLY =~ ^[Yy]$ ]] || { print_info "$(msg 'uninstall_cancelled')"; exit 0; }
    fi

    print_info "$(msg 'stopping_service')"
    systemctl stop "$SERVICE_NAME" 2>/dev/null || true
    systemctl disable "$SERVICE_NAME" 2>/dev/null || true

    print_info "$(msg 'removing_files')"
    rm -f /etc/systemd/system/${SERVICE_NAME}.service
    systemctl daemon-reload

    # Ask about data/config removal
    local remove_data=false
    if [ "${PURGE:-}" = "true" ]; then
        remove_data=true
    elif is_interactive; then
        read -p "$(msg 'purge_prompt')" -n 1 -r < /dev/tty
        echo
        [[ $REPLY =~ ^[Yy]$ ]] && remove_data=true
    fi

    if [ "$remove_data" = true ]; then
        print_info "$(msg 'removing_config_dir')"
        rm -rf "$INSTALL_DIR"
    else
        # Remove binary but keep data dir
        rm -f "$INSTALL_DIR/oa-nsdiy" "$INSTALL_DIR/.version" "$INSTALL_DIR/oa-nsdiy.backup"*
        print_warning "$(msg 'config_not_removed'): $INSTALL_DIR/data"
        print_warning "$(msg 'remove_manually'): rm -rf $INSTALL_DIR"
    fi

    print_info "$(msg 'removing_user')"
    userdel "$SERVICE_USER" 2>/dev/null || true

    print_success "$(msg 'uninstall_complete')"
}

# ============================================================
# Main
# ============================================================
main() {
    local target_version=""
    local positional_args=()

    while [[ $# -gt 0 ]]; do
        case "$1" in
            -y|--yes) FORCE_YES="true"; shift ;;
            --purge)  PURGE="true"; shift ;;
            -v|--version)
                if [ -n "${2:-}" ] && [[ ! "$2" =~ ^- ]]; then
                    target_version="$2"; shift 2
                else
                    echo "Error: --version requires a version argument"; exit 1
                fi
                ;;
            --version=*)
                target_version="${1#*=}"
                [ -z "$target_version" ] && { echo "Error: --version requires a version argument"; exit 1; }
                shift
                ;;
            *) positional_args+=("$1"); shift ;;
        esac
    done
    set -- "${positional_args[@]}"

    select_language

    echo ""
    echo "=============================================="
    echo "       $(msg 'install_title')"
    echo "=============================================="
    echo ""

    case "${1:-}" in
        upgrade|update)
            check_root; detect_platform; check_dependencies
            [ -n "$target_version" ] && install_version "$target_version" || upgrade
            exit 0 ;;
        install)
            check_root; detect_platform; check_dependencies
            if [ -n "$target_version" ] && [ -f "$INSTALL_DIR/oa-nsdiy" ]; then
                install_version "$target_version"
            else
                configure_server
                [ -n "$target_version" ] && { LATEST_VERSION=$(validate_version "$target_version"); fetch_release_by_version "$LATEST_VERSION"; } || get_latest_version
                download_and_extract
                create_user
                setup_directories
                generate_env
                install_service
                get_public_ip
                start_service
                enable_autostart
                print_completion
            fi
            exit 0 ;;
        rollback)
            [ -z "$target_version" ] && [ -n "${2:-}" ] && target_version="$2"
            [ -z "$target_version" ] && {
                print_error "$(msg 'opt_version')"
                echo ""; echo "Usage: $0 rollback -v <version>"; echo ""
                list_versions; exit 1
            }
            check_root; detect_platform; check_dependencies
            install_version "$target_version"
            exit 0 ;;
        list-versions|versions)
            list_versions; exit 0 ;;
        uninstall|remove)
            check_root; uninstall; exit 0 ;;
        --help|-h)
            echo "$(msg 'usage'): $0 [command] [options]"
            echo ""
            echo "Commands:"
            echo "  $(msg 'cmd_none')            $(msg 'cmd_install')"
            echo "  install              $(msg 'cmd_install')"
            echo "  upgrade              $(msg 'cmd_upgrade')"
            echo "  rollback <version>   $(msg 'cmd_install_version')"
            echo "  list-versions        $(msg 'cmd_list_versions')"
            echo "  uninstall            $(msg 'cmd_uninstall')"
            echo ""
            echo "Options:"
            echo "  -v, --version <ver>  $(msg 'opt_version')"
            echo "  -y, --yes            Skip confirmation prompts (for uninstall)"
            echo ""
            echo "Examples:"
            echo "  $0                        # Install latest version"
            echo "  $0 install -v v1.0.0      # Install specific version"
            echo "  $0 upgrade                # Upgrade to latest"
            echo "  $0 rollback v1.0.0        # Rollback to v1.0.0"
            echo "  $0 list-versions          # List available versions"
            echo ""
            exit 0 ;;
    esac

    # Default: fresh install with latest version
    check_root; detect_platform; check_dependencies
    configure_server
    [ -n "$target_version" ] && { LATEST_VERSION=$(validate_version "$target_version"); fetch_release_by_version "$LATEST_VERSION"; } || get_latest_version
    download_and_extract
    create_user
    setup_directories
    generate_env
    install_service
    get_public_ip
    start_service
    enable_autostart
    print_completion
}

main "$@"
