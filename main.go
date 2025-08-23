package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ash543210/go-chat/logger"
	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // Step1: Process the request first.

		// Step2: Check if any errors were added to the context
		if len(c.Errors) > 0 {
			// Step3: Use the last error
			err := c.Errors.Last().Err
			c.Error(errors.New(err.Error()))
			// Step4: Respond with a generic error message
			c.JSON(http.StatusInternalServerError, map[string]any{
				"success": false,
				"message": err.Error(),
			})
		}

		// Any other steps if no errors are found
	}
}

func main() {
	LOGGER, err := logger.New(logger.DEBUG, true, "")
	if err != nil {
		fmt.Println("Logger initializtion failed!")
		return
	}

	defer LOGGER.Close()

	router := gin.Default()

	router.Use(ErrorHandler())

	router.GET("/someJSON", func(c *gin.Context) {
		data := map[string]interface{}{
			"lang": "GO语言",
			"tag":  "<br>",
		}

		// will output : {"lang":"GO\u8bed\u8a00","tag":"\u003cbr\u003e"}
		c.AsciiJSON(200, data)
	})

	// Listen and serve on 0.0.0.0:8080
	router.Run(":8080")
}
