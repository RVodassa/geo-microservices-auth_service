package app

import (
	"fmt"
	"github.com/RVodassa/geo-microservices-auth_service/app/config"
	"github.com/RVodassa/geo-microservices-auth_service/internal/domain/logger"
	"github.com/RVodassa/geo-microservices-auth_service/internal/serve"
	serviceAuth "github.com/RVodassa/geo-microservices-auth_service/internal/service"
	serviceUser "github.com/RVodassa/geo-microservices-user_service/proto/generated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
)

const UserServiceName = "user-service"
const SecretJWT = "JWT_SECRET"

type App struct {
	log logger.Logger
	cfg *config.Config
}

func NewApp(log logger.Logger, cfg *config.Config) *App {
	return &App{
		log: log,
		cfg: cfg,
	}
}

func (app *App) Run() error {

	log.Println("Установка соединений с gRPC сервисами...")

	grpcConns, err := setupGRPCClient(app.cfg.GRPCServices)
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

	secretJWT := os.Getenv(SecretJWT)
	authService := serviceAuth.NewAuthService(client, secretJWT)
	log.Println("Сервис авторизации успешно создан")

	if err = serve.RunServe(authService, app.cfg.ServePort); err != nil {
		return fmt.Errorf("ошибка при запуске сервера: %v", err)
	}

	return nil
}

func setupGRPCClient(grpcServices []config.GRPCService) (map[string]*grpc.ClientConn, error) {

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
