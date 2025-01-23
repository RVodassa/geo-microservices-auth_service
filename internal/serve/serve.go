package serve

import (
	"github.com/RVodassa/geo-microservices-auth_service/internal/handlers/grpcHandlers"
	"github.com/RVodassa/geo-microservices-auth_service/internal/service"
	proto "github.com/RVodassa/geo-microservices-auth_service/proto/generated"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
	"time"
)

type Serve struct {
}

// RunServe Запуск gRPC сервера
func RunServe(authService *service.AuthService, port string) error {
	authGrpcService := grpcHandlers.NewAuthServiceServer(authService)
	grpcServer := grpc.NewServer()
	proto.RegisterAuthServiceServer(grpcServer, authGrpcService)

	var wg sync.WaitGroup
	errChan := make(chan error)

	wg.Add(1)
	go func() {
		defer wg.Done()

		listener, err := net.Listen("tcp", port)
		if err != nil {
			errChan <- err
			return
		}

		if err = grpcServer.Serve(listener); err != nil {
			errChan <- err
			return
		}
	}()

	// Обработка ошибок
	select {
	case err := <-errChan:
		log.Fatalf("Ошибка при запуске gRPC сервера: %v", err)
		return err
	case <-time.After(1 * time.Second):
		log.Println("gRPC сервер успешно запущен на порту :20202")
	}

	wg.Wait()
	return nil
}
