package config

import (
	"log"
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

func Load() *Config {
	configPath, exists := os.LookupEnv("CONFIG_PATH")
	if !exists || configPath == "" {
		log.Fatalf("env CONFIG_PATH is not set")
	}

	cfg := &Config{}

	file, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("failed to read config file: %s", err)
	}

	if err := yaml.Unmarshal(file, cfg); err != nil {
		log.Fatalf("failed to parse config file: %s", err)
	}

	return cfg
}
