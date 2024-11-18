package authGrpcServer

import (
	"context"
	"github.com/RVodassa/geo-microservices-auth_service/internal/domain/entity"
	authService "github.com/RVodassa/geo-microservices-auth_service/internal/service"
	pb "github.com/RVodassa/geo-microservices-auth_service/proto/generated"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServiceServer struct {
	pb.UnimplementedAuthServiceServer
	authService authService.AuthServiceProvider
}

func NewAuthServiceServer(authService authService.AuthServiceProvider) *AuthServiceServer {
	return &AuthServiceServer{authService: authService}
}

func (s *AuthServiceServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// Создание сущности пользователя без хеширования пароля
	user := &entity.User{
		Login:    req.Login,
		Password: req.Password, // Пароль в открытом виде
	}

	// Регистрация пользователя через сервис
	id, err := s.authService.Register(ctx, user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &pb.RegisterResponse{
		Id: id,
	}

	return resp, nil
}

func (s *AuthServiceServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {

	token, err := s.authService.Login(ctx, req.Login, req.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	resp := &pb.LoginResponse{
		Token: token,
	}
	return resp, nil
}

func (s *AuthServiceServer) CheckToken(ctx context.Context, req *pb.CheckTokenRequest) (*pb.CheckTokenResponse, error) {
	verify, err := s.authService.CheckToken(ctx, req.Token)
	if err != nil {
		return nil, err
	}
	resp := &pb.CheckTokenResponse{
		Status: verify,
	}
	return resp, nil
}
