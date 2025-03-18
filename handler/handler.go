package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-deck/routeflow/ctx"
	"github.com/go-deck/routeflow/validator"
	"gorm.io/gorm"
)

// Convert user-defined handler to a Gin-compatible handler
func WrapHandler(userHandler func(*ctx.Context) (interface{}, int), validation map[string]interface{}, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		context := ctx.NewContext(c, db)

		// Validate body params if validation exist
		if len(validation) > 0 {
			err := validator.ValidateBody(context.BodyData, validation)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}

		// Call the actual user-defined handler
		response, statusCode := userHandler(context)

		// Send response
		c.JSON(statusCode, response)
	}
}
