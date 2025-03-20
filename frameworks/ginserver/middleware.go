package ginserver

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-deck/routeflow/internal/loader"
)

// LoadMiddlewares applies middlewares dynamically
func LoadMiddlewares(r *gin.Engine, cfg *loader.Config) {
	for _, mw := range cfg.Middlewares.Global {
		switch mw {
		case "cors":
			corsConfig := cors.DefaultConfig()
			corsConfig.AllowOrigins = cfg.Server.AllowedOrigins
			r.Use(cors.New(corsConfig))
		case "logging":
			r.Use(gin.Logger())
		case "recovery":
			r.Use(gin.Recovery())
		}
	}
}
