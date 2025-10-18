package routes

import (
	"github.com/ash543210/go-chat/internal/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.GET("/health", handlers.HealthCheck)
		api.GET("/user/:id", handlers.GetUser)
		api.POST("/signup", handlers.CreateUser)
	}
}
