package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// LogMiddleware 作为中间件记录请求的日志
func LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		log.Printf("Request to %s took %s", c.Request.URL.Path, duration)
	}
}

func main() {
	r := gin.Default()

	// 使用 LogMiddleware
	r.Use(LogMiddleware())

	r.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	r.Run(":8080")
}
