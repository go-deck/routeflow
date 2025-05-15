package ctx

import (
	framework "github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const ContextKey = "RF_CONTEXT"

// Get (update existing NewContext to this)
func Get(c *framework.Context) *Context {
	return c.MustGet(ContextKey).(*Context)
}

// Initialize middleware (add this new function)
func Middleware(db *gorm.DB) framework.HandlerFunc {
	return func(c *framework.Context) {
		// Create context once per request
		cc := &Context{
			GinContext:  c,
			PathParams:  make(map[string]string),
			QueryParams: make(map[string]string),
			BodyData:    make(map[string]interface{}),
			DB:          db,
		}

		// Parse data once
		for _, param := range c.Params {
			cc.PathParams[param.Key] = param.Value
		}

		for key, values := range c.Request.URL.Query() {
			cc.QueryParams[key] = values[0]
		}

		if c.Request.Method == "POST" || c.Request.Method == "PUT" {
			if err := c.ShouldBindJSON(&cc.BodyData); err != nil {
				cc.BodyData = map[string]interface{}{"error": "Invalid JSON"}
			}
		}

		c.Set(ContextKey, cc)
		c.Next()
	}
}

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
