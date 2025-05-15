package frameworks

import (
	"github.com/go-deck/routeflow/internal/ctx"
	"github.com/go-deck/routeflow/internal/loader"
	"gorm.io/gorm"
)

type Server interface {
	Start(cfg *loader.Config, handlerMap map[string]ctx.HandlerFunc, db *gorm.DB, middlewareMap map[string]ctx.HandlerFunc)
}
