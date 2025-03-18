package loader

// Config struct for YAML parsing
type Config struct {
	Server struct {
		Port           int      `yaml:"port"`
		Timeout        string   `yaml:"timeout"`
		AllowCORS      bool     `yaml:"allow_cors"`
		AllowedOrigins []string `yaml:"allowed_origins"`
		Cookie         struct {
			Secure   bool   `yaml:"secure"`
			HTTPOnly bool   `yaml:"http_only"`
			SameSite string `yaml:"same_site"`
		} `yaml:"cookie"`
	} `yaml:"server"`

	Framework string `yaml:"framework"`

	Middlewares struct {
		Global []string `yaml:"global"`
	} `yaml:"middlewares"`

	Database struct {
		Type               string `yaml:"type"`                 // Database type (postgres, mysql, sqlite3)
		Host               string `yaml:"host"`                 // Hostname (not needed for SQLite)
		Port               string `yaml:"port"`                 // Port (not needed for SQLite)
		Username           string `yaml:"username"`             // DB username
		Password           string `yaml:"password"`             // DB password
		Database           string `yaml:"database"`             // DB name or SQLite file path
		SSLMode            string `yaml:"sslmode"`              // SSL mode for PostgreSQL
		MaxIdleConnections int    `yaml:"max_idle_connections"` // Max idle connections
		MaxOpenConnections int    `yaml:"max_open_connections"` // Max open connections
		ConnMaxLifetime    string `yaml:"conn_max_lifetime"`    // Connection max lifetime (e.g., "1h")
	} `yaml:"database"`

	Routes struct {
		Groups []struct {
			Base   string `yaml:"base"`
			Routes []struct {
				Path       string `yaml:"path"`
				Handler    string `yaml:"handler"`
				Method     string `yaml:"method"`
				BodyParams []struct {
					Name       string                 `yaml:"name"`
					Type       string                 `yaml:"type"`
					Validation map[string]interface{} `yaml:"validation"`
				} `yaml:"body_params"`
			} `yaml:"routes"`
		} `yaml:"groups"`
	} `yaml:"routes"`
}
