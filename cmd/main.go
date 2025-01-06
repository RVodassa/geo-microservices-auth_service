package main

import (
	authGrpcServer "github.com/RVodassa/geo-microservices-auth_service/internal/grpc-server"
	authService "github.com/RVodassa/geo-microservices-auth_service/internal/service"
	proto "github.com/RVodassa/geo-microservices-auth_service/proto/generated"
	userService "github.com/RVodassa/geo-microservices-user_service/proto/generated"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"os"
	"sync"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Ошибка загрузки .env файла")
	}

	jwtSecret := os.Getenv("JWT_SECRET")

	address := "user-service:10101"
	// Создание сервера
	grpcServer := grpc.NewServer()

	// Устанавливаем соединение с сервером
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Не удалось подключиться к серверу: %v", err)
	}
	defer conn.Close()

	client := userService.NewUserServiceClient(conn)

	// Создание сервиса авторизации
	authService := authService.NewAuthService(client, jwtSecret)

	// Регистрация сервиса в gRPC сервере
	authGrpcService := authGrpcServer.NewAuthServiceServer(authService)

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
			return
		}
	}

	wg.Wait()
}
