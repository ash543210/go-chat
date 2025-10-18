package helper

import "github.com/gin-gonic/gin"

func RespondError(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{"error": msg})
}
