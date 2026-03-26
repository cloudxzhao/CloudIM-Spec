package ws

import (
	"cloudim/apps/server/internal/model"
	"cloudim/apps/server/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strings"
)

// GlobalHub 全局 Hub 实例
var GlobalHub *Hub

// Init 初始化 WebSocket 模块
func Init(repo *repository.MessageRepository) {
	GlobalHub = NewHub()
	SetMessageRepository(repo)

	// 启动 Hub
	go GlobalHub.Run()
}

// HandleWebSocket 处理 WebSocket 连接
func HandleWebSocket(c *gin.Context) {
	// 从 URL 参数获取 Token
	tokenString := c.Query("token")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, model.Error(model.CodeUnauthorized, "缺少认证 Token"))
		return
	}

	// 验证 Token
	userID, err := verifyToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.Error(model.CodeTokenInvalid, "无效的 Token"))
		return
	}

	// 升级 WebSocket 连接
	conn, err := Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	// 创建客户端
	client := NewClient(userID, conn, GlobalHub)

	// 启动客户端
	client.Start()

	// 注意：不要在这里关闭连接，由 client 协程管理

	// 检查并推送离线消息
	go deliverOfflineMessages(userID)
}

// verifyToken 验证 JWT Token
func verifyToken(tokenString string) (int64, error) {
	// 从 Authorization 头获取 Token（如果是 Bearer 格式）
	if strings.HasPrefix(tokenString, "Bearer ") {
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("cloudim-secret-key-change-in-production"), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userID, ok := claims["sub"].(float64); ok {
			return int64(userID), nil
		}
	}

	return 0, nil
}

// deliverOfflineMessages 推送离线消息
func deliverOfflineMessages(userID int64) {
	if messageRepo == nil {
		return
	}

	// 获取待推送的消息
	messages, err := messageRepo.FindPendingByReceiver(userID)
	if err != nil {
		log.Printf("Failed to fetch offline messages for user %d: %v", userID, err)
		return
	}

	if len(messages) > 0 {
		// 推送消息
		BroadcastOfflineMessage(GlobalHub, userID, messages)

		// 标记为已送达
		messageIDs := make([]int64, len(messages))
		for i, msg := range messages {
			messageIDs[i] = msg.ID
		}
		if err := messageRepo.MarkAsDelivered(userID, messageIDs); err != nil {
			log.Printf("Failed to mark messages as delivered: %v", err)
		}
	}
}
