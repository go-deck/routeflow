package routeflow

import (
	"log"

	"github.com/go-deck/routeflow/db"
	"github.com/go-deck/routeflow/frameworks/ginframework"
	"github.com/go-deck/routeflow/frameworks/ginframework/ctx"
	"github.com/go-deck/routeflow/loader"
)

// StartServer initializes the correct framework
func StartServer(configPath string, handlerMap map[string]func(*ctx.Context) (interface{}, int)) {
	cfg, err := loader.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	dbConn, err := db.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	switch cfg.Framework {
	case "gin":
		ginframework.StartGinServer(cfg, handlerMap, dbConn)
	default:
		log.Fatalf("Unsupported framework: %s", cfg.Framework)
	}
}
