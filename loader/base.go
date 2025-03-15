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

	Routes struct {
		Groups []struct {
			Base   string `yaml:"base"`
			Routes []struct {
				Path       string `yaml:"path"`
				Handler    string `yaml:"handler"`
				Method     string `yaml:"method"`
				BodyParams []struct {
					Name string `yaml:"name"`
					Type string `yaml:"type"`
				} `yaml:"body_params"`
			} `yaml:"routes"`
		} `yaml:"groups"`
	} `yaml:"routes"`
}
