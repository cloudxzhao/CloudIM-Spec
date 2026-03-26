# CloudIM

CloudIM 是一个 Windows PC 端即时通讯软件，采用 Electron + Vue 3 前端和 Go + Gin 后端架构。

## 项目结构

```
CloudIM/
├── apps/
│   ├── client/     # Electron 客户端
│   │   ├── electron/       # Electron 主进程和预加载脚本
│   │   ├── src/            # Vue 3 渲染进程源码
│   │   │   ├── api/        # API 封装
│   │   │   ├── components/ # 通用组件
│   │   │   ├── stores/     # Pinia 状态管理
│   │   │   ├── router/     # 路由配置
│   │   │   ├── views/      # 页面组件
│   │   │   └── main.js     # 应用入口
│   │   ├── package.json
│   │   └── vite.config.js
│   │
│   └── server/     # Go 后端服务
│       ├── cmd/server/     # 应用入口
│       ├── internal/       # 内部包
│       │   ├── config/     # 配置管理
│       │   ├── controller/ # HTTP 控制器
│       │   ├── database/   # 数据库连接
│       │   ├── middleware/ # 中间件
│       │   ├── model/      # 数据模型
│       │   ├── repository/ # 数据访问层
│       │   ├── router/     # 路由配置
│       │   └── service/    # 业务逻辑层
│       ├── ws/             # WebSocket 模块
│       ├── db/             # 数据库迁移
│       └── go.mod
│
├── docs/           # 文档（待创建）
├── openspec/       # 变更管理
├── DEPLOYMENT.md   # 部署文档
├── Makefile        # 构建脚本
└── package.json    # 根工作区配置
```

## 技术栈

### 后端
- **语言**: Go 1.22+
- **框架**: Gin
- **数据库**: PostgreSQL 15+
- **WebSocket**: gorilla/websocket
- **认证**: JWT (HS256, 24 小时有效期)
- **密码加密**: bcrypt (cost=12)

### 前端
- **框架**: Electron 33+
- **渲染**: Vue 3 + Vite
- **UI 库**: Element Plus
- **状态管理**: Pinia
- **本地存储**: SQLite (better-sqlite3)
- **HTTP 客户端**: Axios

## 快速开始

### 1. 环境准备

确保已安装：
- Go 1.22+
- Node.js 18+
- PostgreSQL 15+

### 2. 数据库初始化

```bash
# 创建数据库
createdb cloudim

# 或者通过 Docker
sudo docker exec -it agent-postgres psql -U postgres -c "CREATE DATABASE cloudim;"
```

### 3. 后端启动

```bash
cd apps/server

# 运行迁移
golang-migrate -path db/migrations -database "postgres://postgres:Admin_2026@localhost:5432/cloudim?sslmode=disable" up

# 启动服务
go run cmd/server/main.go
```

### 4. 前端启动

```bash
cd apps/client

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

## API 接口

| 接口 | 方法 | 描述 |
|------|------|------|
| `/api/v1/auth/register` | POST | 用户注册 |
| `/api/v1/auth/login` | POST | 用户登录 |
| `/api/v1/auth/captcha` | POST | 发送验证码 |
| `/api/v1/user/info` | GET | 获取用户信息 |
| `/ws` | WebSocket | WebSocket 连接 |

## MVP 功能范围

### 已实现
- [x] 用户注册/登录
- [x] JWT 认证
- [x] WebSocket 实时连接
- [x] 心跳机制
- [x] 消息收发
- [x] 离线消息存储
- [x] Electron 主窗口
- [x] 系统托盘
- [x] 登录/注册界面
- [x] 主界面框架

### 未实现（后续迭代）
- [ ] 图片/文件消息传输
- [ ] 音视频通话
- [ ] 消息云端漫游
- [ ] 群组聊天
- [ ] 好友管理

## 开发说明

### 验证码

MVP 阶段使用固定验证码 `123456`，无需对接真实短信服务。

### 测试账号

可以通过注册流程创建测试账号：
- 手机号：任意 11 位手机号
- 验证码：`123456`
- 密码：至少 8 位，包含字母和数字

## 许可证

MIT License
