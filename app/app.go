package app

import (
	"fmt"
	"github.com/RVodassa/geo-microservices-auth_service/internal/serve"
	serviceAuth "github.com/RVodassa/geo-microservices-auth_service/internal/service"
	serviceUser "github.com/RVodassa/geo-microservices-user_service/proto/generated"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
)

const UserServiceName = "user-service"

func RunApp(configPath string) error {
	// Загрузка .env файла
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Ошибка загрузки .env файла: %v", err)
		return fmt.Errorf("ошибка загрузки .env файла: %v", err)
	}

	// Проверка обязательных переменных окружения
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return fmt.Errorf("JWT_SECRET не установлен в .env файле")
	}

	log.Println("Загрузка конфигурации...")
	cfg, err := LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("ошибка загрузки конфигурации: %v", err)
	}
	log.Printf("Конфигурация загружена: %+v", cfg)

	log.Println("Установка соединений с gRPC сервисами...")
	grpcConns, err := setupGRPCClient(cfg.GRPCServices)
	if err != nil {
		return fmt.Errorf("ошибка подключения к gRPC сервисам: %v", err)
	}

	defer func() {
		log.Println("Закрытие gRPC соединений...")
		for name, conn := range grpcConns {
			if err = conn.Close(); err != nil {
				log.Printf("Ошибка при закрытии соединения с %s: %v", name, err)
			}
		}
	}()

	client := serviceUser.NewUserServiceClient(grpcConns[UserServiceName])
	log.Println("Подключение к serviceUser успешно установлено")

	authService := serviceAuth.NewAuthService(client, jwtSecret)
	log.Println("Сервис авторизации успешно создан")

	if err = serve.RunServe(authService, cfg.ServePort); err != nil {
		return fmt.Errorf("ошибка при запуске сервера: %v", err)
	}

	return nil
}

func setupGRPCClient(grpcServices []GRPCService) (map[string]*grpc.ClientConn, error) {

	conns := make(map[string]*grpc.ClientConn) // service_name:client

	for _, grpcServ := range grpcServices {
		addr := fmt.Sprintf("%s:%d", grpcServ.Addr, grpcServ.Port)
		conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			for _, c := range conns {
				_ = c.Close()
			}
			return nil, fmt.Errorf("ошибка подключения к %s: %v", addr, err)
		}
		if _, exists := conns[grpcServ.Name]; exists {
			return nil, fmt.Errorf("дубликат имени сервиса: %s", grpcServ.Name)
		}
		conns[grpcServ.Name] = conn
		log.Printf("Успешно подключен к gRPC сервису: %s", addr)
	}

	return conns, nil
}
