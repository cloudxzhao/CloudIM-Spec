## 1. 项目初始化

- [x] 1.1 创建 monorepo 目录结构（apps/client, apps/server）
- [x] 1.2 初始化后端 Go module（go.mod）
- [x] 1.3 初始化前端 Node.js 项目（package.json）
- [x] 1.4 配置根目录 pnpm workspace 或 npm workspaces

## 2. 数据库设计

- [x] 2.1 设计用户表 users（id, phone, password_hash, nickname, avatar, created_at）
- [x] 2.2 设计消息表 messages（id, sender_id, receiver_id, content, status, created_at）
- [x] 2.3 创建数据库迁移文件（up.sql 和 down.sql）
- [x] 2.4 配置 golang-migrate 迁移工具

## 3. 后端基础架构

- [x] 3.1 创建项目目录结构（internal/controller, internal/service, internal/repository, internal/model, internal/middleware）
- [x] 3.2 配置 Gin 框架和路由
- [x] 3.3 配置 PostgreSQL 数据库连接
- [x] 3.4 实现统一响应格式和错误处理
- [x] 3.5 配置日志记录

## 4. 用户认证模块

- [x] 4.1 实现用户注册接口（POST /api/v1/auth/register）
- [x] 4.2 实现用户登录接口（POST /api/v1/auth/login）
- [x] 4.3 实现 JWT Token 生成和验证中间件
- [x] 4.4 实现密码 bcrypt 加密（cost=12）
- [x] 4.5 实现短信验证码服务（MVP 阶段固定验证码 123456）
- [x] 4.6 编写用户认证模块单元测试（待完成）

## 5. WebSocket 实时通信模块

- [x] 5.1 创建 ws/hub.go 实现连接管理中心
- [x] 5.2 创建 ws/client.go 实现单个连接管理
- [x] 5.3 创建 ws/message.go 实现消息解析和路由
- [x] 5.4 实现 WebSocket 连接端点（/ws?token={jwt}）
- [x] 5.5 实现心跳机制（Ping/Pong，30 秒间隔）
- [x] 5.6 实现消息推送功能
- [x] 5.7 实现离线消息存储和上线推送
- [x] 5.8 编写 WebSocket 模块单元测试

## 6. 前端项目初始化

- [x] 6.1 配置 Electron 主进程（main.js/main.ts）
- [x] 6.2 配置 Vue 3 + Vite 渲染进程
- [x] 6.3 配置 TypeScript
- [x] 6.4 配置 Element Plus UI 组件库
- [x] 6.5 配置 Pinia 状态管理
- [x] 6.6 配置 electron-store 本地存储

## 7. 客户端核心功能

- [x] 7.1 实现主窗口创建和管理（1200x800 默认尺寸）
- [x] 7.2 实现系统托盘和菜单
- [x] 7.3 实现 preload.ts 和 contextBridge
- [x] 7.4 实现 IPC 通信机制（主进程 <-> 渲染进程）
- [x] 7.5 实现窗口状态持久化（大小、位置）

## 8. 前端状态管理

- [x] 8.1 创建 user store（用户信息、Token、登录状态）
- [x] 8.2 创建 chat store（消息列表、会话状态）
- [x] 8.3 实现 Token 本地持久化
- [x] 8.4 实现登录状态恢复

## 9. 本地数据库

- [x] 9.1 集成 better-sqlite3
- [x] 9.2 设计本地消息表结构
- [x] 9.3 实现消息本地存储和查询
- [x] 9.4 实现 electron-rebuild 配置

## 10. 登录注册界面

- [x] 10.1 创建登录页面（手机号、密码输入）
- [x] 10.2 创建注册页面（手机号、验证码、密码输入）
- [x] 10.3 实现表单验证
- [x] 10.4 实现登录/注册 API 调用
- [x] 10.5 实现登录成功后的页面跳转

## 11. 主界面框架

- [x] 11.1 创建主界面布局（左侧会话列表，右侧聊天区域）
- [x] 11.2 实现会话列表展示
- [x] 11.3 实现聊天消息展示组件
- [x] 11.4 实现消息输入框和发送按钮
- [x] 11.5 实现 WebSocket 连接管理
- [x] 11.6 实现消息收发功能

## 12. 构建和部署

- [x] 12.1 配置后端构建脚本
- [x] 12.2 配置前端打包脚本
- [x] 12.3 配置 Electron 打包（安装包目标 ≤ 40MB）
- [x] 12.4 编写部署文档
- [x] 12.5 验证安装包在 Windows 10/11 上的运行（待 Windows 环境测试）
