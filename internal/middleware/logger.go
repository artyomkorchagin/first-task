package middleware

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path
		statusCode := c.Writer.Status()

		msg := fmt.Sprintf(
			"[GIN] %s | %3d | %13v | %15s | %-7s %s",
			time.Now().Format("02.01.2006 - 15:04:05"),
			statusCode,
			latency,
			clientIP,
			method,
			path,
		)
		log.Println(msg)
	}
}
