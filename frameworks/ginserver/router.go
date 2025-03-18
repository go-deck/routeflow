package ginserver

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-deck/routeflow/ctx"
	"github.com/go-deck/routeflow/handler"
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
			validation := make(map[string]interface{})
			if route.BodyParams != nil {
				for _, param := range route.BodyParams {
					validation[param.Name] = param.Validation
				}
			}

			// Convert to a Gin-compatible handler
			wrappedHandler := handler.WrapHandler(userHandler, validation, db)

			// Register route dynamically
			api.Handle(route.Method, route.Path, wrappedHandler)
		}
	}
}
