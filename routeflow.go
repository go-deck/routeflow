package routeflow

import (
	"log"

	"github.com/go-deck/routeflow/ctx"
	"github.com/go-deck/routeflow/db"
	base "github.com/go-deck/routeflow/frameworks"
	ginserver "github.com/go-deck/routeflow/frameworks/ginserver"
	"github.com/go-deck/routeflow/loader"
	"github.com/go-deck/routeflow/utils"
	"gorm.io/gorm"
)

// App represents the main application structure
type App struct {
	Config     *loader.Config
	DB         *gorm.DB
	HandlerMap map[string]func(*ctx.Context) (interface{}, int)
}

type Context = ctx.Context

func New(configPath string, handlerStructs ...interface{}) (*App, error) {
	cfg, err := loader.LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	app := &App{Config: cfg}

	discoveredHandlers := utils.DiscoverHandlers(handlerStructs...)
	app.HandlerMap = utils.MapHandlersFromYAML(cfg, discoveredHandlers)

	return app, nil
}

// InitDB initializes the database connection
func (app *App) InitDB() error {
	dbConn, err := db.ConnectDB(app.Config)
	if err != nil {
		return err
	}
	app.DB = dbConn
	return nil
}

// Serve starts the API server with the given handler mappings
func (app *App) Serve(handlerMap map[string]func(*ctx.Context) (interface{}, int)) {
	var server base.Server
	switch app.Config.Framework {
	case "gin":
		server = &ginserver.GinServer{}
	default:
		log.Fatalf("Unsupported framework: %s", app.Config.Framework)
	}

	server.Start(app.Config, handlerMap, app.DB)
}
