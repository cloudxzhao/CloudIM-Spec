package middleware

import (
	"cloudim/apps/server/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

// JWTAuth JWT 认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Authorization 头获取 Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, model.Error(model.CodeUnauthorized, "未提供认证信息"))
			c.Abort()
			return
		}

		// 检查 Bearer 格式
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, model.Error(model.CodeTokenInvalid, "无效的 Token 格式"))
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 验证 Token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 这里使用全局配置的 secret，实际应该从上下文获取
			return []byte("cloudim-secret-key-change-in-production"), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, model.Error(model.CodeTokenInvalid, "Token 无效"))
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// 从 Token 中提取用户 ID
			if userID, ok := claims["sub"].(float64); ok {
				c.Set("userID", int64(userID))
				c.Next()
				return
			}
		}

		c.JSON(http.StatusUnauthorized, model.Error(model.CodeTokenInvalid, "Token 无效"))
		c.Abort()
		return
	}
}
