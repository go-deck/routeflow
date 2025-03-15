package loader

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// LoadConfig loads the YAML file into a Config struct
func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Error parsing YAML: %v", err)
		return nil, err
	}

	return &config, nil
}
