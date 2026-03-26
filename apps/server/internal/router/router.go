package router

import (
	"cloudim/apps/server/internal/controller"
	"cloudim/apps/server/internal/middleware"
	"cloudim/apps/server/ws"
	"github.com/gin-gonic/gin"
)

// Setup 设置路由
func Setup() *gin.Engine {
	r := gin.New()

	// 全局中间件
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1 路由
	v1 := r.Group("/api/v1")
	{
		// 认证相关
		auth := v1.Group("/auth")
		{
			auth.POST("/register", controller.Register)
			auth.POST("/login", controller.Login)
			auth.POST("/captcha", controller.SendCaptcha)
		}

		// 需要认证的路由
		authRequired := v1.Group("")
		authRequired.Use(middleware.JWTAuth())
		{
			authRequired.GET("/user/info", controller.GetUserInfo)
			authRequired.PUT("/user/profile", controller.UpdateProfile)
		}
	}

	// WebSocket 端点
	r.GET("/ws", ws.HandleWebSocket)

	return r
}
