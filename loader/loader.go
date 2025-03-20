package loader

import (
	"os"

	log "github.com/sirupsen/logrus"

	"gopkg.in/yaml.v3"
)

// LoadConfig loads the YAML file into a Config struct
func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to load config file")
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to unmarshal config file")
		return nil, err
	}

	return &config, nil
}
