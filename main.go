package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ash543210/go-chat/internal/middleware"
	"github.com/ash543210/go-chat/internal/routes"
	"github.com/ash543210/go-chat/mongo"
	"github.com/ash543210/go-chat/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
	LOGGER, err := logger.New(logger.DEBUG, true, "./logs/test.txt")
	if err != nil {
		fmt.Println("Logger initializtion failed!")
		return
	}
	LOGGER.Debug(context.Background(), "Logger initialized.")

	defer mongo.Disconnect()

	envErr := godotenv.Load(".env")
	if envErr != nil {
		LOGGER.Error(context.Background(), "Error loading .env file")
		return
	}

	defer LOGGER.Close()

	router := gin.Default()

	router.Use(ErrorHandler())
	router.Use(middleware.LoggerMiddleware(LOGGER))

	routes.RegisterRoutes(router)

	// Listen and serve on 0.0.0.0:8080
	srvAddr := ":8080"
	go func() {
		if err := router.Run(srvAddr); err != nil {
			LOGGER.Error(context.Background(), "Server failed to start: %v", err)
		}
	}()
	fmt.Printf("Server started at %s\n", srvAddr)

	// Wait for interrupt signal to gracefully shut down
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Shutting down server gracefully...")
}
