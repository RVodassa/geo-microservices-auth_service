package main

import (
	"fmt"
	"github.com/RVodassa/geo-microservices-auth_service/app"
	"github.com/RVodassa/geo-microservices-auth_service/app/config"
	myLogger "github.com/RVodassa/slog_utils/slog_logger"
	"github.com/joho/godotenv"
	"log"
	"log/slog"
	"os"
)

// TODO: Gracefully SD for Serve
// TODO: Refactoring project
// TODO: структурные логи
// TODO: Tests

const (
	defaultConfigPath = "/app/auth-service/configs/config.yaml" // Путь в контейнере
	configEnvVar      = "CONFIG_PATH"
)

func main() {
	const op = "main.main"

	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("%s: ошибка при загрузке cfg %v", op, err)
	}

	logger := setupLogger(cfg)
	logger.Info("Slog инициализирован", "op", op, "LogLevel", cfg.LogLevel, "Env", cfg.Env)

	// Инициализация и запуск приложения
	newApp := app.NewApp(logger, cfg)
	if err = newApp.Run(); err != nil {
		logger.Error(op, err)
		os.Exit(1)
	}
}

func loadConfig() (*config.Config, error) {
	const op = "main.loadConfig"

	if err := godotenv.Load(".env"); err != nil {
		return nil, fmt.Errorf("ошибка загрузки .env файла: %w", err)
	}

	configPath := os.Getenv(configEnvVar)
	if configPath == "" {
		configPath = defaultConfigPath
		log.Printf("%s: Путь к конфигурации по умолчанию: %s", op, configPath)
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("%s:failed to load config: %w", op, err)
	}
	return cfg, nil
}

func setupLogger(cfg *config.Config) *slog.Logger {
	var logLevel slog.Level

	switch cfg.LogLevel {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	logger := myLogger.SetupLogger(cfg.Env, logLevel)
	return logger
}
