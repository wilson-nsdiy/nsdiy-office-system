# OA-NSDIY

工作室 OA 管理系统，采用 Go (Gin + Ent) 后端 + React (Next.js 14) 前端，前端静态导出嵌入 Go 二进制，**单一进程部署**。

| 仓库 | 地址 |
|------|------|
| Gitee | https://gitee.com/zhouws-chn/oa-nsdiy |
| GitHub | https://github.com/wilson-nsdiy/nsdiy-office-system |

## 技术栈

| 层级 | 技术 |
|------|------|
| 后端框架 | Gin |
| ORM | Ent v0.14.5 (类型安全查询 + Auto Migration) |
| 数据库 | SQLite (默认) / PostgreSQL |
| 依赖注入 | Google Wire |
| 配置管理 | Viper（环境变量优先，YAML 可选） |
| 日志 | Zap + Lumberjack (结构化日志 + 轮转) |
| 认证 | JWT 双 Token (golang-jwt) |
| 密码存储 | bcrypt + salt |
| 前端框架 | Next.js 14 (React 18, App Router, Static Export) |
| 样式 | Tailwind CSS |
| HTTP 客户端 | Axios |
| 包管理 | npm |

## 项目结构

```
oa-nsdiy/
├── backend/                    # Go 后端
│   ├── cmd/server/             # 程序入口 + Wire DI
│   ├── ent/                    # Ent 生成代码 (go generate)
│   │   └── schema/             # 数据库模型定义 (源文件)
│   ├── internal/
│   │   ├── config/             # 配置加载 (环境变量/YAML) + Wire 注入
│   │   ├── db/                 # 数据库连接 + Ent Client 初始化
│   │   ├── handler/            # HTTP 处理器
│   │   ├── middleware/         # JWT 认证、CORS 等
│   │   ├── pkg/                # 通用工具 (errors, response)
│   │   ├── repository/         # 数据访问层 (Ent Client 封装)
│   │   ├── server/             # 路由注册、中间件编排
│   │   │   └── routes/         # 按域分组的路由
│   │   ├── service/            # 业务逻辑层
│   │   └── web/                # 前端静态文件嵌入 (embed 构建标签)
├── frontend/                   # Next.js 前端 (Static Export)
│   ├── src/
│   │   ├── api/                # Axios API 客户端
│   │   ├── app/                # App Router 页面
│   │   ├── contexts/           # React Context (Auth)
│   │   └── types/              # TypeScript 类型定义
│   └── package.json
├── deploy/                     # 部署相关文件
│   ├── .env.example            # 环境变量配置模板
│   ├── build.sh                # 生产构建脚本 (Linux/macOS)
│   ├── build.bat               # 生产构建脚本 (Windows)
│   ├── install.sh              # 安装脚本 (Linux)
│   └── install.bat             # 安装脚本 (Windows)
```

## 功能模块

### 后端 API 模块

| 模块 | 路由前缀 | 功能 |
|------|----------|------|
| Auth | `/api/auth` | 登录、注册、密码重置、Token 刷新 |
| Setup | `/api/setup` | 系统初始化状态检查、创建初始管理员 |
| Roles | `/api/roles` | RBAC 角色管理 |
| Permissions | `/api/permissions` | RBAC 权限管理 (树形结构) |
| News Groups | `/api/news-groups` | 新闻分类管理 |
| News | `/api/news` | 新闻 CRUD |
| Articles | `/api/articles` | 文章管理，支持版本历史 |
| Projects | `/api/projects` | 项目管理，含成员分配 |
| Tasks | `/api/tasks` | 任务管理，支持子任务 |
| Media | `/api/media/accounts`, `/api/media/contents` | 社交媒体账号 + 内容管理 |
| Files | `/api/files` | 文件上传下载 |
| API Tokens | `/api/api-tokens` | API Token 管理 |

### 前端页面

| 路径 | 功能 |
|------|------|
| `/setup/` | 系统初始化 |
| `/login/` | 登录 |
| `/dashboard/` | 仪表盘 |
| `/articles/` | 文章管理 |
| `/news/` | 新闻管理 |
| `/news-groups/` | 新闻分组 |
| `/projects/` | 项目管理 |
| `/media/accounts/` | 媒体账号 |
| `/media/contents/` | 媒体内容 |
| `/roles/` | 角色管理 |
| `/permissions/` | 权限管理 |
| `/files/` | 文件管理 |
| `/api-tokens/` | API Token 管理 |

## 开发步骤

### 环境要求

- Go 1.25+
- Node.js 20+
- npm

### 1. 配置

应用所有字段都有默认值，可不提供任何配置直接启动。如需修改（如 JWT 密钥），二选一：

```bash
cp deploy/.env.example .env       # 由 shell/IDE 加载，或手动 export
```

### 2. 后端

```bash
cd backend

# 安装依赖
go mod download

# 生成 Ent 代码
go generate ./ent && go generate ./cmd/server

# 开发模式运行 (端口 3001)
go run ./cmd/server
```

### 3. 前端

```bash
cd frontend

# 安装依赖
npm install

# 开发模式运行 (端口 3000，API 请求代理到后端 3001)
npm run dev
```

开发模式下，前端运行在 `http://localhost:3000`，后端运行在 `http://localhost:3001`。
Next.js 开发服务器通过 rewrites 将 `/api/*` 请求代理到后端，无需跨域配置。

## 生产部署

详见 [生产模式部署指南](docs/production-deploy.md)。

## 常用命令

### 后端

