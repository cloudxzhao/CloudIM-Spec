package ws

import (
	"sync"
)

// Hub 维护所有活跃的 WebSocket 连接
type Hub struct {
	// 注册的客户端
	clients map[int64][]*Client

	// 注册通道
	register chan *Client

	// 注销通道
	unregister chan *Client

	// 消息广播通道
	broadcast chan *BroadcastMessage

	// 锁
	mu sync.RWMutex
}

// BroadcastMessage 广播消息
type BroadcastMessage struct {
	ReceiverID int64
	Message    []byte
}

// NewHub 创建新的 Hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[int64][]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *BroadcastMessage, 256),
	}
}

// Run 运行 Hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			h.broadcastMessage(message)
		}
	}
}

// registerClient 注册客户端
func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// 支持同一用户多连接
	h.clients[client.UserID] = append(h.clients[client.UserID], client)
}

// unregisterClient 注销客户端
func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[client.UserID]; ok {
		// 从用户连接列表中移除
		clients := h.clients[client.UserID]
		for i, c := range clients {
			if c == client {
				h.clients[client.UserID] = append(clients[:i], clients[i+1:]...)
				break
			}
		}
		// 如果用户没有连接了，删除该用户
		if len(h.clients[client.UserID]) == 0 {
			delete(h.clients, client.UserID)
		}
	}
}

// broadcastMessage 广播消息
func (h *Hub) broadcastMessage(message *BroadcastMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// 发送给接收者的所有连接
	if clients, ok := h.clients[message.ReceiverID]; ok {
		for _, client := range clients {
			select {
			case client.send <- message.Message:
			default:
				// 发送失败，客户端可能已断开
				close(client.send)
			}
		}
	}
}

// SendToUser 发送消息给指定用户
func (h *Hub) SendToUser(userID int64, message []byte) {
	h.broadcast <- &BroadcastMessage{
		ReceiverID: userID,
		Message:    message,
	}
}

// GetUserConnections 获取用户的连接数
func (h *Hub) GetUserConnections(userID int64) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, ok := h.clients[userID]; ok {
		return len(clients)
	}
	return 0
}
