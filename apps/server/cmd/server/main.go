package main

import (
	"cloudim/apps/server/internal/config"
	"cloudim/apps/server/internal/controller"
	"cloudim/apps/server/internal/database"
	"cloudim/apps/server/internal/repository"
	"cloudim/apps/server/internal/router"
	"cloudim/apps/server/internal/service"
	"cloudim/apps/server/ws"
	"fmt"
	"log"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化数据库
	if err := database.Init(&cfg.Database); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// 初始化仓库
	userRepo := repository.NewUserRepository(database.DB)
	messageRepo := repository.NewMessageRepository(database.DB)

	// 初始化服务
	authService := service.NewAuthService(userRepo, &cfg.JWT)
	controller.SetAuthService(authService)

	// 初始化 WebSocket 模块
	ws.Init(messageRepo)

	// 设置路由
	r := router.Setup()

	// 启动服务
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Starting server on %s", addr)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
