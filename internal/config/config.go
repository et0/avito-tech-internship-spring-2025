package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	HTTP HTTP     `yaml:"http_server"`
	DB   Database `yaml:"database"`
	GRPC GRPC     `yaml:"grpc"`
}

type HTTP struct {
	Port      string `yaml:"port"`
	JWTSecret string `yaml:"jwt_secret"`
}

type GRPC struct {
	Port string `yaml:"port"`
}

type Database struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Basename string `yaml:"basename"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func Load() (*Config, error) {
	configPath, exists := os.LookupEnv("CONFIG_PATH")
	if !exists || configPath == "" {
		return nil, fmt.Errorf("env CONFIG_PATH is not set")
	}

	cfg := &Config{}

	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %s", err)
	}

	if err := yaml.Unmarshal(file, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %s", err)
	}

	return cfg, nil
}
