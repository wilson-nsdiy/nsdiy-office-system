#!/bin/bash
# =============================================================================
# OA-NSDIY Docker 部署准备脚本
# =============================================================================
# 本脚本会在 deploy/ 目录下生成部署所需的文件:
#   - 由 .env.example 复制出 .env
#   - 自动生成 JWT_SECRET 随机密钥
#   - 创建本地数据目录 data/
#
# 运行后请编辑 .env 填写必填项 (尤其 IMAGE_REGISTRY)，然后:
#   docker-compose up -d
# =============================================================================

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_info()    { echo -e "${BLUE}[INFO]${NC} $1"; }
print_success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }
print_warning() { echo -e "${YELLOW}[WARNING]${NC} $1"; }
print_error()   { echo -e "${RED}[ERROR]${NC} $1"; }

# Generate random secret (32 bytes hex)
generate_secret() {
    openssl rand -hex 32
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

main() {
    echo ""
    echo "=========================================="
    echo "  OA-NSDIY 部署准备"
    echo "=========================================="
    echo ""

    # 1. 依赖检查
    if ! command_exists openssl; then
        print_error "未安装 openssl，请先安装。"
        exit 1
    fi

    # 2. 确认脚本所在目录 (deploy/)
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    cd "$SCRIPT_DIR"
    if [ ! -f ".env.example" ]; then
        print_error "未找到 .env.example，请确认在 deploy/ 目录下运行。"
        exit 1
    fi

    # 3. 已存在则确认覆盖
    if [ -f ".env" ]; then
        print_warning "当前目录已存在 .env"
        read -p "是否覆盖？(y/N): " -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            print_info "已取消。"
            exit 0
        fi
    fi

    # 4. 由模板生成 .env
    print_info "正在生成 .env ..."
    cp .env.example .env

    # 5. 生成并写入 JWT_SECRET
    print_info "正在生成随机密钥 ..."
    JWT_SECRET=$(generate_secret)

    # 跨平台 sed: GNU (Linux) 用 -i，BSD (macOS) 用 -i ''
    if sed --version >/dev/null 2>&1; then
        sed -i "s/^JWT_SECRET=.*/JWT_SECRET=${JWT_SECRET}/" .env
    else
        sed -i '' "s/^JWT_SECRET=.*/JWT_SECRET=${JWT_SECRET}/" .env
    fi
    print_success "已生成 .env 并写入随机密钥"

    # 6. 创建数据目录
    print_info "正在创建数据目录 ..."
    mkdir -p data
    print_success "已创建 data/"

    # 7. 收紧 .env 权限 (仅所有者可读写)
    chmod 600 .env 2>/dev/null || true

    # 8. 完成提示
    echo ""
    echo "=========================================="
    echo "  准备完成！"
    echo "=========================================="
    echo ""
    echo "已生成的密钥 (请妥善保管，切勿公开):"
    echo "  JWT_SECRET:  ${JWT_SECRET}"
    echo ""
    print_warning "密钥已保存到 .env，请勿提交到 git 或对外分享。"
    echo ""
    echo "目录结构:"
    echo "  docker-compose.yml   - Docker Compose 编排"
    echo "  .env                 - 环境变量 (含生成的密钥)"
    echo "  .env.example         - 模板 (仅作参考)"
    echo "  data/                - 应用数据 (SQLite/日志，运行时生成)"
    echo ""
    echo "下一步:"
    echo "  1. 编辑 .env，填写以下必填项:"
    echo "     - IMAGE_REGISTRY      (Gitee 容器仓库地址)"
    echo "     - DATABASE_SOURCE     (若使用 PostgreSQL)"
    echo "     - CORS_ALLOW_ORIGINS  (真实前端域名)"
    echo ""
    echo "  2. 启动服务:"
    echo "     docker-compose up -d"
    echo ""
    echo "  3. 查看日志:"
    echo "     docker-compose logs -f oa-nsdiy"
    echo ""
    echo "  4. 访问 Web:"
    echo "     http://localhost:3001   (默认账号 admin / admin123)"
    echo ""
}

main "$@"
