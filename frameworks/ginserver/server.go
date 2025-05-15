package ginserver

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	framework "github.com/gin-gonic/gin"
	"github.com/go-deck/routeflow/internal/ctx"
	"github.com/go-deck/routeflow/internal/loader"
	"gorm.io/gorm"
)

type GinServer struct{}

func (g *GinServer) Start(cfg *loader.Config, handlerMap map[string]ctx.HandlerFunc, db *gorm.DB, middlewareMap map[string]ctx.HandlerFunc) {
	r := framework.New()

	r.Use(ctx.Middleware(db))

	framework.SetMode(framework.ReleaseMode)

	// Load middleware
	LoadMiddlewares(r, cfg, middlewareMap, db)

	// Register routes
	InitGinRouter(r, cfg, handlerMap, db, middlewareMap)

	// Start server
	port := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Server running on %s", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
