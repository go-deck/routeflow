package routeflow

import (
	"log"

	"github.com/go-deck/routeflow/ctx"
	"github.com/go-deck/routeflow/db"
	base "github.com/go-deck/routeflow/frameworks"
	ginserver "github.com/go-deck/routeflow/frameworks/ginserver"
	"github.com/go-deck/routeflow/loader"
	"gorm.io/gorm"
)

// App represents the main application structure
type App struct {
	Config *loader.Config
	DB     *gorm.DB
}

// New creates a new instance of App
func New(configPath string) (*App, error) {
	cfg, err := loader.LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	return &App{Config: cfg}, nil
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
