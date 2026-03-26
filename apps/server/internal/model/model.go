package model

import "time"

// User 用户模型
type User struct {
	ID           int64     `json:"id"`
	Phone        string    `json:"phone"`
	PasswordHash string    `json:"-"` // 不返回给客户端
	Nickname     string    `json:"nickname"`
	Avatar       string    `json:"avatar"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Message 消息模型
type Message struct {
	ID         int64     `json:"id"`
	SenderID   int64     `json:"sender_id"`
	ReceiverID int64     `json:"receiver_id"`
	Content    string    `json:"content"`
	Status     string    `json:"status"` // pending, delivered
	CreatedAt  time.Time `json:"created_at"`
}
