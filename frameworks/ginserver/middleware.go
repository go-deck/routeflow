package ginserver

import (
	"github.com/gin-contrib/cors"
	framework "github.com/gin-gonic/gin"
	"github.com/go-deck/routeflow/internal/ctx"
	"github.com/go-deck/routeflow/internal/loader"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Change to accept ctx.HandlerFunc explicitly
func MiddlewareAdapter(db *gorm.DB, fn ctx.HandlerFunc) framework.HandlerFunc {
	return func(c *framework.Context) {
		cc := ctx.Get(c) // Use shared context
		response, status := fn(cc)
		if status > 0 {
			c.AbortWithStatusJSON(status, response)
		} else {
			c.Next()
		}
	}
}

func LoadMiddlewares(r *framework.Engine, cfg *loader.Config, middlewareMap map[string]ctx.HandlerFunc, db *gorm.DB) {
	// Built-in middlewares
	for _, mw := range cfg.Middlewares.BuiltIn {
		switch mw {
		case "cors":
			corsConfig := cors.DefaultConfig()
			corsConfig.AllowOrigins = cfg.Server.AllowedOrigins
			r.Use(cors.New(corsConfig))
		case "logging":
			r.Use(framework.Logger())
		case "recovery":
			r.Use(framework.Recovery())
		}
	}

	// Custom middlewares
	for _, name := range cfg.Middlewares.Custom {
		if mw, exists := middlewareMap[name]; exists {
			r.Use(MiddlewareAdapter(db, mw))
		} else {
			log.Printf("⚠️ Middleware not found: %s", name)
		}
	}
}
