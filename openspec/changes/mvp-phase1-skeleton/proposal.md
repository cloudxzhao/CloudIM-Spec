## Why

CloudIM 需要搭建一个可运行的 MVP 骨架，验证技术架构的可行性，并为后续功能迭代奠定基础。这是从零开始构建 Windows PC 端即时通讯软件的第一步。

## What Changes

- 创建 monorepo 项目结构（apps/client + apps/server）
- 搭建 Go + Gin 后端服务框架，包含 HTTP API 和 WebSocket 服务
- 搭建 Electron + Vue 3 前端客户端框架，包含窗口管理和系统托盘
- 实现 PostgreSQL 数据库表结构和迁移管理
- 实现用户注册/登录和 JWT 认证
- 实现 WebSocket 实时通信基础能力（连接、心跳、消息收发）

## Capabilities

### New Capabilities

- `user-auth`: 用户注册、登录、JWT 认证机制
- `realtime-messaging`: WebSocket 实时消息通信
- `client-core`: Electron 客户端核心框架（窗口、托盘、IPC）

### Modified Capabilities

<!-- 无现有能力需要修改 -->

## Impact

**后端 (apps/server)**
- Go 1.22+ / Gin 框架
- gorilla/websocket 用于实时通信
- PostgreSQL 作为主数据库
- golang-migrate 用于数据库迁移

**前端 (apps/client)**
- Electron 33+ 桌面框架
- Vue 3 + Vite 渲染进程
- Pinia 状态管理
- Element Plus UI 组件库
- SQLite (better-sqlite3) 本地存储

**开发环境**
- PostgreSQL 15+ 需要本地部署
- Node.js 18+ 用于前端构建
- Go 1.22+ 用于后端编译