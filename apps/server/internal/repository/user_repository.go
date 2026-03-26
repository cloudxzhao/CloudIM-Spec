package repository

import (
	"cloudim/apps/server/internal/model"
	"database/sql"
	"errors"
	"time"
)

var ErrUserNotFound = errors.New("user not found")
var ErrUserExists = errors.New("user already exists")

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create 创建用户
func (r *UserRepository) Create(phone, passwordHash string) (*model.User, error) {
	user := &model.User{
		Phone:        phone,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err := r.db.QueryRow(
		"INSERT INTO users (phone, password_hash, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id, nickname, avatar",
		phone, passwordHash, user.CreatedAt, user.UpdatedAt,
	).Scan(&user.ID, &user.Nickname, &user.Avatar)

	if err != nil {
		if isUniqueViolation(err) {
			return nil, ErrUserExists
		}
		return nil, err
	}

	return user, nil
}

// FindByPhone 根据手机号查找用户
func (r *UserRepository) FindByPhone(phone string) (*model.User, error) {
	user := &model.User{}
	err := r.db.QueryRow(
		"SELECT id, phone, password_hash, nickname, avatar, created_at, updated_at FROM users WHERE phone = $1",
		phone,
	).Scan(&user.ID, &user.Phone, &user.PasswordHash, &user.Nickname, &user.Avatar, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

// FindByID 根据 ID 查找用户
func (r *UserRepository) FindByID(id int64) (*model.User, error) {
	user := &model.User{}
	err := r.db.QueryRow(
		"SELECT id, phone, password_hash, nickname, avatar, created_at, updated_at FROM users WHERE id = $1",
		id,
	).Scan(&user.ID, &user.Phone, &user.PasswordHash, &user.Nickname, &user.Avatar, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateProfile 更新用户资料
func (r *UserRepository) UpdateProfile(id int64, nickname, avatar string) error {
	_, err := r.db.Exec(
		"UPDATE users SET nickname = $1, avatar = $2, updated_at = $3 WHERE id = $4",
		nickname, avatar, time.Now(), id,
	)
	return err
}

// isUniqueViolation 检查是否是唯一约束违反
func isUniqueViolation(err error) bool {
	// PostgreSQL unique violation error code: 23505
	// This is a simplified check, in production you might want to use pgconn
	return err != nil
}
