package ginserver

import (
	"github.com/gin-contrib/cors"
	framework "github.com/gin-gonic/gin"
	"github.com/go-deck/routeflow/internal/ctx"
	"github.com/go-deck/routeflow/internal/loader"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Initialize Gin router with the provided configuration and handler map
func InitGinRouter(r *framework.Engine, cfg *loader.Config, handlerMap map[string]ctx.HandlerFunc, db *gorm.DB, middlewareMap map[string]ctx.HandlerFunc) {
	for _, group := range cfg.Routes.Groups {
		api := r.Group(group.Base)
		for _, mw := range group.Middlewares.BuiltIn {
			switch mw {
			case "cors":
				corsConfig := cors.DefaultConfig()
				corsConfig.AllowOrigins = cfg.Server.AllowedOrigins
				api.Use(cors.New(corsConfig))
			case "logging":
				api.Use(framework.Logger())
			case "recovery":
				api.Use(framework.Recovery())
			}
		}
		for _, name := range group.Middlewares.Custom {
			if mw, exists := middlewareMap[name]; exists {
				api.Use(MiddlewareAdapter(db, mw))
			}
		}
		// Inside InitGinRouter function...
		for _, route := range group.Routes {
			handlers, exists := handlerMap[route.Handler]
			if !exists {
				log.WithFields(log.Fields{
					"route":   route.Path,
					"handler": route.Handler,
				}).Error("Handler not found")
				continue
			}

			// Collect route-specific middlewares
			var routeHandlers []framework.HandlerFunc

			// Add built-in middlewares for this route
			for _, mw := range route.Middlewares.BuiltIn {
				switch mw {
				case "cors":
					corsConfig := cors.DefaultConfig()
					corsConfig.AllowOrigins = cfg.Server.AllowedOrigins
					routeHandlers = append(routeHandlers, cors.New(corsConfig))
				case "logging":
					routeHandlers = append(routeHandlers, framework.Logger())
				case "recovery":
					routeHandlers = append(routeHandlers, framework.Recovery())
				}
			}

			// Add custom middlewares for this route
			for _, name := range route.Middlewares.Custom {
				if mw, exists := middlewareMap[name]; exists {
					routeHandlers = append(routeHandlers, MiddlewareAdapter(db, mw))
				}
			}

			// Add the actual handler as the final step
			routeHandlers = append(routeHandlers, MiddlewareAdapter(db, handlers))

			// Register the route with collected middlewares and handler
			api.Handle(route.Method, route.Path, routeHandlers...)
		}
	}
}
