package ginframework

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-deck/routeflow/loader"
)

// StartGinServer initializes the Gin server
func StartGinServer(cfg *loader.Config, handlerMap map[string]func(*Context) (interface{}, int)) {
	r := gin.New()

	// Set Gin to release mode
	gin.SetMode(gin.ReleaseMode)

	// Load middleware
	LoadMiddlewares(r, cfg)

	// Register routes
	InitGinRouter(r, cfg, handlerMap)

	// Start server
	port := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Server running on %s", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
