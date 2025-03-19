package ctx

import (
	framework "github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Context wraps request data
type Context struct {
	GinContext  *framework.Context     // Raw Gin context (if needed)
	PathParams  map[string]string      // Path parameters
	QueryParams map[string]string      // Query parameters
	BodyData    map[string]interface{} // JSON body data
	DB          *gorm.DB               // Database connection
}
