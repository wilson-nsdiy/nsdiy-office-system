# OA-NSDIY

工作室 OA 管理系统，采用 Go (Gin + Ent) 后端 + React (Next.js 14) 前端，前端静态导出嵌入 Go 二进制，**单镜像部署**。

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
| 容器化 | Docker (单镜像多阶段构建) |

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
│   ├── Dockerfile
│   └── Makefile
├── frontend/                   # Next.js 前端 (Static Export)
│   ├── src/
│   │   ├── api/                # Axios API 客户端
│   │   ├── app/                # App Router 页面
│   │   ├── contexts/           # React Context (Auth)
│   │   └── types/              # TypeScript 类型定义
│   └── package.json
├── deploy/                     # 部署相关文件
│   ├── Dockerfile              # 多阶段构建 (前端构建 → Go embed → Alpine)
│   ├── docker-compose.yml      # Docker Compose 编排 (远程镜像 + .env)
│   ├── .env.example            # 环境变量配置模板（推荐，部署时复制为 .env）
│   ├── docker-deploy.sh        # Docker 部署准备脚本 (生成 .env + 密钥)
│   ├── install.sh              # 脚本安装 (裸机/systemd，从 Gitee Release 拉取)
│   ├── commands.sh             # Linux/macOS 构建脚本
│   └── commands.ps1            # Windows 构建脚本
```

## 功能模块

### 后端 API 模块

| 模块 | 路由前缀 | 功能 |
|------|----------|------|
| Auth | `/api/auth` | 登录、注册、密码重置、Token 刷新 |
| Roles | `/api/roles` | RBAC 角色管理 |
| Permissions | `/api/permissions` | RBAC 权限管理 (树形结构) |
| News Groups | `/api/news-groups` | 新闻分类管理 |
| News | `/api/news` | 新闻 CRUD |
| Articles | `/api/articles` | 文章管理，支持版本历史 |
| Projects | `/api/projects` | 项目管理，含成员分配 |
| Tasks | `/api/tasks` | 任务管理，支持子任务 |
| Media | `/api/media` | 社交媒体账号 + 内容管理 |
| Files | `/api/files` | 文件上传下载 |
| API Tokens | `/api/api-tokens` | API Token 管理 |

### 前端页面

| 路径 | 功能 |
|------|------|
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

- Go 1.23+
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
make generate

# 开发模式运行 (端口 3001)
make dev
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

### 方式一：脚本安装（推荐）

一键安装脚本，从 Gitee Releases 下载预编译二进制（已嵌入前端），自动配置 systemd 服务。**无需 Docker、无需源码、无需编译环境**，适合裸机/VPS 直接部署。

#### 一键安装

```bash
curl -fsSL https://gitee.com/zhouws-chn/oa-nsdiy/raw/master/deploy/install.sh | sudo bash
```

脚本会交互式询问：
- 语言（中文/English）
- 监听地址、端口（默认 `0.0.0.0:3001`）

然后自动完成：
1. 下载最新版本二进制（含 sha256 校验）→ `/opt/oa-nsdiy/server`
2. 创建系统用户 `oa-nsdiy`
3. 生成 `.env`（含随机 `JWT_SECRET`）→ `/opt/oa-nsdiy/.env`
4. 注册并启动 systemd 服务（开机自启）

#### 安装后步骤

```bash
# 1. 编辑配置（修改 JWT_SECRET、CORS_ENABLED、CORS_ALLOW_ORIGINS 等）
sudo nano /opt/oa-nsdiy/.env

# 2. 重启服务使配置生效
sudo systemctl restart oa-nsdiy

# 3. 访问
#    http://<服务器IP>:3001  （默认账号 admin / admin123）
```

#### 常用命令

| 命令 | 说明 |
|------|------|
| `sudo systemctl status oa-nsdiy` | 查看状态 |
| `sudo journalctl -u oa-nsdiy -f` | 实时日志 |
| `sudo systemctl restart oa-nsdiy` | 重启 |
| `sudo systemctl stop oa-nsdiy` | 停止 |

#### 升级 / 回滚 / 卸载

```bash
# 升级到最新版本
curl -fsSL https://gitee.com/zhouws-chn/oa-nsdiy/raw/master/deploy/install.sh | sudo bash -s -- upgrade

