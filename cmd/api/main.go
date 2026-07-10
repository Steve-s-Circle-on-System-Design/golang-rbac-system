// Command rbac-system is the entrypoint for the RBAC service.
package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/health", healthHandler)

	addr := ":8080"
	log.Printf("server listening on %s", addr)

	if err := router.Run(addr); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
