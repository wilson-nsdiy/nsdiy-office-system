# AGENTS.md

项目开发规范与约定，供 AI 辅助开发时参考。

## 项目概述

OA-NSDIY 是工作室 OA 管理系统，采用 Go (Gin + Ent) 后端 + Next.js 14 前端的分离架构。

## 技术栈

- **后端**: Go 1.25+, Gin, Ent ORM v0.14.5, Google Wire (DI), Viper (配置), JWT 双 Token 认证
- **前端**: Next.js 14 (App Router), React 18, TypeScript, Tailwind CSS, Axios
- **数据库**: SQLite (默认), 可切换 PostgreSQL
- **日志**: Zap + Lumberjack (结构化 + 轮转)
- **测试**: in-memory SQLite + Ent Client (Repository 层)
- **包管理**: npm (前端), Go Modules (后端)
- **CI/CD**: Gitee Go 流水线

## 项目结构

```
oa-nsdiy/
├── backend/                 # Go 后端
│   ├── cmd/server/          # 入口 + Wire DI
│   ├── ent/schema/          # Ent 数据模型 (源文件)
│   ├── ent/                 # Ent 生成代码 (go generate)
│   ├── migrations/          # 版本化数据库迁移
│   └── internal/
│       ├── config/          # 配置 + Wire
│       ├── db/              # 数据库连接 + Ent Client 初始化
│       ├── handler/         # HTTP 处理器
│       ├── middleware/      # 中间件 (JWT, CORS)
│       ├── pkg/             # 工具 (errors, response)
│       ├── repository/      # 数据访问层 (Ent Client 封装)
│       ├── server/routes/   # 路由注册
│       └── service/         # 业务逻辑
├── frontend/                # Next.js 前端
│   └── src/
│       ├── api/             # API 客户端
│       ├── app/             # App Router 页面
│       ├── contexts/        # React Context
│       └── types/           # TypeScript 类型
├── deploy/                  # Docker 部署配置
│   ├── Dockerfile           # 多阶段构建
│   ├── docker-compose.yml   # 容器编排
│   └── config.example.yaml  # 配置文件模板
└── .workflow/               # CI/CD 流水线
```

## 开发规范

### 后端

- 分层架构: Handler → Service → Repository → Ent Client
- Handler 只做参数校验和响应格式化，业务逻辑放 Service
- Repository 封装 Ent 查询，Service 不直接操作 Ent client
- Repository 返回 Ent 原始错误 (不做 sql.ErrNoRows 转换)，Service 用 `HandleRepoErr()` 统一处理
- 依赖注入通过 Google Wire 管理，在 `cmd/server` 和 `internal/config/wire.go` 中配置
- API 响应统一格式: 成功 `{ data: T }`, 失败 `{ error: string, code: string }`
- 使用 `internal/pkg/response` 统一响应，`internal/pkg/errors` 统一错误码
- 数据模型变更: 修改 `ent/schema/` 后执行 `make generate`
- 密码存储使用 bcrypt + salt
- JWT 双 Token 机制: Access Token (30m) + Refresh Token (168h)

### Repository 模式

- 构造函数: `NewXxxRepository(client *ent.Client) *XxxRepository`
- CRUD: `client.Xxx.Get(ctx, id)` / `Create().SetXxx().Save(ctx)` / `UpdateOneID(id).SetXxx().Save(ctx)` / `DeleteOneID(id).Exec(ctx)`
- 类型转换: `toXxx(e *ent.Xxx) *Xxx` 将 Ent 类型转为 Service 层实体
- 分页: `q.Count(ctx)` + `q.Order().Limit().Offset().All(ctx)`
- 关联加载: `WithXxx()` 替代 LEFT JOIN (如 `WithAuthor()`, `WithGroup()`)
- 错误处理: 直接返回 Ent 错误，不做 sql.ErrNoRows 转换

### Service 错误处理

- 使用 `HandleRepoErr(err, reason, message)` 将 Ent NotFound 转为 404 业务错误
- 外键验证的 NotFound 应返回 400 BadRequest (如 group_id 不存在)
- 其他 DB 错误直接透传为 500

### 前端

- 使用 App Router，页面为 `'use client'` 组件
- 认证状态通过 `contexts/auth.tsx` 的 React Context + localStorage 管理
- API 客户端在 `api/client.ts`，自动附加 Bearer Token，401 时自动刷新 Token
- 类型定义集中在 `types/index.ts`
- 表单中枚举类型字段 (如 ProjectStatus, ProjectPriority) 需显式类型标注，避免 TS 推断为 string
- select 的 onChange 需要 `as` 类型断言: `e.target.value as ProjectStatus`

### 通用

