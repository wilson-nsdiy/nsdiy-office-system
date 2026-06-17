# 生产模式部署指南

生产模式使用 Go `embed` 将前端静态文件嵌入后端二进制，构建单一可执行文件，同时服务 API 和前端。

## 部署方式

### 方式一：从 Release 页面下载

1. 访问项目的 Release 页面，下载对应平台的二进制文件
2. 上传至目标服务器并配置环境：

```bash
chmod +x server-linux-amd64
cp deploy/.env.example deploy/.env
# 编辑 .env 文件，配置必要的环境变量（至少设置 JWT_SECRET）
./server-linux-amd64
```

### 方式二：手动编译

```bash
# 1. 构建前端
cd frontend && npm install && npm run build

# 2. 复制前端文件到后端
cp -r frontend/out/ backend/internal/web/dist

# 3. 构建后端（-tags embed 启用前端嵌入，-ldflags="-s -w" 去除调试信息减小体积）
cd backend
CGO_ENABLED=0 go build -tags embed -ldflags="-s -w" -o bin/server ./cmd/server

# 4. 部署至目标服务器
scp bin/server user@server:/path/to/deploy/
```

目标服务器上同样需要配置 `.env` 并启动：

```bash
cp deploy/.env.example deploy/.env
# 编辑 .env 文件，配置必要的环境变量
./server
```

## 环境变量

从模板复制并填写：`cp deploy/.env.example deploy/.env`

| 变量 | 说明 | 默认值 |
|------|------|--------|
| SERVER_HOST | 监听地址 | 0.0.0.0 |
| SERVER_PORT | 监听端口 | 3001 |
| SERVER_MODE | 运行模式 | release |
| DATABASE_DRIVER | 数据库驱动 | sqlite |
| DATABASE_DSN | 数据库连接串 | (SQLite 自动推导) |
| JWT_SECRET | JWT 密钥 | **必须配置** |
| JWT_ACCESS_EXPIRY | Access Token 过期时间 | 30m |
| JWT_REFRESH_EXPIRY | Refresh Token 过期时间 | 168h |

## 注意事项

- SQLite 数据库文件默认在运行目录下，建议配置绝对路径
- 如需切换 PostgreSQL，设置 `DATABASE_DRIVER=postgresql` 和 `DATABASE_DSN`
