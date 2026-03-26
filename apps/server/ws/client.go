package ws

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	// 写入超时时间
	writeWait = 10 * time.Second

	// 心跳超时时间（60 秒未收到消息判定为超时）
	pongWait = 60 * time.Second

	// 发送心跳间隔（30 秒）
	pingPeriod = (pongWait * 9) / 10

	// 最大消息大小
	maxMessageSize = 4096
)

// Client WebSocket 客户端连接
type Client struct {
	// 用户 ID
	UserID int64

	// WebSocket 连接
	Conn *websocket.Conn

	// 发送消息通道
	send chan []byte

	// Hub 实例
	hub *Hub

	// 锁
	mu sync.Mutex

	// 最后活动时间
	lastActive time.Time
}

// Upgrader WebSocket 升级器
var Upgrader = websocket.Upgrader{
	// 允许所有来源（生产环境应该限制）
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// NewClient 创建新的客户端
func NewClient(userID int64, conn *websocket.Conn, hub *Hub) *Client {
	return &Client{
		UserID:     userID,
		Conn:       conn,
		send:       make(chan []byte, 256),
		hub:        hub,
		lastActive: time.Now(),
	}
}

// Start 启动客户端（读写协程）
func (c *Client) Start() {
	// 注册到 Hub
	c.hub.register <- c

	// 启动写协程
	go c.writePump()

	// 启动读协程
	go c.readPump()
}

// readPump 读取消息
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.lastActive = time.Now()
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// 处理收到的消息
		c.handleMessage(message)
	}
}

// writePump 写入消息
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub 关闭了通道
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// Send 发送消息
func (c *Client) Send(message []byte) {
	select {
	case c.send <- message:
	default:
		// 通道已满
		log.Printf("Client send channel full for user %d", c.UserID)
	}
}

// Close 关闭连接
func (c *Client) Close() {
	c.Conn.Close()
}
