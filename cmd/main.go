package main

import (
	grpc_service "github.com/RVodassa/geo-microservices-auth/internal/grpc-server"
	"github.com/RVodassa/geo-microservices-auth/internal/service"
	proto "github.com/RVodassa/geo-microservices-auth/proto/generated"
	userService "github.com/RVodassa/geo-microservices-user/proto/generated"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"sync"
)

const userServiceAddress = "user-service:10101"

func main() {
	if err := godotenv.Load("/app/.env"); err != nil {
		log.Println("Ошибка загрузки .env файла")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	// Создание сервера
	grpcServer := grpc.NewServer()

	// Устанавливаем соединение с сервером
	conn, err := grpc.Dial(userServiceAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Не удалось подключиться к серверу: %v", err)
	}
	defer conn.Close()

	client := userService.NewUserServiceClient(conn)

	// Создание сервиса авторизации
	authService := service.NewAuthService(client, jwtSecret)

	// Регистрация сервиса в gRPC сервере
	authGrpcService := grpc_service.NewAuthServiceServer(authService)

	// Регистрируем сервис в gRPC сервере
	proto.RegisterAuthServiceServer(grpcServer, authGrpcService)

	// Запуск gRPC сервера
	var wg sync.WaitGroup
	errChan := make(chan error)

	wg.Add(1)
	go func() {
		defer wg.Done()

		listener, err := net.Listen("tcp", ":20202")
		if err != nil {
			errChan <- err
			return
		}

		if err := grpcServer.Serve(listener); err != nil {
			errChan <- err
			return
		}
	}()

	// Обработка ошибок
	select {
	case err := <-errChan:
		if err != nil {
			log.Fatalf("Ошибка при запуске gRPC сервера: %v", err)
		}
	}

	wg.Wait()
}
