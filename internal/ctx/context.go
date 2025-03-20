package ctx

import (
	framework "github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// NewContext creates a new `routeflow.Context` from Gin context
func NewContext(c *framework.Context, db *gorm.DB) *Context {
	// Extract path parameters
	pathParams := make(map[string]string)
	for _, param := range c.Params {
		pathParams[param.Key] = param.Value
	}

	// Extract query parameters
	queryParams := make(map[string]string)
	for key, value := range c.Request.URL.Query() {
		queryParams[key] = value[0] // Take first value
	}

	// Extract JSON body data
	bodyData := make(map[string]interface{})
	if c.Request.Method == "POST" || c.Request.Method == "PUT" {
		if err := c.ShouldBindJSON(&bodyData); err != nil {
			bodyData = map[string]interface{}{"error": "Invalid JSON"}
		}
	}

	return &Context{
		GinContext:  c,
		PathParams:  pathParams,
		QueryParams: queryParams,
		BodyData:    bodyData,
		DB:          db,
	}
}
