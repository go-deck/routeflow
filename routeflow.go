package routeflow

import (
	log "github.com/sirupsen/logrus"

	base "github.com/go-deck/routeflow/frameworks"
	ginserver "github.com/go-deck/routeflow/frameworks/ginserver"
	"github.com/go-deck/routeflow/internal/ctx"
	"github.com/go-deck/routeflow/internal/db"
	"github.com/go-deck/routeflow/internal/loader"
	"github.com/go-deck/routeflow/internal/utils"
	"gorm.io/gorm"
)

// App represents the main application structure
type App struct {
	Config        *loader.Config
	DB            *gorm.DB
	HandlerMap    map[string]ctx.HandlerFunc
	MiddlewareMap map[string]ctx.HandlerFunc
}

type Context = ctx.Context

// Register middleware struct instances
func (app *App) RegisterMiddlewareStructs(middlewares ...interface{}) {
	discoverMiddleware := utils.DiscoverHandlers(middlewares...)
	app.MiddlewareMap = utils.MapMiddlewareHandlersFromYAML(app.Config, discoverMiddleware)
}

func New(configPath string, handlerStructs ...interface{}) (*App, error) {
	cfg, err := loader.LoadConfig(configPath)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to load config")
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
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to connect to database")
		return err
	}

	app.DB = dbConn
	return nil
}

// Serve starts the API server with the given handler mappings
func (app *App) Serve() {
	var server base.Server
	switch app.Config.Framework {
	case "gin":
		server = &ginserver.GinServer{}
	default:
		log.Fatalf("Unsupported framework: %s", app.Config.Framework)
	}

	server.Start(app.Config, app.HandlerMap, app.DB, app.MiddlewareMap)
}