- **⚠️ 必须提供 `config.yaml` 配置文件**：程序本身不包含默认配置，开发和生产环境都需要
- 数据库文件不提交到 Git
- 前端使用 npm，不用 pnpm 或 yarn
- `package-lock.json` 需要提交到仓库

## 常用命令

```bash
# 后端
cd backend
go mod download          # 安装依赖
make generate            # 生成 Ent 代码
make dev                 # 开发模式运行
make build               # 构建 (开发模式，不含前端)
make test                # 运行测试
make lint                # golangci-lint 检查
make migrate-new         # 创建新的迁移文件

# 前端
cd frontend
npm install              # 安装依赖
npm run dev              # 开发模式运行
npm run build            # 构建 (静态导出到 out/)
npx tsc --noEmit         # 类型检查

# 生产构建 (前端嵌入后端)
cd frontend && npm run build
cp -r out/ ../backend/internal/web/dist
cd ../backend
go build -tags embed -o bin/server ./cmd/server

# Docker
make docker-up           # 启动服务
make docker-down         # 停止服务
```

## 构建模式

### 开发模式 (默认)
- 前后端独立运行
- `make build` 构建不含前端的后端二进制
- Next.js 开发服务器代理 API 请求

### 生产模式 (`-tags embed`)
- 前端静态文件通过 Go `embed` 嵌入后端二进制
- 构建步骤: 前端 `npm run build` → 复制到 `backend/internal/web/dist/` → `go build -tags embed`
- 单一进程同时服务 API 和前端
- 相关代码: `internal/web/embed_on.go` (启用) / `embed_off.go` (禁用)

## 配置

**⚠️ 必须提供 `config.yaml` 配置文件**：程序本身不包含默认配置。

| 场景 | 配置文件位置 | 生效方式 |
|------|-------------|----------|
| Docker 部署 | `deploy/config.yaml` | volume 挂载到 `/app/config.yaml` |
| 开发环境 | `config.yaml`（项目根目录）| Viper 自动搜索 |

后端配置文件 `config.yaml` (从 `config.example.yaml` 复制)，主要配置:

- `server`: 监听地址和端口 (默认 3001)、运行模式 (debug/release)
- `database`: 驱动和连接串（SQLite 留空自动推导）
- `jwt`: 密钥和过期时间
- `cors`: 跨域配置
- `log`: 日志级别、格式、输出方式

前端环境变量（仅开发模式需要）:

- `NEXT_PUBLIC_API_URL`: 前端 API 地址 (开发时默认 http://localhost:3001，生产同源部署无需配置)

## 经验教训

### 后端分层纪律

- Handler 禁止直接 import repository 或 db 包，必须通过 Service 层访问数据；使用 depguard lint 规则强制执行
- 当 Handler 需要查询数据（如当前用户信息），应在 Service 中提供方法（如 `GetUserByID`），Handler 调用 Service 获取 DTO

### Repository 错误处理

- Repository 层不要将 `ent.IsNotFound(err)` 转换为 `sql.ErrNoRows`，这会丢失语义信息
- Service 层使用 `HandleRepoErr()` 统一处理: NotFound → 404, 其他错误透传
- 外键校验的 NotFound 应返回 400 (如 "group_not_found")，而非 404

### Go 编码规范

- 类型断言必须使用双值形式 `v, ok := x.(T)`，避免 panic
- 资源关闭（rows.Close / tx.Rollback / DB.Close）必须检查返回错误
- Go 1.18+ 使用 `any` 替代 `interface{}`
- 条件分支可合并时遵循 staticcheck QF1007 建议

### Ent ORM 使用

- Schema 定义在 `ent/schema/`，通过 `go generate ./ent` 生成 Client 代码
- 生成的代码应提交到 Git (确定性输出)
- Task 的自引用关系在 Ent 中用纯字段 (parent_id) 而非自引用 Edge
- Repository 层的实体类型保持指针风格 (*string, *int)，通过转换器桥接 Ent 生成的值类型

### Token 刷新机制

- 并发刷新需用队列/锁防止多次请求，保证只刷新一次
- 主动刷新：Token 过期前（如 120 秒）提前刷新，减少用户感知的 401

### 前端通用模式

- 表格加载逻辑抽取为 `useTableLoader` hook，统一处理 loading/error/pagination
- 分页大小用 `usePersistedPageSize` 持久化到 localStorage
- 导航状态用 `useNavigationLoading` hook 管理
- 部署更新后 chunk 加载失败，用 `ChunkErrorRecovery` 组件自动 reload

### 日志

- 生产环境使用 zap 结构化日志 + lumberjack 日志轮转
- 前端 API 请求自动注入时区和语言头，便于后端日志关联
