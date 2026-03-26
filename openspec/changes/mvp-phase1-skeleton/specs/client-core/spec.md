## ADDED Requirements

### Requirement: 应用可以创建主窗口

系统支持创建 Electron 主窗口用于显示应用界面。

- **窗口尺寸**：默认 1200x800，最小 800x600
- **窗口行为**：可调整大小，可最小化/最大化/关闭
- **安全性**：禁用 Node.js 集成，启用上下文隔离

#### Scenario: 创建主窗口
- **WHEN** 应用启动时
- **THEN** 创建 BrowserWindow，加载 Vue 应用入口 URL

#### Scenario: 窗口状态持久化
- **WHEN** 用户调整窗口大小或位置后关闭
- **THEN** 保存窗口状态，下次启动时恢复

#### Scenario: 窗口关闭行为
- **WHEN** 用户点击关闭按钮
- **THEN** MVP 阶段直接退出应用（后续可改为最小化到托盘）

### Requirement: 应用可以创建系统托盘

系统支持在系统托盘区显示应用图标和菜单。

- **托盘图标**：应用 Logo
- **菜单项**：显示/隐藏主窗口、退出应用
- **点击行为**：单击显示/切换主窗口

#### Scenario: 创建系统托盘
- **WHEN** 应用启动时
- **THEN** 在系统托盘区显示应用图标

#### Scenario: 托盘菜单 - 显示主窗口
- **WHEN** 用户点击托盘菜单"显示"项
- **THEN** 主窗口显示并聚焦

#### Scenario: 托盘菜单 - 退出应用
- **WHEN** 用户点击托盘菜单"退出"项
- **THEN** 应用正常退出，清理资源

### Requirement: 主进程与渲染进程可以通信

系统通过 IPC 实现主进程与渲染进程之间的通信。

- **IPC 方向**：渲染进程 → 主进程（invoke/handle），主进程 → 渲染进程（webContents.send）
- **预加载脚本**：preload.ts 暴露安全的 API 给渲染进程
- **安全性**：使用 contextBridge，不暴露 Node.js API

#### Scenario: 渲染进程调用主进程 API
- **WHEN** 渲染进程需要访问系统功能（如关闭窗口、访问文件系统）
- **THEN** 通过 IPC invoke 调用主进程暴露的 handle 方法

#### Scenario: 主进程推送事件到渲染进程
- **WHEN** 主进程收到系统事件（如网络状态变化）
- **THEN** 通过 webContents.send 推送到渲染进程

### Requirement: 应用可以管理用户会话状态

系统使用 Pinia 管理用户登录状态和用户信息。

- **状态内容**：用户 ID、手机号、昵称、头像 URL、JWT Token、登录状态
- **持久化**：Token 存储到本地（electron-store）
- **状态同步**：登录/登出时更新状态

#### Scenario: 用户登录成功
- **WHEN** 用户完成登录认证
- **THEN** Pinia store 更新用户信息和 Token，状态标记为 logged-in

#### Scenario: 用户登出
- **WHEN** 用户点击登出
- **THEN** Pinia store 清空用户信息和 Token，状态标记为 logged-out，跳转到登录页

#### Scenario: 应用启动恢复登录状态
- **WHEN** 应用启动时本地存在有效 Token
- **THEN** 从本地读取 Token，恢复用户登录状态

### Requirement: 客户端可以存储本地数据

系统使用 SQLite 存储本地数据用于离线访问。

- **数据库**：SQLite (better-sqlite3)
- **存储内容**：会话列表、消息历史、用户设置
- **位置**：应用数据目录

#### Scenario: 存储消息历史
- **WHEN** 收到新消息或发送消息
- **THEN** 消息同时存储到 SQLite 数据库

#### Scenario: 读取本地消息历史
- **WHEN** 用户打开聊天窗口
- **THEN** 从 SQLite 加载该聊天的历史消息

#### Scenario: 存储用户设置
- **WHEN** 用户修改设置（如通知开关、主题）
- **THEN** 设置存储到 SQLite，下次启动时恢复
