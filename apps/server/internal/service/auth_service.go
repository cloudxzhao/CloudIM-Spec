package service

import (
	"cloudim/apps/server/internal/config"
	"cloudim/apps/server/internal/model"
	"cloudim/apps/server/internal/repository"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	ErrInvalidPassword  = errors.New("invalid password")
	ErrUserNotFound     = errors.New("user not found")
	ErrUserExists       = errors.New("user already exists")
	ErrInvalidCaptcha   = errors.New("invalid captcha")
	ErrWeakPassword     = errors.New("password must be at least 8 characters and contain letters and numbers")
)

type AuthService struct {
	userRepo *repository.UserRepository
	config   *config.JWTConfig
}

func NewAuthService(userRepo *repository.UserRepository, cfg *config.JWTConfig) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		config:   cfg,
	}
}

// Register 用户注册
func (s *AuthService) Register(phone, captcha, password string) (*model.User, string, error) {
	// 验证验证码（MVP 阶段固定验证码 123456）
	if captcha != "123456" {
		return nil, "", ErrInvalidCaptcha
	}

	// 验证密码强度
	if !validatePassword(password) {
		return nil, "", ErrWeakPassword
	}

	// 密码加密
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	// 创建用户
	user, err := s.userRepo.Create(phone, string(hash))
	if err != nil {
		if errors.Is(err, repository.ErrUserExists) {
			return nil, "", ErrUserExists
		}
		return nil, "", err
	}

	// 生成 Token
	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// Login 用户登录
func (s *AuthService) Login(phone, password string) (*model.User, string, error) {
	// 查找用户
	user, err := s.userRepo.FindByPhone(phone)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, "", ErrUserNotFound
		}
		return nil, "", err
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", ErrInvalidPassword
	}

	// 生成 Token
	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// GetUserByID 根据 ID 获取用户
func (s *AuthService) GetUserByID(id int64) (*model.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

// validatePassword 验证密码强度
func validatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	hasLetter := false
	hasNumber := false
	for _, c := range password {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
			hasLetter = true
		}
		if c >= '0' && c <= '9' {
			hasNumber = true
		}
	}
	return hasLetter && hasNumber
}

// generateToken 生成 JWT Token
func (s *AuthService) generateToken(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Duration(s.config.Expiration) * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.Secret))
}

// VerifyToken 验证 JWT Token
func (s *AuthService) VerifyToken(tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.Secret), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := int64(claims["sub"].(float64))
		return userID, nil
	}

	return 0, errors.New("invalid token")
}
