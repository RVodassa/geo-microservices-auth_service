package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	GRPCServices []GRPCService `yaml:"grpc_services"`
	Env          string        `yaml:"env"`
	LogLevel     string        `yaml:"log_level"`
	ServePort    string        `yaml:"serve_port"`
}

type GRPCService struct {
	Name string `yaml:"name"`
	Addr string `yaml:"addr"`
	Port int    `yaml:"port"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения конфигурации: %v", err)
	}

	var cfg Config
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("ошибка парсинга конфигурации: %v", err)
	}

	log.Printf("Загружена конфигурация: %+v", cfg)
	return &cfg, nil
}
