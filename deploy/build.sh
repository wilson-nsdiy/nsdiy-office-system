#!/bin/bash
# =============================================================================
# OA-NSDIY 生产构建脚本
# =============================================================================
# 一键构建生产环境可执行文件（前端嵌入后端）
#
# 用法:
#   ./deploy/build.sh           # 构建当前平台
#   ./deploy/build.sh linux     # 交叉编译 Linux amd64
#   ./deploy/build.sh darwin    # 交叉编译 macOS arm64
# =============================================================================

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

print_info()    { echo -e "${GREEN}[INFO]${NC} $1"; }
print_warn()    { echo -e "${YELLOW}[WARN]${NC} $1"; }
print_error()   { echo -e "${RED}[ERROR]${NC} $1"; exit 1; }

# Parse target platform
TARGET_OS="${1:-}"
VERSION=$(git -C "$PROJECT_ROOT" describe --tags --always --dirty 2>/dev/null || echo "dev")

case "$TARGET_OS" in
    linux)
        GOOS=linux GOARCH=amd64
        OUTPUT_NAME="oa-nsdiy-linux-amd64"
        ;;
    darwin|macos)
        GOOS=darwin GOARCH=arm64
        OUTPUT_NAME="oa-nsdiy-darwin-arm64"
        ;;
    "")
        GOOS=$(go env GOOS)
        GOARCH=$(go env GOARCH)
        OUTPUT_NAME="oa-nsdiy"
        ;;
    *)
        print_error "不支持的目标平台: $TARGET_OS (仅支持 linux/darwin)"
        ;;
esac

OUTPUT_DIR="$PROJECT_ROOT/backend/bin"

echo ""
echo "=========================================="
echo "  OA-NSDIY 生产构建"
echo "=========================================="
echo ""

# 1. Build frontend
print_info "构建前端..."
cd "$PROJECT_ROOT/frontend"
npm install --silent
npm run build
print_info "前端构建完成"

# 2. Copy frontend to embed directory
print_info "复制前端文件到后端..."
DIST_DIR="$PROJECT_ROOT/backend/internal/web/dist"
rm -rf "$DIST_DIR"/* 2>/dev/null || true
mkdir -p "$DIST_DIR"
cp -r "$PROJECT_ROOT/frontend/out/"* "$DIST_DIR/"
print_info "前端文件复制完成"

# 3. Build backend
print_info "构建后端 (GOOS=$GOOS GOARCH=$GOARCH)..."
cd "$PROJECT_ROOT/backend"
CGO_ENABLED=0 GOOS="$GOOS" GOARCH="$GOARCH" \
    go build -tags embed -ldflags="-s -w -X main.Version=$VERSION" \
    -trimpath -o "$OUTPUT_DIR/$OUTPUT_NAME" ./cmd/server

# 4. Done
BINARY_PATH="$OUTPUT_DIR/$OUTPUT_NAME"
FILE_SIZE=$(du -h "$BINARY_PATH" | cut -f1)

echo ""
echo "=========================================="
echo -e "  ${GREEN}构建成功!${NC}"
echo "=========================================="
echo ""
echo "  文件: $BINARY_PATH"
echo "  大小: $FILE_SIZE"
echo "  版本: $VERSION"
echo ""
echo "部署步骤:"
echo "  1. 上传 $OUTPUT_NAME 到目标服务器"
echo "  2. Linux: 使用安装脚本 curl -sSL https://raw.githubusercontent.com/wilson-nsdiy/nsdiy-office-system/master/deploy/install.sh | sudo bash"
echo "  3. 或手动部署: cp deploy/.env.example deploy/.env && ./$OUTPUT_NAME"
echo ""