package ginframework

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-deck/routeflow/loader"
)

// Initialize Gin router with the provided configuration and handler map
func InitGinRouter(r *gin.Engine, cfg *loader.Config, handlerMap map[string]func(*Context) (interface{}, int)) {
	for _, group := range cfg.Routes.Groups {
		api := r.Group(group.Base)
		for _, route := range group.Routes {
			handler, exists := handlerMap[route.Handler]
			if !exists {
				log.Fatalf("Handler not found for route: %s", route.Handler)
			}

			// Convert to a Gin-compatible handler
			wrappedHandler := wrapHandler(handler)

			// Register routes dynamically
			api.Handle(route.Method, route.Path, wrappedHandler)
		}
	}
}

// Convert user-defined handler to a Gin-compatible `gin.HandlerFunc`
func wrapHandler(userHandler func(*Context) (interface{}, int)) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Convert Gin context to routeflow.Context
		ctx := NewContext(c)

		// Call the user handler
		response, statusCode := userHandler(ctx)

		// Send response
		c.JSON(statusCode, response)
	}
}
