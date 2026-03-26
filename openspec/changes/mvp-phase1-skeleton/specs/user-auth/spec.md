## ADDED Requirements

### Requirement: 用户可以通过手机号注册

系统允许新用户通过手机号和验证码进行注册。

- **前置条件**：用户未注册过该手机号
- **输入**：手机号、6 位数字验证码、密码
- **处理**：验证验证码有效性，检查手机号是否已存在，加密存储密码
- **输出**：注册成功返回用户 ID 和 JWT Token

#### Scenario: 注册成功
- **WHEN** 用户提供未注册的手机号、正确的验证码和符合要求的密码
- **THEN** 系统创建用户账户，返回用户 ID 和 JWT Token

#### Scenario: 手机号已存在
- **WHEN** 用户使用已注册的手机号发起注册请求
- **THEN** 系统返回错误"该手机号已注册"，HTTP 409 Conflict

#### Scenario: 验证码错误
- **WHEN** 用户输入错误的验证码
- **THEN** 系统返回错误"验证码无效或已过期"，HTTP 400 Bad Request

#### Scenario: 密码格式不符合要求
- **WHEN** 用户密码长度小于 8 位或不包含字母和数字
- **THEN** 系统返回错误"密码必须至少 8 位，包含字母和数字"，HTTP 400 Bad Request

### Requirement: 用户可以通过手机号和密码登录

系统允许已注册用户通过手机号和密码进行登录。

- **前置条件**：用户已完成注册
- **输入**：手机号、密码
- **处理**：验证手机号存在性，比对密码哈希
- **输出**：登录成功返回 JWT Token 和用户基本信息

#### Scenario: 登录成功
- **WHEN** 用户提供正确的手机号和密码
- **THEN** 系统返回 JWT Token（有效期 24 小时）和用户基本信息（ID、手机号、昵称、头像 URL）

#### Scenario: 用户不存在
- **WHEN** 用户使用未注册的手机号登录
- **THEN** 系统返回错误"用户不存在"，HTTP 404 Not Found

#### Scenario: 密码错误
- **WHEN** 用户提供错误密码
- **THEN** 系统返回错误"密码错误"，HTTP 401 Unauthorized

### Requirement: 系统使用 JWT 进行身份认证

系统对需要认证的 API 请求进行 JWT Token 验证。

- **Token 格式**：Bearer Token
- **有效期**：24 小时
- **签名算法**：HS256
- **Claims**：sub (用户 ID), exp (过期时间), iat (签发时间)

#### Scenario: 携带有效 Token 访问受保护接口
- **WHEN** 请求头包含有效的 Bearer Token 且未过期
- **THEN** 系统允许访问，从 Token 中提取用户 ID 注入上下文

#### Scenario: Token 缺失
- **WHEN** 请求受保护接口但未携带 Authorization 头
- **THEN** 系统返回错误"未提供认证信息"，HTTP 401 Unauthorized

#### Scenario: Token 过期
- **WHEN** 请求携带已过期的 Token
- **THEN** 系统返回错误"Token 已过期"，HTTP 401 Unauthorized

#### Scenario: Token 格式无效
- **WHEN** 请求携带格式错误的 Token（非 Bearer 格式或 JWT 解析失败）
- **THEN** 系统返回错误"无效的 Token 格式"，HTTP 401 Unauthorized

### Requirement: WebSocket 连接通过 Token 认证

WebSocket 连接建立时需要通过 URL 参数传递 JWT Token 进行认证。

- **连接地址**：ws://{host}/ws?token={jwt_token}
- **认证时机**：WebSocket 握手阶段
- **失败处理**：认证失败立即关闭连接

#### Scenario: WebSocket 连接认证成功
- **WHEN** 客户端携带有效 JWT Token 建立 WebSocket 连接
- **THEN** 服务端完成握手，建立连接并关联用户 ID

#### Scenario: WebSocket 连接 Token 无效
- **WHEN** 客户端携带无效或过期的 Token 建立 WebSocket 连接
- **THEN** 服务端拒绝握手，返回 HTTP 401 并关闭连接

#### Scenario: WebSocket 连接未提供 Token
- **WHEN** 客户端未在 URL 参数中提供 Token
- **THEN** 服务端拒绝握手，返回错误"缺少认证 Token"，HTTP 401 Unauthorized

### Requirement: 密码加密存储

用户密码必须使用 bcrypt 算法加密存储。

- **加密算法**：bcrypt
- **Cost 因子**：12
- **密码验证**：使用 bcrypt.CompareHashAndPassword

#### Scenario: 新用户密码加密
- **WHEN** 用户注册时提交明文密码
- **THEN** 系统使用 bcrypt 加密后存储到数据库，不保存明文

#### Scenario: 登录密码验证
- **WHEN** 用户登录时提交密码
- **THEN** 系统使用 bcrypt 比对密码哈希，验证通过才返回 Token
