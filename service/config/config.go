package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	GRPC struct {
		Address string `yaml:"address"`
	} `yaml:"grpc"`
}

func Load() (*Config, error) {
	data, err := os.ReadFile("config/config.yaml")
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
