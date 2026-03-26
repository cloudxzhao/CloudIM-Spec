package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

// Logger 日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		log.Printf("[%d] %s %s %v", statusCode, path, query, latency)
	}
}

// Recovery 恢复中间件
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				c.AbortWithStatusJSON(500, gin.H{
					"code":    50001,
					"message": "internal server error",
				})
			}
		}()
		c.Next()
	}
}
