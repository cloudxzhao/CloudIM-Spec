## ADDED Requirements

### Requirement: 客户端可以建立 WebSocket 连接

系统支持客户端通过 WebSocket 协议与服务端建立持久连接。

- **协议**：WebSocket (ws:// 或 wss://)
- **认证方式**：URL 参数传递 JWT Token
- **连接管理**：Hub 模式管理所有连接

#### Scenario: 成功建立 WebSocket 连接
- **WHEN** 客户端携带有效 Token 向 /ws 端点发起 WebSocket 连接请求
- **THEN** 服务端完成握手，建立连接，将连接注册到 Hub，关联用户 ID

#### Scenario: 重复连接处理
- **WHEN** 同一用户从多个客户端建立 WebSocket 连接
- **THEN** 服务端允许多连接，所有连接均保持活跃

### Requirement: 客户端可以发送心跳保持连接

系统通过心跳机制检测连接活性，防止空闲连接超时断开。

- **心跳间隔**：客户端每 30 秒发送一次 Ping 消息
- **超时检测**：服务端 60 秒未收到消息判定为超时
- **心跳消息格式**：{ "type": "ping", "timestamp": <unix_timestamp> }

#### Scenario: 心跳响应
- **WHEN** 客户端发送 Ping 消息
- **THEN** 服务端在 5 秒内回复 Pong 消息 { "type": "pong", "timestamp": <unix_timestamp> }

#### Scenario: 连接超时清理
- **WHEN** 服务端 60 秒内未收到客户端任何消息
- **THEN** 服务端主动关闭连接，从 Hub 中移除

### Requirement: 用户可以发送文本消息

系统支持用户通过 WebSocket 连接发送文本消息。

- **消息格式**：JSON { "type": "message", "data": { "to": <target_id>, "content": <text> } }
- **消息类型**：私聊消息（MVP 阶段仅支持私聊）
- **消息确认**：服务端返回发送结果确认

#### Scenario: 发送消息成功
- **WHEN** 用户通过 WebSocket 发送格式正确的消息
- **THEN** 服务端返回 ACK 确认 { "type": "ack", "msg_id": <message_id> }

#### Scenario: 接收方不在线
- **WHEN** 消息目标用户未建立 WebSocket 连接
- **THEN** 服务端存储离线消息，待目标用户上线后推送

#### Scenario: 消息格式错误
- **WHEN** 客户端发送的消息 JSON 格式无效或缺少必需字段
- **THEN** 服务端返回错误 { "type": "error", "code": "INVALID_MESSAGE" }

### Requirement: 用户可以接收消息

系统支持用户通过 WebSocket 连接接收其他用户发送的消息。

- **推送格式**：JSON { "type": "message", "data": { "from": <sender_id>, "content": <text>, "timestamp": <unix_timestamp> } }
- **消息顺序**：按时间顺序投递

#### Scenario: 在线接收消息
- **WHEN** 目标用户 WebSocket 连接在线
- **THEN** 服务端立即推送消息到目标客户端

#### Scenario: 离线消息上线推送
- **WHEN** 用户上线建立 WebSocket 连接
- **THEN** 服务端推送该用户所有未读离线消息（按时间顺序）

### Requirement: 系统存储离线消息

系统为离线用户存储消息，待用户上线后推送。

- **存储介质**：PostgreSQL 数据库
- **存储内容**：消息 ID、发送者 ID、接收者 ID、消息内容、时间戳、状态
- **消息状态**：pending（待推送）、delivered（已推送）

#### Scenario: 存储离线消息
- **WHEN** 消息目标用户不在线
- **THEN** 服务端将消息存入数据库，状态标记为 pending

#### Scenario: 上线拉取离线消息
- **WHEN** 用户建立 WebSocket 连接
- **THEN** 服务端查询该用户所有 pending 状态消息，按时间顺序推送，推送后标记为 delivered

### Requirement: 客户端可以断开 WebSocket 连接

系统支持客户端主动断开 WebSocket 连接，服务端进行清理。

- **断开方式**：客户端调用 close() 或网络异常
- **清理操作**：从 Hub 移除连接，释放资源

#### Scenario: 客户端主动断开
- **WHEN** 客户端调用 WebSocket.close()
- **THEN** 服务端收到 close 事件，从 Hub 移除该连接

#### Scenario: 网络异常断开
- **WHEN** 网络中断导致 WebSocket 连接异常
- **THEN** 服务端在 60 秒超时检测后清理连接