# 列出可用版本
curl -fsSL https://gitee.com/zhouws-chn/oa-nsdiy/raw/master/deploy/install.sh | sudo bash -s -- list-versions

# 安装/回退到指定版本
curl -fsSL https://gitee.com/zhouws-chn/oa-nsdiy/raw/master/deploy/install.sh | sudo bash -s -- install -v v1.0.0

# 卸载（保留数据；加 --purge 连数据一起删）
curl -fsSL https://gitee.com/zhouws-chn/oa-nsdiy/raw/master/deploy/install.sh | sudo bash -s -- uninstall
```

> **前置条件**：Gitee Releases 页面需有对应版本的 `oa-nsdiy_<version>_linux_{amd64,arm64}.tar.gz` 制品。CI 流水线（`release-pipeline.yml`）在推送 `v*` tag 时自动构建，需手动上传到 Release 页面（详见 `.workflow/release-pipeline.yml` 注释）。

#### 配置说明

脚本安装方式通过 systemd 的 `EnvironmentFile=/opt/oa-nsdiy/.env` 注入环境变量，应用经 viper.AutomaticEnv 读取。所有可用变量及说明见 `deploy/.env.example`。

---

### 方式二：Docker 部署

单镜像架构：前端静态导出嵌入 Go 二进制，只需一个容器。配置通过 `.env` 环境变量注入，数据持久化到本地 `./data` 目录。

#### 1. 准备配置

```bash
cd deploy

# 一键生成 .env（自动填充随机 JWT_SECRET），或手动 cp .env.example .env
bash docker-deploy.sh

# 编辑 .env，必须填写：
#   JWT_SECRET           （随机密钥，生成: openssl rand -hex 32）
#   CORS_ENABLED         （单镜像同源部署时可设为 false）
#   CORS_ALLOW_ORIGINS   （真实前端域名；单机同源可保持默认）
nano .env
```

#### 2. 拉取镜像并启动

```bash
docker-compose up -d
```

`docker-compose.yml` 默认走远程镜像拉取模式；若需本地构建（开发/无远程仓库时），将其中 `image:` 行注释、取消 `build:` 块的注释即可。

#### 3. 验证

```bash
# 查看日志
docker-compose logs -f oa-nsdiy

# 健康检查
curl http://localhost:3001/api/health
```

访问 `http://localhost:3001`（默认账号 admin / admin123）。

#### 配置与数据

| 项目 | 位置 | 说明 |
|------|------|------|
| 配置 | `deploy/.env` | 由 `env_file` 注入容器，应用经 viper.AutomaticEnv 读取 |
| 数据 | `deploy/data/` | SQLite 数据库 + 日志，本地目录映射便于整体打包迁移 |

**Docker 镜像构建流程（Dockerfile 多阶段构建）：**

1. **Stage 1** - Node.js 构建前端 (`npm run build` → `out/` 静态文件)
2. **Stage 2** - Go 编译后端 (`-tags embed` 将前端嵌入二进制)
3. **Stage 3** - Alpine 运行时 (仅包含单一二进制，约 50MB)

#### 整体迁移

`.env` 和数据都在本地目录，可整体打包迁移到新服务器：

```bash
tar czf oa-nsdiy-backup.tar.gz docker-compose.yml .env data/
# 在新服务器解压后 docker-compose up -d 即可
```

### 方式三：手动构建

```bash
# 1. 构建前端
cd frontend
npm install
npm run build    # 输出到 out/ 目录

# 2. 将前端产物复制到后端 embed 目录
cp -r out/ ../backend/internal/web/dist

# 3. 构建后端 (带 embed 标签)
cd ../backend
CGO_ENABLED=0 go build -tags embed -ldflags="-s -w" -o bin/server ./cmd/server

# 4. 运行
./bin/server
```

访问 `http://localhost:3001`，前端和 API 由同一进程提供服务。

## 常用命令

### 后端

