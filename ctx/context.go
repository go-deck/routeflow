package ctx

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Context wraps request data
type Context struct {
	GinContext  *gin.Context           // Raw Gin context (if needed)
	PathParams  map[string]string      // Path parameters
	QueryParams map[string]string      // Query parameters
	BodyData    map[string]interface{} // JSON body data
	DB          *gorm.DB               // Database connection
}

// NewContext creates a new `routeflow.Context` from Gin context
func NewContext(c *gin.Context, db *gorm.DB) *Context {
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
