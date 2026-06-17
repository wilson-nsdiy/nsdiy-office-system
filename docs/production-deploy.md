# 生产模式部署指南

## 概述

生产模式使用 Go `embed` 将前端静态文件嵌入后端二进制，构建单一可执行文件，同时服务 API 和前端。

## 部署方式

### 方式一：从 Release 页面下载成果物部署

1. 访问项目的 Release 页面
2. 下载对应版本的二进制文件（如 `server-linux-amd64`、`server-darwin-arm64` 等）
3. 上传至目标服务器

```bash
# 赋予执行权限
chmod +x server-linux-amd64

# 配置环境
cp deploy/.env.example deploy/.env
# 编辑 .env 文件，配置必要的环境变量

# 启动服务
./server-linux-amd64
```

### 方式二：手动编译获取成果物后部署

#### 1. 构建前端

```bash
cd frontend
npm install          # 安装前端依赖
npm run build        # 构建静态文件，输出到 out/
```

#### 2. 复制前端文件到后端

```bash
cp -r frontend/out/ backend/internal/web/dist
```

#### 3. 构建后端二进制

```bash
cd backend
CGO_ENABLED=0 go build -tags embed -ldflags="-s -w" -o bin/server ./cmd/server
```

**参数说明：**
- `-tags embed`：启用 embed 模式，嵌入前端静态文件
- `-ldflags="-s -w"`：去除调试信息，减小二进制体积
- `-o bin/server`：输出文件路径

#### 4. 部署至目标服务器

```bash
# 上传 bin/server 到目标服务器
scp bin/server user@server:/path/to/deploy/

# 在目标服务器上配置环境
cp deploy/.env.example deploy/.env
# 编辑 .env 文件，配置必要的环境变量

# 启动服务
./server
```

## 环境变量配置

环境变量通过 `.env` 文件配置。部署时从模板复制并填写：

```bash
cp deploy/.env.example deploy/.env
# 编辑 .env 填写必要配置
```

| 变量 | 说明 | 默认值 |
|------|------|--------|
| SERVER_HOST | 监听地址 | 0.0.0.0 |
| SERVER_PORT | 监听端口 | 3001 |
| SERVER_MODE | 运行模式 | release |
| DATABASE_DRIVER | 数据库驱动 | sqlite |
| DATABASE_DSN | 数据库连接串 | (SQLite 自动推导) |
| JWT_SECRET | JWT 密钥 | (必须配置) |
| JWT_ACCESS_EXPIRY | Access Token 过期时间 | 30m |
| JWT_REFRESH_EXPIRY | Refresh Token 过期时间 | 168h |

## 注意事项

- 生产环境必须配置 `JWT_SECRET`
- SQLite 数据库文件默认在运行目录下，建议配置绝对路径
- 如需切换 PostgreSQL，设置 `DATABASE_DRIVER=postgresql` 和 `DATABASE_DSN`