| 命令 | 说明 |
|------|------|
| `go generate ./ent && go generate ./cmd/server` | 生成 Ent ORM 代码 |
| `go build -ldflags="-s -w -X main.Version=$(cat cmd/server/VERSION)" -trimpath -o bin/server ./cmd/server` | 构建 Go 后端 |
| `go run ./cmd/server` | 启动后端开发服务 |
| `go test ./...` | 运行所有测试 |
| `go test -tags=unit ./...` | 运行单元测试 |
| `golangci-lint run ./...` | golangci-lint 代码检查 |
| `touch migrations/$(date +%Y%m%d%H%M%S)_name.sql` | 创建新的迁移文件 |

### 前端

| 命令 | 说明 |
|------|------|
| `npm install` | 安装依赖 |
| `npm run dev` | 启动前端开发服务 |
| `npm run build` | 构建前端 (静态导出到 out/) |
| `npx tsc --noEmit` | TypeScript 类型检查 |

### 全局

| 命令 | 说明 |
|------|------|
| `go build -ldflags="-s -w -X main.Version=$(cat cmd/server/VERSION)" -trimpath -o bin/server ./cmd/server` | 构建后端 |
| `deploy/build.sh` | 生产构建 (Linux/macOS) |
| `deploy\\build.bat` | 生产构建 (Windows) |

## 配置说明

配置采用**环境变量优先、YAML 可选**的策略：应用经 viper 的 `AutomaticEnv()` 读取环境变量（key 中点号转下划线、自动大写），所有字段都有内置默认值，**无需任何配置文件即可启动**。

| 场景 | 推荐方式 | 说明 |
|------|----------|------|
| 生产部署 | `deploy/.env` | 从 `.env.example` 复制并填写 |
| 脚本安装 | `/opt/oa-nsdiy/.env` | 经 systemd 的 `EnvironmentFile` 注入（见方式一） |
| 开发环境 | 项目根目录 `.env` | 由 shell/IDE 加载；或直接 `export` 环境变量 |

**配置模板：**

```bash
# 生产部署（环境变量）
cp deploy/.env.example deploy/.env
cp deploy/.env.example .env             # 本地直跑时由 shell/IDE 加载
```

完整变量清单及说明见 `deploy/.env.example`（每个字段都标注了 `[必须修改]`/`[建议修改]`/`[保持默认]`）。

**关键配置项（环境变量形式）：**

```bash
# [必须修改] JWT 密钥 —— 生成命令: openssl rand -hex 32
JWT_SECRET=

# 服务
SERVER_HOST=0.0.0.0
SERVER_PORT=3001
SERVER_MODE=release              # release 或 debug

# 数据库: sqlite(默认) 或 postgres
DATABASE_DRIVER=sqlite
DATABASE_SOURCE=                 # 留空则自动推导为 ./data/db/oa_nsdiy.db
# PostgreSQL 示例:
# DATABASE_SOURCE=host=pg-host port=5432 user=pg-user password=pg-pass dbname=oa_nsdiy sslmode=require

# 日志
LOG_LEVEL=info
LOG_FORMAT=json
LOG_OUTPUT_TO_STDOUT=true

# CORS（单镜像同源部署时可设为 false 关闭）
CORS_ENABLED=true
CORS_ALLOW_ORIGINS=http://localhost:3000
CORS_ALLOW_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOW_HEADERS=Content-Type,Authorization
```

**生成 JWT 密钥：**

```bash
openssl rand -hex 32
```

## 架构说明

采用经典的分层架构：

```
Handler → Service → Repository → Ent Client → Database
```

- **Handler**: HTTP 请求处理、参数校验、响应格式化
- **Service**: 业务逻辑、事务编排、错误转换 (`HandleRepoErr` 统一 NotFound → 404)
- **Repository**: 数据访问、Ent Client 查询封装、实体类型转换
- **Middleware**: JWT 认证、CORS、日志

### 单一二进制架构

生产部署采用单一二进制架构，前端静态文件通过 Go `embed` 嵌入后端二进制：

```
┌─────────────────────────────────────────┐
│  单一 Go 二进制 (oa-nsdiy)              │
│                                         │
│  /api/*  → Gin 路由 (API)               │
│  其他    → 静态文件服务 (嵌入的前端)      │
│           SPA fallback → index.html     │
└─────────────────────────────────────────┘
```

- 构建标签 `embed` 控制是否嵌入前端（`embed_on.go` / `embed_off.go`）
- 开发模式：前后端独立运行，Next.js rewrites 代理 API
- 生产模式：单一进程同时服务 API 和前端静态文件

关键设计：
- Repository 返回 Ent 原始错误，Service 层通过 `ent.IsNotFound()` / `HandleRepoErr()` 转换为业务错误
- API 响应统一格式: 成功 `{ code: 0, message: "success", data: T }`, 失败 `{ code: <int>, message: string, reason: string }`
- 数据库变更通过 Ent Schema 定义 + `go generate ./ent` 生成类型安全代码

## 测试

后端测试使用内存 SQLite + Ent Client：

```bash
cd backend

# 运行所有测试
go test ./...

# 运行 Repository 层测试
go test ./internal/repository/ -v
```

## 数据库迁移

项目采用双重迁移策略：

1. **Ent Auto-Migration**: 开发阶段，启动时自动同步 Schema 到数据库
2. **版本化 SQL**: 生产部署，`migrations/` 目录下维护增量迁移文件

```bash
# 从 Ent Schema 生成迁移
go run -mod=mod entgo.io/ent/cmd/ent describe ./ent/schema

# 创建新的空迁移文件
touch migrations/$(date +%Y%m%d%H%M%S)_name.sql
```

## CI/CD

项目使用 GitHub Actions 流水线：

- **Release**: 推送标签触发，GoReleaser 自动构建并发布 Release