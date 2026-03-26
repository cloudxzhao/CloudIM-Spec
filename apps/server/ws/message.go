package ws

import (
	"cloudim/apps/server/internal/model"
	"cloudim/apps/server/internal/repository"
	"encoding/json"
	"log"
	"time"
)

// 消息类型
const (
	MessageTypePing      = "ping"
	MessageTypePong      = "pong"
	MessageTypeMessage   = "message"
	MessageTypeAck       = "ack"
	MessageTypeError     = "error"
)

// WSMessage WebSocket 消息基础结构
type WSMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data,omitempty"`
}

// PingMessage 心跳请求
type PingMessage struct {
	Timestamp int64 `json:"timestamp"`
}

// PongMessage 心跳响应
type PongMessage struct {
	Timestamp int64 `json:"timestamp"`
}

// ChatMessage 聊天消息
type ChatMessage struct {
	To        int64  `json:"to"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

// AckMessage 确认消息
type AckMessage struct {
	MsgID int64 `json:"msg_id"`
}

// ErrorMessage 错误消息
type ErrorMessage struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// MessageRepository 消息仓库接口
var messageRepo *repository.MessageRepository

// SetMessageRepository 设置消息仓库
func SetMessageRepository(repo *repository.MessageRepository) {
	messageRepo = repo
}

// handleMessage 处理收到的消息
func (c *Client) handleMessage(data []byte) {
	var msg WSMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		c.sendError("INVALID_JSON", "invalid message format")
		return
	}

	switch msg.Type {
	case MessageTypePing:
		c.handlePing(msg.Data)
	case MessageTypeMessage:
		c.handleChatMessage(msg.Data)
	default:
		c.sendError("UNKNOWN_TYPE", "unknown message type")
	}
}

// handlePing 处理心跳
func (c *Client) handlePing(data json.RawMessage) {
	var ping PingMessage
	if err := json.Unmarshal(data, &ping); err != nil {
		c.sendError("INVALID_PING", "invalid ping format")
		return
	}

	// 回复 Pong
	pong := WSMessage{
		Type: MessageTypePong,
		Data: mustMarshal(PongMessage{Timestamp: time.Now().Unix()}),
	}

	c.sendJSON(pong)
}

// handleChatMessage 处理聊天消息
func (c *Client) handleChatMessage(data json.RawMessage) {
	var chatMsg ChatMessage
	if err := json.Unmarshal(data, &chatMsg); err != nil {
		c.sendError("INVALID_MESSAGE", "invalid message format")
		return
	}

	// 验证消息内容
	if chatMsg.To <= 0 || chatMsg.Content == "" {
		c.sendError("INVALID_MESSAGE", "missing required fields")
		return
	}

	// 创建消息记录
	if messageRepo != nil {
		msg, err := messageRepo.Create(c.UserID, chatMsg.To, chatMsg.Content)
		if err != nil {
			log.Printf("Failed to save message: %v", err)
			// 继续发送，不阻塞
		} else {
			// 返回 ACK
			ack := WSMessage{
				Type: MessageTypeAck,
				Data: mustMarshal(AckMessage{MsgID: msg.ID}),
			}
			c.sendJSON(ack)
		}
	}

	// 构建推送消息
	pushMsg := WSMessage{
		Type: MessageTypeMessage,
		Data: mustMarshal(ChatMessage{
			To:        chatMsg.To,
			Content:   chatMsg.Content,
			Timestamp: time.Now().Unix(),
		}),
	}

	// 发送给接收者
	c.hub.SendToUser(chatMsg.To, mustMarshal(pushMsg))
}

// sendError 发送错误消息
func (c *Client) sendError(code, message string) {
	errMsg := WSMessage{
		Type: MessageTypeError,
		Data: mustMarshal(ErrorMessage{Code: code, Message: message}),
	}
	c.sendJSON(errMsg)
}

// sendJSON 发送 JSON 消息
func (c *Client) sendJSON(v interface{}) {
	data := mustMarshal(v)
	c.Send(data)
}

// mustMarshal JSON 序列化
func mustMarshal(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		log.Printf("Failed to marshal: %v", err)
		return []byte("{}")
	}
	return data
}

// BroadcastOfflineMessage 广播离线消息（当用户上线时）
func BroadcastOfflineMessage(hub *Hub, userID int64, messages []*model.Message) {
	for _, msg := range messages {
		pushMsg := WSMessage{
			Type: MessageTypeMessage,
			Data: mustMarshal(ChatMessage{
				To:        msg.ReceiverID,
				Content:   msg.Content,
				Timestamp: msg.CreatedAt.Unix(),
			}),
		}
		hub.SendToUser(userID, mustMarshal(pushMsg))
	}
}
