package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	YagptKey string `yaml:"yagpt-key"`
}

func LoadConfig(path string) (*Config, error) {
	log.Printf("Loading config from %s", path)

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	errUm := yaml.Unmarshal(data, &config)
	if errUm != nil {
		return nil, errUm
	}

	return &config, nil
}
