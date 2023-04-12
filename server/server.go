package server

import (
	"backend/handlers"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func StartServer(addr string) {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Server is running"})
	})

	router.GET("/ws", func(c *gin.Context) {
		handlers.UpgradeToWebSocket(c.Writer, c.Request)
	})

	router.POST("/start", func(c *gin.Context) {
		// Placeholder for starting a chat session
	})

	router.POST("/leave", func(c *gin.Context) {
		// Placeholder for leaving a chat session
	})

	log.Printf("Starting server at %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
