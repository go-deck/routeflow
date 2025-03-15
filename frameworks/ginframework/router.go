package ginframework

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-deck/routeflow/loader"
)

// Initialize Gin router with the provided configuration and handler map
func InitGinRouter(r *gin.Engine, cfg *loader.Config, handlerMap map[string]func(map[string]string, map[string]string, map[string]interface{}) (interface{}, int)) {
	for _, group := range cfg.Routes.Groups {
		api := r.Group(group.Base)
		for _, route := range group.Routes {
			handler, exists := handlerMap[route.Handler]
			if !exists {
				log.Fatalf("Handler not found for route: %s", route.Handler)
			}

			// Convert handler to the correct type
			wrappedHandler := convertToGinHandler(handler)

			// Register routes dynamically
			api.Handle(route.Method, route.Path, wrappedHandler)
		}
	}
}

// Extract parameters from the request and store them in a generic map
func extractParams(c *gin.Context) map[string]interface{} {
	params := make(map[string]interface{})

	// Extract path parameters
	pathParams := make(map[string]string)
	for _, param := range c.Params {
		pathParams[param.Key] = param.Value
	}
	params["pathParams"] = pathParams

	// Extract query parameters
	queryParams := make(map[string]string)
	for key, value := range c.Request.URL.Query() {
		queryParams[key] = value[0] // Take the first value
	}
	params["queryParams"] = queryParams

	// Extract JSON body data (if applicable)
	bodyData := make(map[string]interface{})
	if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut {
		if err := c.ShouldBindJSON(&bodyData); err != nil {
			log.Printf("Error parsing request body: %v", err)
		}
	}
	params["bodyData"] = bodyData

	return params
}

// Convert generic handler to Gin's `gin.HandlerFunc`
func convertToGinHandler(handler interface{}) gin.HandlerFunc {
	h, ok := handler.(func(map[string]string, map[string]string, map[string]interface{}) (interface{}, int))
	if !ok {
		log.Fatalf("Invalid handler function type")
	}

	return func(c *gin.Context) {
		// Extract request parameters
		params := extractParams(c)

		// Call user-defined handler
		response, statusCode := h(
			params["pathParams"].(map[string]string),
			params["queryParams"].(map[string]string),
			params["bodyData"].(map[string]interface{}),
		)

		// Send response
		c.JSON(statusCode, response)
	}
}
