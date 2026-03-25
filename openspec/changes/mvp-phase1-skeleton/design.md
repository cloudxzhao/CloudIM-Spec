## Context

CloudIM MVP 第一期从零开始构建，需要搭建完整的技术骨架。当前状态：
- 无现有代码库
- PostgreSQL 已本地部署
- 技术栈已确定：Electron + Vue 3 前端，Go + Gin 后端

约束条件：
- Windows 10/11 兼容
- 安装包 ≤ 40MB
- 内存占用 ≤ 150MB
- 消息延迟 < 200ms

## Goals / Non-Goals

**Goals:**
- 搭建可编译运行的 monorepo 项目结构
- 实现后端 HTTP + WebSocket 双协议服务
- 实现前端 Electron 主进程与 Vue 渲染进程通信
- 实现用户注册/登录的完整认证流程
- 实现 WebSocket 连接、心跳、基础消息收发

**Non-Goals:**
- 图片/文件消息传输
- 音视频通话
- 消息云端漫游
- 敏感词过滤（后续迭代）
- 多端同步

## Decisions

### 1. Monorepo 目录结构

**决定：** 采用 `apps/` 分离前后端

```
CloudIM-spec/
├── apps/
│   ├── client/     # Electron 客户端
│   └── server/     # Go 后端
├── docs/           # 文档
└── openspec/       # 变更管理
```

**理由：** 前后端独立但同仓库管理，便于版本同步和共享类型定义。

**备选方案：**
- 分开仓库 → 增加同步成本，否决
- 单一目录混合 → 结构混乱，否决

### 2. 后端分层架构

**决定：** Controller → Service → Repository 三层架构

```
internal/
├── controller/   # HTTP 处理器
├── service/      # 业务逻辑
├── repository/   # 数据访问
├── model/        # 数据模型
└── middleware/   # 中间件
```

**理由：** 清晰的职责分离，便于测试和维护。

### 3. WebSocket 连接管理

**决定：** Hub-Client 模式 + gorilla/websocket

```
ws/
├── hub.go      # 连接注册、消息路由
├── client.go   # 单个连接的读写
└── message.go  # 消息解析
```

**理由：** Hub 负责广播和路由，Client 负责单个连接的生命周期，职责清晰。

**备选方案：**
- melody 库 → 封装过度，灵活性不足
- 原生 gorilla → 需要自行实现 Hub，当前决定即此方案

### 4. 前端状态管理

**决定：** Pinia + 独立 stores

```
stores/
├── user.ts     # 用户状态、Token
├── chat.ts     # 聊天消息、会话
└── contact.ts  # 好友、群组
```

**理由：** Pinia 是 Vue 3 官方推荐，类型支持好，模块化清晰。

### 5. 本地数据存储

**决定：** 客户端使用 SQLite (better-sqlite3)

**理由：**
- 离线消息缓存
- 会话列表本地化
- 减少服务端压力

**备选方案：**
- IndexedDB → API 复杂，Electron 环境下 SQLite 性能更好
- 文件 JSON → 查询效率低，不支持事务

### 6. JWT 认证

**决定：** JWT + Bearer Token，有效期 24 小时

**理由：** 无状态认证，适合即时通讯场景，WebSocket 连接通过 URL 参数传递 Token。

## Risks / Trade-offs

| 风险 | 缓解措施 |
|------|----------|
| Electron 打包体积超标 | 配置 Vite tree-shaking，禁用 Node.js 未使用模块 |
| WebSocket 连接不稳定 | 实现指数退避重连（5s, 10s, 30s），最多 5 次 |
| bcrypt 性能开销 | cost 设为 12，平衡安全与性能 |
| SQLite 跨平台兼容 | better-sqlite3 需编译，确保 electron-rebuild 正确执行 |
| JWT 无法主动失效 | MVP 阶段接受此限制，后续可引入 Redis 黑名单 |

## Migration Plan

### 部署步骤

1. 初始化 PostgreSQL 数据库 `cloudim`
2. 执行数据库迁移创建表结构
3. 启动后端服务 `go run cmd/server/main.go`
4. 启动前端开发服务器 `npm run dev`
5. 构建生产版本 `npm run build`

### 回滚策略

- 数据库迁移提供 `down.sql` 回滚脚本
- 配置文件版本化管理

## Open Questions

1. **验证码服务** - MVP 阶段使用固定验证码 `123456`，后续需要对接真实短信服务商
2. **头像存储** - MVP 阶段使用默认头像，后续需要确定 OSS/CDN 方案
3. **日志收集** - MVP 阶段本地文件日志，后续需要确定 ELK/云端日志方案