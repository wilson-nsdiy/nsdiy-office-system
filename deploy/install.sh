#!/bin/bash
#
# OA-NSDIY 安装脚本
#
# Usage: curl -fsSL https://raw.githubusercontent.com/wilson-nsdiy/nsdiy-office-system/main/deploy/install.sh | sudo bash
#

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Configuration
GITHUB_OWNER="wilson-nsdiy"
GITHUB_REPO="nsdiy-office-system"
INSTALL_DIR="/opt/oa-nsdiy"
SERVICE_NAME="oa-nsdiy"
SERVICE_USER="oa-nsdiy"
CONFIG_DIR="/etc/oa-nsdiy"

SERVER_HOST="0.0.0.0"
SERVER_PORT="3001"

print_info()    { echo -e "${BLUE}[信息]${NC} $1"; }
print_success() { echo -e "${GREEN}[成功]${NC} $1"; }
print_warning() { echo -e "${YELLOW}[警告]${NC} $1"; }
print_error()   { echo -e "${RED}[错误]${NC} $1"; }

is_interactive() {
    [ -e /dev/tty ] && [ -r /dev/tty ] && [ -w /dev/tty ]
}

# ============================================================
# Root / platform / dependency checks
# ============================================================
check_root() {
    if [ "$(id -u)" -ne 0 ]; then
        print_error "请使用 root 权限运行 (使用 sudo)"
        exit 1
    fi
}

detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)

    case "$ARCH" in
        x86_64) ARCH="amd64" ;;
        aarch64|arm64) ARCH="arm64" ;;
        *) print_error "不支持的架构: $ARCH"; exit 1 ;;
    esac

    case "$OS" in
        linux) OS="linux" ;;
        *) print_error "不支持的操作系统: $OS"; exit 1 ;;
    esac

    print_info "检测到平台: ${OS}_${ARCH}"
}

