package ginframework

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-deck/routeflow/frameworks/ginframework/ctx"
	"github.com/go-deck/routeflow/frameworks/ginframework/handler"
	"github.com/go-deck/routeflow/loader"
	"gorm.io/gorm"
)

// InitGinRouter initializes the Gin router
func InitGinRouter(r *gin.Engine, cfg *loader.Config, handlerMap map[string]func(*ctx.Context) (interface{}, int), db *gorm.DB) {
	for _, group := range cfg.Routes.Groups {
		api := r.Group(group.Base)
		for _, route := range group.Routes {
			userHandler, exists := handlerMap[route.Handler]
			if !exists {
				log.Fatalf("Handler not found for route: %s", route.Handler)
			}

			// Extract validation properties from YAML
			props := make(map[string]interface{})
			if route.BodyParams != nil {
				for _, param := range route.BodyParams {
					props[param.Name] = param.Props
				}
			}

			// Convert to a Gin-compatible handler
			wrappedHandler := handler.WrapHandler(userHandler, props, db)

			// Register route dynamically
			api.Handle(route.Method, route.Path, wrappedHandler)
		}
	}
}
