package ginserver

import (
	"github.com/gin-contrib/cors"
	framework "github.com/gin-gonic/gin"
	"github.com/go-deck/routeflow/internal/ctx"
	"github.com/go-deck/routeflow/internal/handler"
	"github.com/go-deck/routeflow/internal/loader"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Initialize Gin router with the provided configuration and handler map
func InitGinRouter(r *framework.Engine, cfg *loader.Config, handlerMap map[string]func(*ctx.Context) (interface{}, int), db *gorm.DB, middlewareMap map[string]func(*ctx.Context) (interface{}, int)) {
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
		for _, route := range group.Routes {
			handlers, exists := handlerMap[route.Handler]

			// middleware
			for _, mw := range route.Middlewares.BuiltIn {
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
			for _, name := range route.Middlewares.Custom {
				if mw, exists := middlewareMap[name]; exists {
					api.Use(MiddlewareAdapter(db, mw))
				}
			}

			if !exists {
				log.WithFields(log.Fields{
					"route":   route.Path,
					"handler": route.Handler,
				}).Error("Handler not found")
				continue
			}

			// Extract validation properties from YAML
			validation := make(map[string]interface{})
			if route.BodyParams != nil {
				for _, param := range route.BodyParams {
					validation[param.Name] = param.Validation
				}
			}

			// Convert to a Gin-compatible handler
			wrappedHandler := handler.WrapHandler(handlers, validation, db)

			// Register routes dynamically
			api.Handle(route.Method, route.Path, wrappedHandler)
		}
	}
}
