package controller

import (
	"cloudim/apps/server/internal/model"
	"cloudim/apps/server/internal/service"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

// RegisterRequest 注册请求
type RegisterRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Captcha  string `json:"captcha" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// CaptchaRequest 验证码请求
type CaptchaRequest struct {
	Phone string `json:"phone" binding:"required"`
}

// UserInfoResponse 用户信息响应
type UserInfoResponse struct {
	ID       int64  `json:"id"`
	Phone    string `json:"phone"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

// authService 认证服务（全局变量，由 main.go 初始化）
var authService *service.AuthService

// SetAuthService 设置认证服务实例
func SetAuthService(s *service.AuthService) {
	authService = s
}

// Register 用户注册
// @Summary 用户注册
// @Description 通过手机号和验证码注册用户
// @Tags 认证
// @Accept json
// @Produce json
// @Param body body RegisterRequest true "注册信息"
// @Success 200 {object} model.Response
// @Router /api/v1/auth/register [post]
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Error(model.CodeInvalidParam, "invalid request"))
		return
	}

	user, token, err := authService.Register(req.Phone, req.Captcha, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrUserExists) {
			c.JSON(http.StatusConflict, model.Error(model.CodeInvalidParam, "该手机号已注册"))
			return
		}
		if errors.Is(err, service.ErrInvalidCaptcha) {
			c.JSON(http.StatusBadRequest, model.Error(model.CodeInvalidParam, "验证码无效或已过期"))
			return
		}
		if errors.Is(err, service.ErrWeakPassword) {
			c.JSON(http.StatusBadRequest, model.Error(model.CodeInvalidParam, "密码必须至少 8 位，包含字母和数字"))
			return
		}
		c.JSON(http.StatusInternalServerError, model.Error(model.CodeInternalError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.Success(gin.H{
		"user": UserInfoResponse{
			ID:       user.ID,
			Phone:    user.Phone,
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
		},
		"token": token,
	}))
}

// Login 用户登录
// @Summary 用户登录
// @Description 通过手机号和密码登录
// @Tags 认证
// @Accept json
// @Produce json
// @Param body body LoginRequest true "登录信息"
// @Success 200 {object} model.Response
// @Router /api/v1/auth/login [post]
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Error(model.CodeInvalidParam, "invalid request"))
		return
	}

	user, token, err := authService.Login(req.Phone, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, model.Error(model.CodeNotFound, "用户不存在"))
			return
		}
		if errors.Is(err, service.ErrInvalidPassword) {
			c.JSON(http.StatusUnauthorized, model.Error(model.CodeUnauthorized, "密码错误"))
			return
		}
		c.JSON(http.StatusInternalServerError, model.Error(model.CodeInternalError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.Success(gin.H{
		"user": UserInfoResponse{
			ID:       user.ID,
			Phone:    user.Phone,
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
		},
		"token": token,
	}))
}

// SendCaptcha 发送验证码
// @Summary 发送验证码
// @Description 向指定手机号发送验证码（MVP 阶段固定验证码 123456）
// @Tags 认证
// @Accept json
// @Produce json
// @Param body body CaptchaRequest true "手机号"
// @Success 200 {object} model.Response
// @Router /api/v1/auth/captcha [post]
func SendCaptcha(c *gin.Context) {
	var req CaptchaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Error(model.CodeInvalidParam, "invalid request"))
		return
	}

	// MVP 阶段：固定验证码 123456
	// 实际生产环境需要对接短信服务商
	c.JSON(http.StatusOK, model.Success(gin.H{
		"message": "验证码已发送（MVP 阶段固定验证码：123456）",
	}))
}

// GetUserInfo 获取用户信息
// @Summary 获取用户信息
// @Description 获取当前登录用户的信息
// @Tags 用户
// @Produce json
// @Success 200 {object} model.Response
// @Router /api/v1/user/info [get]
func GetUserInfo(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, model.Error(model.CodeUnauthorized, "未认证"))
		return
	}

	user, err := authService.GetUserByID(userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Error(model.CodeInternalError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.Success(UserInfoResponse{
		ID:       user.ID,
		Phone:    user.Phone,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
	}))
}

// UpdateProfile 更新用户资料
// @Summary 更新用户资料
// @Description 更新用户昵称和头像
// @Tags 用户
// @Accept json
// @Produce json
// @Param body body UpdateProfileRequest true "更新信息"
// @Success 200 {object} model.Response
// @Router /api/v1/user/profile [put]
func UpdateProfile(c *gin.Context) {
	// TODO: 实现更新用户资料
	c.JSON(http.StatusOK, model.Success(nil))
}
