package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	YagptKey           string  `yaml:"yagpt-key"`
	YagptModel         string  `yaml:"yagpt-model"`
	YagptTemperature   float32 `yaml:"yagpt-temperature"`
	YagptMaxTokens     int32   `yaml:"yagpt-max-tokens"`
	YagptSystemMessage string  `yaml:"yagpt-system-message"`

	TgToken string `yaml:"tg-token"`

	DbDSN string `yaml:"db-dsn"`
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
