package handlers

import (
	"net/http"

	"github.com/ash543210/go-chat/internal/helper"
	"github.com/ash543210/go-chat/internal/services/auth"
	"github.com/ash543210/go-chat/internal/services/user"
	"github.com/ash543210/go-chat/pkg/logger"
	"github.com/gin-gonic/gin"
)

type SignupPayload struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
}

func GetUser(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"user_id": id,
		"name":    "John Doe",
	})
}

func CreateUser(c *gin.Context) {
	l, exists := c.Get("logger")
	if !exists {
		c.JSON(500, gin.H{
			"error": "Logger not found in context",
		})
		return
	}

	logger, ok := l.(*logger.Logger)
	if !ok {
		c.JSON(500, gin.H{
			"error": "Logger type assertion failed",
		})
		return
	}
	var payload SignupPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid input: " + err.Error(),
		})
	}
	logger.Info(c, "Creating user: %s", payload.Email)
	userService := user.UserService{Logger: logger}
	UserDetails, userID, err := userService.CreateUser(payload.FirstName, payload.LastName, payload.Email, payload.Password, c.Request.Context())
	if err != nil {
		helper.RespondError(c, 500, "User creation failed: "+err.Error())
		return
	}
	authService := (auth.AuthService{Logger: logger})
	_, refreshToken, err := authService.GenerateAuthTokens(userID, nil)
	if err != nil {
		helper.RespondError(c, 500, "Token generation failed: "+err.Error())
		return
	}
	err = authService.InsertRefreshToken(userID, refreshToken, c.Request.Context())
	if err != nil {
		helper.RespondError(c, 500, "Storing refresh token failed: "+err.Error())
		return
	}
	c.JSON(201, gin.H{
		"message": "User created successfully",
		"user":    UserDetails,
	})
}