check_dependencies() {
    local missing=()
    command -v curl &>/dev/null || missing+=("curl")
    command -v tar   &>/dev/null || missing+=("tar")
    [ ${#missing[@]} -gt 0 ] && {
        print_error "缺少依赖: ${missing[*]}"
        print_info "请先安装以上依赖"
        exit 1
    }
}

# ============================================================
# GitHub Release helpers
# ============================================================
API_BASE="https://api.github.com/repos/${GITHUB_OWNER}/${GITHUB_REPO}"

get_latest_version() {
    print_info "正在获取最新版本..."
    RELEASE_JSON=$(curl -s --connect-timeout 10 --max-time 30 "${API_BASE}/releases/latest" 2>/dev/null)

    if [ -z "$RELEASE_JSON" ]; then
        print_error "获取最新版本失败"
        print_info "请到 Release 页面手动下载: https://github.com/${GITHUB_OWNER}/${GITHUB_REPO}/releases"
        exit 1
    fi

    LATEST_VERSION=$(echo "$RELEASE_JSON" | grep -o '"tag_name"[[:space:]]*:[[:space:]]*"[^"]*"' | head -1 | sed -E 's/.*"([^"]+)"$/\1/')

    if [ -z "$LATEST_VERSION" ]; then
        print_error "获取最新版本失败"
        print_info "请到 Release 页面手动下载: https://github.com/${GITHUB_OWNER}/${GITHUB_REPO}/releases"
        exit 1
    fi

    print_info "最新版本: $LATEST_VERSION"
}

list_versions() {
    print_info "正在获取可用版本..."
    local versions
    versions=$(curl -s --connect-timeout 10 --max-time 30 "${API_BASE}/releases" 2>/dev/null | grep -o '"tag_name"[[:space:]]*:[[:space:]]*"[^"]*"' | sed -E 's/.*"([^"]+)"$/\1/' | head -20)

    if [ -z "$versions" ]; then
        print_error "获取版本列表失败"
        exit 1
    fi

    echo ""
    echo "可用版本列表:"
    echo "----------------------------------------"
    echo "$versions" | while read -r version; do echo "  $version"; done
    echo "----------------------------------------"
    echo ""
}

validate_version() {
    local version="$1"
    [ -z "$version" ] && { print_error "指定要安装的版本号 (例如: v1.0.0)"; exit 1; }

    [[ "$version" =~ ^v ]] || version="v$version"
    print_info "正在验证版本 $version"

    local http_code
    http_code=$(curl -s -o /dev/null -w "%{http_code}" --connect-timeout 10 --max-time 30 "${API_BASE}/releases/tags/${version}" 2>/dev/null)
    [ -z "$http_code" ] || ! [[ "$http_code" =~ ^[0-9]+$ ]] && {
        print_error "网络错误: 无法连接 GitHub API"; exit 1
    }
    [ "$http_code" != "200" ] && {
        print_error "指定版本不存在: $version"
        echo ""
        list_versions
        exit 1
    }
    echo "$version"
}

fetch_release_by_version() {
    local version="$1"
    RELEASE_JSON=$(curl -s --connect-timeout 10 --max-time 30 "${API_BASE}/releases/tags/${version}" 2>/dev/null)
}

get_asset_url() {
    local asset_name="$1"

    # GitHub API: assets[].browser_download_url contains the direct download link
    local url
    url=$(echo "$RELEASE_JSON" | python3 -c "
import sys, json
try:
    data = json.load(sys.stdin)
    for a in data.get('assets', []):
        if isinstance(a, dict) and a.get('name') == '${asset_name}':
            print(a.get('browser_download_url', '')); break
except Exception:
    pass
" 2>/dev/null)

    if [ -n "$url" ]; then
        echo "$url"
        return 0
    fi

    # Fallback: construct GitHub-style URL
    if [ -n "$LATEST_VERSION" ]; then
        echo "https://github.com/${GITHUB_OWNER}/${GITHUB_REPO}/releases/download/${LATEST_VERSION}/${asset_name}"
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

    print_info "正在下载 ${archive_name}..."

    TEMP_DIR=$(mktemp -d)
    trap "rm -rf $TEMP_DIR" EXIT

    local download_url
    download_url=$(get_asset_url "$archive_name")
    if [ -z "$download_url" ]; then
        download_url="https://github.com/${GITHUB_OWNER}/${GITHUB_REPO}/releases/download/${LATEST_VERSION}/${archive_name}"
    fi

    if ! curl -fsSL "$download_url" -o "$TEMP_DIR/$archive_name"; then
        print_error "下载失败"
        print_info "请到 Release 页面手动下载: https://github.com/${GITHUB_OWNER}/${GITHUB_REPO}/releases"
        exit 1
    fi

    print_info "正在校验文件..."
    local checksum_url
    checksum_url=$(get_asset_url "$checksum_asset")
    if [ -n "$checksum_url" ] && curl -fsSL "$checksum_url" -o "$TEMP_DIR/checksums.txt" 2>/dev/null; then
        local expected_checksum actual_checksum
        expected_checksum=$(grep "$archive_name" "$TEMP_DIR/checksums.txt" | awk '{print $1}')
        actual_checksum=$(sha256sum "$TEMP_DIR/$archive_name" | awk '{print $1}')
        if [ -n "$expected_checksum" ] && [ "$expected_checksum" != "$actual_checksum" ]; then
            print_error "校验失败"
            print_error "Expected: $expected_checksum"
            print_error "Actual:   $actual_checksum"
            exit 1
        fi
        [ -n "$expected_checksum" ] && print_success "校验通过" || print_warning "无法验证校验和（checksums.txt 未找到该文件）"
    else
        print_warning "无法验证校验和（checksums.txt 未找到）"
    fi

    print_info "正在解压..."
    tar -xzf "$TEMP_DIR/$archive_name" -C "$TEMP_DIR"

    mkdir -p "$INSTALL_DIR"

    if [ -f "$TEMP_DIR/oa-nsdiy" ]; then
        cp "$TEMP_DIR/oa-nsdiy" "$INSTALL_DIR/oa-nsdiy"
    elif [ -f "$TEMP_DIR/server" ]; then
        cp "$TEMP_DIR/server" "$INSTALL_DIR/oa-nsdiy"
    else
        print_error "压缩包中未找到二进制文件"
        exit 1
    fi
    chmod +x "$INSTALL_DIR/oa-nsdiy"

    echo "$LATEST_VERSION" > "$INSTALL_DIR/.version"

    print_success "二进制文件已安装到 $INSTALL_DIR/oa-nsdiy ($LATEST_VERSION)"
}

# ============================================================
# User / directories / .env
# ============================================================
create_user() {
    if id "$SERVICE_USER" &>/dev/null; then
        print_info "用户已存在: $SERVICE_USER"
        local current_shell
        current_shell=$(getent passwd "$SERVICE_USER" 2>/dev/null | cut -d: -f7)
        if [ "$current_shell" = "/bin/false" ] || [ "$current_shell" = "/sbin/nologin" ]; then
            usermod -s /bin/sh "$SERVICE_USER" 2>/dev/null && print_success "用户 shell 已更新为 /bin/sh"
        fi
    else
        print_info "正在创建系统用户 $SERVICE_USER..."
        useradd -r -s /bin/sh -d "$INSTALL_DIR" "$SERVICE_USER"
        print_success "用户已创建"
    fi
}

setup_directories() {
    print_info "正在设置目录..."
    mkdir -p "$INSTALL_DIR" "$INSTALL_DIR/data" "$CONFIG_DIR"
    chown -R "$SERVICE_USER:$SERVICE_USER" "$INSTALL_DIR" "$CONFIG_DIR"
    print_success "目录配置完成"
}

generate_env() {
    local env_file="$INSTALL_DIR/.env"
    [ -f "$env_file" ] && { print_info ".env 已存在，保留现有配置"; return; }

    print_info "正在生成 .env 配置..."
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
    print_success ".env 已生成（含随机 JWT_SECRET）"
}

# ============================================================
# systemd service
# ============================================================
install_service() {
    print_info "正在安装 systemd 服务..."

    cat > /etc/systemd/system/${SERVICE_NAME}.service << EOF
[Unit]
Description=OA-NSDIY - Studio OA Management System
Documentation=https://github.com/${GITHUB_OWNER}/${GITHUB_REPO}
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
    print_success "systemd 服务已安装"
}

# ============================================================
# Service management
# ============================================================
start_service() {
    print_info "正在启动服务..."
    if systemctl start "$SERVICE_NAME"; then
        print_success "服务已启动"
    else
        print_error "服务启动失败，请检查日志"
        print_info "sudo journalctl -u ${SERVICE_NAME} -n 50"
        return 1
    fi
}

enable_autostart() {
    print_info "正在设置开机自启..."
    systemctl enable "$SERVICE_NAME" 2>/dev/null && print_success "开机自启已启用" || print_warning "设置开机自启失败"
}

# ============================================================
# Completion message
# ============================================================
print_completion() {
    echo ""
    echo "=============================================="
    print_success "OA-NSDIY 安装完成！"
    echo "=============================================="
    echo ""
    echo "安装目录: $INSTALL_DIR"
    echo "服务地址: ${SERVER_HOST}:${SERVER_PORT}"
    echo ""
    echo "=============================================="
    echo "  后续步骤"
    echo "=============================================="
    echo ""
    echo "  1. 编辑配置文件（修改 JWT_SECRET 等必填项）:"
    echo "     sudo nano ${INSTALL_DIR}/.env"
    echo ""
    echo "  2. 启动服务:"
    echo "     sudo systemctl restart ${SERVICE_NAME}"
    echo ""
    echo "  3. 设置开机自启:"
    echo "     sudo systemctl enable ${SERVICE_NAME}"
    echo ""
    echo "  4. 访问 Web:"
    echo "     http://YOUR_SERVER_IP:${SERVER_PORT}"
    echo ""
    echo "=============================================="
    echo "  常用命令"
    echo "=============================================="
    echo ""
    echo "  查看状态: sudo systemctl status ${SERVICE_NAME}"
    echo "  查看日志: sudo journalctl -u ${SERVICE_NAME} -f"
    echo "  重启服务: sudo systemctl restart ${SERVICE_NAME}"
    echo "  停止服务: sudo systemctl stop ${SERVICE_NAME}"
    echo ""
}

# ============================================================
# Upgrade / install-version / uninstall
# ============================================================
upgrade() {
    [ ! -f "$INSTALL_DIR/oa-nsdiy" ] && {
        print_error "OA-NSDIY 尚未安装，请先执行全新安装"
        print_info "用法: $0 install"
        exit 1
    }

    print_info "正在升级 OA-NSDIY..."

    local current_version
    current_version=$(get_current_version)
    print_info "当前版本: $current_version"

    systemctl is-active --quiet "$SERVICE_NAME" && {
        print_info "正在停止服务..."
        systemctl stop "$SERVICE_NAME"
    }

    cp "$INSTALL_DIR/oa-nsdiy" "$INSTALL_DIR/oa-nsdiy.backup"
    print_info "备份已创建: $INSTALL_DIR/oa-nsdiy.backup"

    get_latest_version
    download_and_extract
    chown "$SERVICE_USER:$SERVICE_USER" "$INSTALL_DIR/oa-nsdiy"

    print_info "正在启动服务..."
    systemctl start "$SERVICE_NAME"
    print_success "升级完成！"
}

install_version() {
    local target_version="$1"

    [ ! -f "$INSTALL_DIR/oa-nsdiy" ] && {
        print_error "OA-NSDIY 尚未安装，请先执行全新安装"
        print_info "用法: $0 install -v $target_version"
        exit 1
    }

    target_version=$(validate_version "$target_version")
    print_info "正在安装指定版本: $target_version"

    local current_version
    current_version=$(get_current_version)
    print_info "当前版本: $current_version"

    if [ "$current_version" = "$target_version" ] || [ "$current_version" = "${target_version#v}" ]; then
        print_warning "已经是该版本，无需操作"
        exit 0
    fi

    systemctl is-active --quiet "$SERVICE_NAME" && {
        print_info "正在停止服务..."
        systemctl stop "$SERVICE_NAME"
    }

    local backup_name
    if [ "$current_version" != "unknown" ] && [ "$current_version" != "not_installed" ]; then
        backup_name="oa-nsdiy.backup.${current_version}"
    else
        backup_name="oa-nsdiy.backup.$(date +%Y%m%d%H%M%S)"
    fi
    cp "$INSTALL_DIR/oa-nsdiy" "$INSTALL_DIR/$backup_name"
    print_info "备份已创建: $INSTALL_DIR/$backup_name"

    LATEST_VERSION="$target_version"
    fetch_release_by_version "$target_version"
    download_and_extract
    chown "$SERVICE_USER:$SERVICE_USER" "$INSTALL_DIR/oa-nsdiy"

    print_info "正在启动服务..."
    systemctl start "$SERVICE_NAME" && print_success "服务已启动" || {
        print_error "服务启动失败，请检查日志"
        print_info "sudo journalctl -u ${SERVICE_NAME} -n 50"
    }

    local new_version
    new_version=$(get_current_version)
    echo ""
    echo "=============================================="
    print_success "指定版本安装完成！"
    echo "=============================================="
    echo ""
    echo "  当前版本: $new_version"
    echo ""
}

uninstall() {
    print_warning "这将从系统中移除 OA-NSDIY。"

    if ! is_interactive; then
        if [ "${FORCE_YES:-}" != "true" ]; then
            print_error "非交互模式。请使用 'curl ... | bash -s -- uninstall -y' 确认卸载。"
            exit 1
        fi
    else
        read -p "确定要继续吗？(y/N) " -n 1 -r < /dev/tty
        echo
        [[ $REPLY =~ ^[Yy]$ ]] || { print_info "卸载已取消"; exit 0; }
    fi

    print_info "正在停止服务..."
    systemctl stop "$SERVICE_NAME" 2>/dev/null || true
    systemctl disable "$SERVICE_NAME" 2>/dev/null || true

    print_info "正在移除文件..."
    rm -f /etc/systemd/system/${SERVICE_NAME}.service
    systemctl daemon-reload

    local remove_data=false
    if [ "${PURGE:-}" = "true" ]; then
        remove_data=true
    elif is_interactive; then
        read -p "是否同时删除数据目录？这将清除所有数据 [y/N]: " -n 1 -r < /dev/tty
        echo
        [[ $REPLY =~ ^[Yy]$ ]] && remove_data=true
    fi

    if [ "$remove_data" = true ]; then
        print_info "正在移除数据目录..."
        rm -rf "$INSTALL_DIR"
    else
        rm -f "$INSTALL_DIR/oa-nsdiy" "$INSTALL_DIR/.version" "$INSTALL_DIR/oa-nsdiy.backup"*
        print_warning "数据目录未被移除: $INSTALL_DIR/data"
        print_warning "如不再需要，请手动删除: rm -rf $INSTALL_DIR"
    fi

    print_info "正在移除用户..."
    userdel "$SERVICE_USER" 2>/dev/null || true

    print_success "OA-NSDIY 已卸载"
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

    echo ""
    echo "=============================================="
    echo "       OA-NSDIY 安装脚本"
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
                [ -n "$target_version" ] && { LATEST_VERSION=$(validate_version "$target_version"); fetch_release_by_version "$LATEST_VERSION"; } || get_latest_version
                download_and_extract
                create_user
                setup_directories
                generate_env
                install_service
                start_service
                enable_autostart
                print_completion
            fi
            exit 0 ;;
        rollback)
            [ -z "$target_version" ] && [ -n "${2:-}" ] && target_version="$2"
            [ -z "$target_version" ] && {
                print_error "指定要安装的版本号 (例如: v1.0.0)"
                echo ""; echo "用法: $0 rollback -v <version>"; echo ""
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
            echo "用法: $0 [命令] [选项]"
            echo ""
            echo "命令:"
            echo "  (无参数)             安装最新版本"
            echo "  install              安装 OA-NSDIY"
            echo "  upgrade              升级到最新版本"
            echo "  rollback <版本>       安装/回退到指定版本"
            echo "  list-versions        列出可用版本"
            echo "  uninstall            卸载 OA-NSDIY"
            echo ""
            echo "选项:"
            echo "  -v, --version <版本>  指定要安装的版本号 (例如: v1.0.0)"
            echo "  -y, --yes            跳过确认提示 (用于卸载)"
            echo ""
            echo "示例:"
            echo "  $0                        # 安装最新版本"
            echo "  $0 install -v v1.0.0      # 安装指定版本"
            echo "  $0 upgrade                # 升级到最新"
            echo "  $0 rollback v1.0.0        # 回退到 v1.0.0"
            echo "  $0 list-versions          # 列出可用版本"
            echo ""
            exit 0 ;;
    esac

    # Default: fresh install with latest version
    check_root; detect_platform; check_dependencies
    [ -n "$target_version" ] && { LATEST_VERSION=$(validate_version "$target_version"); fetch_release_by_version "$LATEST_VERSION"; } || get_latest_version
    download_and_extract
    create_user
    setup_directories
    generate_env
    install_service
    start_service
    enable_autostart
    print_completion
}

main "$@"