| 命令 | 说明 |
|------|------|
| `make generate` | 生成 Ent ORM 代码 |
| `make build` | 构建 Go 后端 |
| `make dev` | 启动后端开发服务 |
| `make test` | 运行所有测试 |
| `make test-unit` | 运行单元测试 |
| `make lint` | golangci-lint 代码检查 |
| `make migrate-new` | 创建新的迁移文件 |
| `make clean` | 清理构建产物 |

### 前端

| 命令 | 说明 |
|------|------|
| `npm install` | 安装依赖 |
| `npm run dev` | 启动前端开发服务 |
| `npm run build` | 构建前端 (静态导出到 out/) |
| `npx tsc --noEmit` | TypeScript 类型检查 |

### 全局

Linux/macOS 使用 `./deploy/commands.sh`，Windows 使用 `.\deploy\commands.ps1`：

| 命令 | 说明 |
|------|------|
| `build` | 构建后端 + 前端 |
| `build-embed` | 构建前端并嵌入后端二进制 (生产) |
| `dev-backend` | 启动后端开发服务 |
| `dev-frontend` | 启动前端开发服务 |
| `test` | 运行所有测试 |
| `lint` | 代码检查 |
| `docker-up` | Docker Compose 启动 |
| `docker-down` | Docker Compose 停止 |
| `clean` | 清理构建产物 |
| `generate` | 生成 Ent 代码 |
| `migrate-new` | 创建新的迁移文件 |

### 前端

| 命令 | 说明 |
|------|------|
| `npm install` | 安装依赖 |
| `npm run dev` | 启动前端开发服务 |
| `npm run build` | 构建前端 (静态导出到 out/) |
| `npx tsc --noEmit` | TypeScript 类型检查 |

### 全局

Linux/macOS 使用 `./deploy/commands.sh`，Windows 使用 `.\deploy\commands.ps1`：

| 命令 | 说明 |
|------|------|
| `build` | 构建后端 + 前端 |
| `dev-backend` | 启动后端开发服务 |
| `dev-frontend` | 启动前端开发服务 |
| `docker-up` | Docker Compose 启动 |
| `docker-down` | Docker Compose 停止 |
| `clean` | 清理构建产物 |

## 配置说明

配置采用**环境变量优先、YAML 可选**的策略：应用经 viper 的 `AutomaticEnv()` 读取环境变量（key 中点号转下划线、自动大写），所有字段都有内置默认值，**无需任何配置文件即可启动**。

| 场景 | 推荐方式 | 说明 |
|------|----------|------|
| Docker 部署 | `deploy/.env` | 经 docker-compose 的 `env_file` 注入容器（见方式二） |
| 脚本安装 | `/opt/oa-nsdiy/.env` | 经 systemd 的 `EnvironmentFile` 注入（见方式一） |
| 开发环境 | 项目根目录 `.env` | 由 shell/IDE 加载；或直接 `export` 环境变量 |

**配置模板：**

```bash
# Docker / 脚本部署（环境变量）
cp deploy/.env.example deploy/.env      # Docker
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

### 单镜像架构

生产部署采用单镜像架构，前端静态文件通过 Go `embed` 嵌入后端二进制：

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
- API 响应统一格式: 成功 `{ data: T }`, 失败 `{ error: string, code: string }`
- 数据库变更通过 Ent Schema 定义 + `make generate` 生成类型安全代码

## 测试

后端测试使用内存 SQLite + Ent Client：

```bash
cd backend

# 运行所有测试
make test

# 运行 Repository 层测试
go test ./internal/repository/ -v
```

## 数据库迁移

项目采用双重迁移策略：

1. **Ent Auto-Migration**: 开发阶段，启动时自动同步 Schema 到数据库
2. **版本化 SQL**: 生产部署，`migrations/` 目录下维护增量迁移文件

```bash
# 从 Ent Schema 生成迁移
make migrate-diff

# 创建新的空迁移文件
make migrate-new
```

## CI/CD

项目使用 Gitee Go 流水线：

- **master-pipeline**: master 分支推送触发，编译 + 构建镜像 + 发布
- **branch-pipeline**: 非 master 分支推送触发，编译检查
- **pr-pipeline**: PR 到 master 触发，编译检查