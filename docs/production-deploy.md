# 生产模式部署指南

生产模式使用 Go `embed` 将前端静态文件嵌入后端二进制，构建单一可执行文件，同时服务 API 和前端。

| 仓库 | 地址 |
|------|------|
| Gitee | https://gitee.com/zhouws-chn/oa-nsdiy |
| GitHub | https://github.com/wilson-nsdiy/nsdiy-office-system |

## 部署方式

### 方式一：使用安装脚本（推荐）

生产环境推荐使用 root 模式安装，个人开发/测试可使用 `--user` 模式。

**Linux (root 模式 - 生产环境推荐):**

```bash
curl -sSL https://raw.githubusercontent.com/wilson-nsdiy/nsdiy-office-system/master/deploy/install.sh | sudo bash
```

**Linux (--user 模式 - 无需 sudo):**

```bash
curl -sSL https://raw.githubusercontent.com/wilson-nsdiy/nsdiy-office-system/master/deploy/install.sh | bash -s -- --user
```

安装脚本支持以下命令：

| 命令 | 说明 |
|------|------|
| `install` | 安装最新版本 |
| `upgrade` | 升级到最新版本 |
| `rollback -v <版本>` | 回退到指定版本 |
| `list-versions` | 列出可用版本 |
| `uninstall` | 卸载 |

**Windows:**

```cmd
REM 以管理员权限运行
install.bat

REM 安装指定版本
install.bat -v v1.0.0

REM 升级到最新版本
install.bat upgrade
```

### 方式二：从 Release 页面下载

1. 访问项目的 [Release 页面](https://github.com/wilson-nsdiy/nsdiy-office-system/releases)，下载对应平台的二进制文件
2. 上传至目标服务器并配置环境：

```bash
# Linux
chmod +x oa-nsdiy-linux-amd64
cp deploy/.env.example deploy/.env
# 编辑 .env 文件，配置必要的环境变量（至少设置 JWT_SECRET）
./oa-nsdiy-linux-amd64

# Windows
copy deploy\.env.example deploy\.env
REM 编辑 .env 文件
oa-nsdiy.exe
```

### 方式三：手动编译

```bash
# 使用构建脚本（推荐）
./deploy/build.sh           # 构建当前平台
./deploy/build.sh linux     # 交叉编译 Linux amd64

# 或手动构建
# 1. 构建前端
cd frontend && npm install && npm run build

# 2. 复制前端文件到后端
cp -r frontend/out/ backend/internal/web/dist

# 3. 构建后端（-tags embed 启用前端嵌入，-ldflags="-s -w" 去除调试信息减小体积）
cd backend
CGO_ENABLED=0 go build -tags embed -ldflags="-s -w" -o bin/oa-nsdiy ./cmd/server

# 4. 部署至目标服务器
scp bin/oa-nsdiy user@server:/path/to/deploy/
```

目标服务器上同样需要配置 `.env` 并启动：

```bash
cp deploy/.env.example deploy/.env
# 编辑 .env 文件，配置必要的环境变量
./oa-nsdiy
```

## 环境变量

从模板复制并填写：`cp deploy/.env.example deploy/.env`

| 变量 | 说明 | 默认值 |
|------|------|--------|
| SERVER_HOST | 监听地址 | 0.0.0.0 |
| SERVER_PORT | 监听端口 | 3001 |
| SERVER_MODE | 运行模式 | release |
| DATABASE_DRIVER | 数据库驱动 | sqlite |
| DATABASE_SOURCE | 数据库连接串 | (SQLite 自动推导) |
| JWT_SECRET | JWT 密钥 | **必须配置** |
| JWT_ACCESS_EXPIRY | Access Token 过期时间 | 30m |
| JWT_REFRESH_EXPIRY | Refresh Token 过期时间 | 168h |

## 安装模式对比

| | root 模式 (生产推荐) | --user 模式 |
|---|---|---|
| 安装目录 | `/opt/oa-nsdiy` | `~/.local/bin` |
| 配置目录 | `/etc/oa-nsdiy` | `~/.config/oa-nsdiy` |
| 系统用户 | 创建 `oa-nsdiy` | 使用当前用户 |
| systemd 服务 | 系统级 | 用户级 |
| 需要 sudo | 是 | 否 |

## 注意事项

- SQLite 数据库文件默认在运行目录下，建议配置绝对路径
- 如需切换 PostgreSQL，设置 `DATABASE_DRIVER=postgres` 和 `DATABASE_SOURCE`
- 生产环境建议使用安装脚本自动配置 systemd 服务
