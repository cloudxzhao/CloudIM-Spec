# CloudIM 部署文档

## 环境要求

### 后端
- Go 1.22+
- PostgreSQL 15+
- 操作系统：Linux / macOS / Windows

### 前端
- Node.js 18+
- npm 或 pnpm
- 操作系统：Windows 10/11（Electron 打包）

## 快速开始

### 1. 数据库初始化

```bash
# 通过 Docker 容器连接 PostgreSQL
sudo docker exec -it agent-postgres psql -U postgres

# 创建数据库
CREATE DATABASE cloudim;

# 退出
\q
```

### 2. 后端部署

```bash
cd apps/server

# 运行数据库迁移（需要安装 golang-migrate）
golang-migrate -path db/migrations -database "postgres://postgres:Admin_2026@localhost:5432/cloudim?sslmode=disable" up

# 启动后端服务
go run cmd/server/main.go

# 或者构建后运行
go build -o bin/server ./cmd/server
./bin/server
```

后端服务将在 `http://localhost:8080` 启动。

### 3. 前端部署（开发模式）

```bash
cd apps/client

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

前端将在 `http://localhost:5173` 启动。

### 4. Electron 客户端（开发模式）

```bash
cd apps/client

# 安装依赖
npm install

# 启动 Electron 开发模式
npm run electron:dev
```

### 5. Electron 客户端（生产打包）

```bash
cd apps/client

# 构建并打包
npm run electron:build
```

打包后的文件位于 `dist-electron/` 目录。

## API 接口

### 认证接口

| 接口 | 方法 | 描述 |
|------|------|------|
| `/api/v1/auth/register` | POST | 用户注册 |
| `/api/v1/auth/login` | POST | 用户登录 |
| `/api/v1/auth/captcha` | POST | 发送验证码 |
| `/api/v1/user/info` | GET | 获取用户信息 |
| `/api/v1/user/profile` | PUT | 更新用户资料 |

### WebSocket 连接

```
ws://localhost:8080/ws?token=<jwt_token>
```

## 配置说明

### 后端环境变量

| 变量名 | 默认值 | 描述 |
|--------|--------|------|
| `SERVER_PORT` | `8080` | 服务端端口 |
| `DB_HOST` | `localhost` | 数据库主机 |
| `DB_PORT` | `5432` | 数据库端口 |
| `DB_USER` | `postgres` | 数据库用户 |
| `DB_PASSWORD` | `Admin_2026` | 数据库密码 |
| `DB_NAME` | `cloudim` | 数据库名称 |
| `JWT_SECRET` | `cloudim-secret-key-change-in-production` | JWT 密钥 |
| `JWT_EXPIRATION` | `24` | Token 有效期（小时） |

## 常见问题

### 1. 后端启动失败

检查 PostgreSQL 是否正常运行：
```bash
sudo docker ps | grep agent-postgres
```

### 2. 前端依赖安装失败

尝试使用国内镜像：
```bash
npm config set registry https://registry.npmmirror.com
npm install
```

### 3. Electron 打包失败

确保已安装 Python 和构建工具：
```bash
# Windows
npm install --global windows-build-tools

# Linux
sudo apt install build-essential
```

## 开发测试

### 使用固定验证码

MVP 阶段，短信验证码固定为 `123456`，无需对接真实短信服务。

### 测试账号

可以通过注册流程创建测试账号，使用手机号 + 固定验证码 `123456` + 自定义密码。
