package main

import (
	"github.com/RVodassa/geo-microservices-auth_service/app"
	"log"
	"os"
)

// TODO: Gracefully SD for Serve
// TODO: Refactoring project
// TODO: Slog
// TODO: Tests

const (
	defaultConfigPath = "/app/auth-service/configs/config.yaml" // Путь в контейнере
	configEnvVar      = "CONFIG_PATH"
)

func main() {
	configPath := os.Getenv(configEnvVar)
	if configPath == "" {
		configPath = defaultConfigPath
		log.Printf("Используется конфигурация по умолчанию: %s", configPath)
	}

	err := app.RunApp(configPath)
	if err != nil {
		log.Fatalf("Ошибка при запуске приложения: %v", err)
	}
}
