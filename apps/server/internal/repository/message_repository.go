package repository

import (
	"cloudim/apps/server/internal/model"
	"database/sql"
	"time"
)

type MessageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

// Create 创建消息
func (r *MessageRepository) Create(senderID, receiverID int64, content string) (*model.Message, error) {
	msg := &model.Message{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
		Status:     "pending",
		CreatedAt:  time.Now(),
	}

	err := r.db.QueryRow(
		"INSERT INTO messages (sender_id, receiver_id, content, status, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		senderID, receiverID, content, msg.Status, msg.CreatedAt,
	).Scan(&msg.ID)

	if err != nil {
		return nil, err
	}

	return msg, nil
}

// FindPendingByReceiver 查找接收者的未读消息
func (r *MessageRepository) FindPendingByReceiver(receiverID int64) ([]*model.Message, error) {
	rows, err := r.db.Query(
		"SELECT id, sender_id, receiver_id, content, status, created_at FROM messages WHERE receiver_id = $1 AND status = 'pending' ORDER BY created_at ASC",
		receiverID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*model.Message
	for rows.Next() {
		msg := &model.Message{}
		err := rows.Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.Content, &msg.Status, &msg.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, rows.Err()
}

// MarkAsDelivered 将消息标记为已送达
func (r *MessageRepository) MarkAsDelivered(receiverID int64, messageIDs []int64) error {
	if len(messageIDs) == 0 {
		return nil
	}

	_, err := r.db.Exec(
		"UPDATE messages SET status = 'delivered' WHERE id = ANY($1) AND receiver_id = $2",
		messageIDs, receiverID,
	)
	return err
}

// FindByConversation 查找会话消息
func (r *MessageRepository) FindByConversation(userID, otherUserID int64, limit int, offset int) ([]*model.Message, error) {
	rows, err := r.db.Query(
		`SELECT id, sender_id, receiver_id, content, status, created_at
		 FROM messages
		 WHERE (sender_id = $1 AND receiver_id = $2) OR (sender_id = $2 AND receiver_id = $1)
		 ORDER BY created_at DESC
		 LIMIT $3 OFFSET $4`,
		userID, otherUserID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*model.Message
	for rows.Next() {
		msg := &model.Message{}
		err := rows.Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.Content, &msg.Status, &msg.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, rows.Err()
}
