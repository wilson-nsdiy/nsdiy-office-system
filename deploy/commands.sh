#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

VERSION=$(git -C "$PROJECT_ROOT" describe --tags --always --dirty 2>/dev/null || echo "dev")

usage() {
    echo "Usage: $0 <command>"
    echo ""
    echo "Commands:"
    echo "  build           - Build frontend + backend (without embedding frontend)"
    echo "  build-embed     - Build frontend and embed into backend binary (production)"
    echo "  build-backend   - Build Go backend (without embed)"
    echo "  build-frontend  - Build Next.js frontend (static export)"
    echo "  test            - Run all tests"
    echo "  test-backend    - Run backend tests"
    echo "  test-frontend   - Run frontend tests"
    echo "  lint            - Run linters"
    echo "  dev-backend     - Run backend in development mode"
    echo "  dev-frontend    - Run frontend in development mode"
    echo "  docker-up       - Start all services with Docker Compose"
    echo "  docker-down     - Stop all services"
    echo "  clean           - Clean build artifacts"
    echo "  generate        - Generate Ent code"
    echo "  migrate-new     - Create a new migration file"
}

cmd_build_backend() {
    cd "$PROJECT_ROOT/backend"
    CGO_ENABLED=0 go build -ldflags="-s -w -X main.Version=$VERSION" -trimpath -o bin/server ./cmd/server
}

cmd_build_frontend() {
    cd "$PROJECT_ROOT/frontend"
    npm run build
}

cmd_build() {
    cmd_build_frontend
    cmd_build_backend
}

cmd_build_embed() {
    cmd_build_frontend
    echo "Copying frontend output to backend embed directory..."
    cd "$PROJECT_ROOT/backend/internal/web/dist"
    rm -f * 2>/dev/null || true
    cp -r "$PROJECT_ROOT/frontend/out/"* .
    cmd_build_backend
    echo "Embedded build complete."
}

cmd_test_backend() {
    cd "$PROJECT_ROOT/backend"
    go test ./...
}

cmd_test_frontend() {
    cd "$PROJECT_ROOT/frontend"
    npm run lint:check
    npx tsc --noEmit
}

cmd_test() {
    cmd_test_backend
    cmd_test_frontend
}

cmd_lint() {
    cd "$PROJECT_ROOT/backend"
    golangci-lint run ./...
}

cmd_dev_backend() {
    cd "$PROJECT_ROOT/backend"
    go run ./cmd/server
}

cmd_dev_frontend() {
    cd "$PROJECT_ROOT/frontend"
    npm run dev
}

cmd_docker_up() {
    cd "$SCRIPT_DIR"
    docker-compose up -d
}

cmd_docker_down() {
    cd "$SCRIPT_DIR"
    docker-compose down
}

cmd_clean() {
    rm -rf "$PROJECT_ROOT/backend/bin"
    rm -rf "$PROJECT_ROOT/frontend/.next"
}

cmd_generate() {
    cd "$PROJECT_ROOT/backend"
    go generate ./ent
    go generate ./cmd/server
}

cmd_migrate_new() {
    read -p "Enter migration name: " name
    local ts
    ts=$(date +%Y%m%d%H%M%S)
    local file="$PROJECT_ROOT/backend/migrations/${ts}_${name}.sql"
    touch "$file"
    echo "Created $file"
}

if [[ $# -lt 1 ]]; then
    usage
    exit 1
fi

command="$1"
shift

case "$command" in
    build)          cmd_build ;;
    build-embed)    cmd_build_embed ;;
    build-backend)  cmd_build_backend ;;
    build-frontend) cmd_build_frontend ;;
    test)           cmd_test ;;
    test-backend)   cmd_test_backend ;;
    test-frontend)  cmd_test_frontend ;;
    lint)           cmd_lint ;;
    dev-backend)    cmd_dev_backend ;;
    dev-frontend)   cmd_dev_frontend ;;
    docker-up)      cmd_docker_up ;;
    docker-down)    cmd_docker_down ;;
    clean)          cmd_clean ;;
    generate)       cmd_generate ;;
    migrate-new)    cmd_migrate_new ;;
    help|--help|-h) usage ;;
    *)
        echo "Unknown command: $command"
        usage
        exit 1
        ;;
esac
