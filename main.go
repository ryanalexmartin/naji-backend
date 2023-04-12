package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	addr := "localhost:8080"
	startServer(addr)
}

func startServer(addr string) {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Server is running"})
	})

	log.Printf("Starting server at %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

