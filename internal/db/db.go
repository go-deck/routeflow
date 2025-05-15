package db

import (
	"fmt"
	"log"
	"time"

	"github.com/go-deck/routeflow/internal/loader"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ConnectDB initializes the database connection based on YAML configuration.
func ConnectDB(cfg *loader.Config) (*gorm.DB, error) {
	dbConfig := cfg.Database
	var dsn string
	var dialector gorm.Dialector

	// Choose database type
	switch dbConfig.Type {
	case "postgres":
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			dbConfig.Host, dbConfig.Username, dbConfig.Password, dbConfig.Database, dbConfig.Port, dbConfig.SSLMode,
		)
		dialector = postgres.Open(dsn)

	case "mysql":
		dsn = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Database,
		)
		dialector = mysql.Open(dsn)

	case "sqlite3":
		dsn = dbConfig.Database // SQLite uses just the file path
		dialector = sqlite.Open(dsn)

	default:
		return nil, fmt.Errorf("❌ Unsupported database type: %s", dbConfig.Type)
	}

	// Open database connection
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Enable logging for queries
	})
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to connect to %s: %v", dbConfig.Type, err)
	}

	// Get the underlying SQL database for connection settings
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to get SQL database: %v", err)
	}

	// Set database connection pool settings
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConnections)
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConnections)

	connMaxLifetime, err := time.ParseDuration(dbConfig.ConnMaxLifetime)
	if err == nil {
		sqlDB.SetConnMaxLifetime(connMaxLifetime)
	}

	log.Printf("✅ Database connected: %s", dbConfig.Type)
	return db, nil
}
