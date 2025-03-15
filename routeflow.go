package routeflow

import (
	"log"

	"github.com/go-deck/routeflow/frameworks/ginframework"
	"github.com/go-deck/routeflow/loader"
)

// StartServer initializes the correct framework
func StartServer(configPath string, handlerMap map[string]func(map[string]string, map[string]string, map[string]interface{}) (interface{}, int)) {
	cfg, err := loader.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	switch cfg.Framework {
	case "gin":
		ginframework.StartGinServer(cfg, handlerMap)
	default:
		log.Fatalf("Unsupported framework: %s", cfg.Framework)
	}
}
