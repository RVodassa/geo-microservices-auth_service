package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/RVodassa/geo-microservices-auth_service/internal/domain/entity"
	userService "github.com/RVodassa/geo-microservices-user_service/proto/generated"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const expTime = time.Minute // Токен будет действовать 1 час

type AuthServiceProvider interface {
	Register(ctx context.Context, user *entity.User) (id uint64, err error)
	Login(ctx context.Context, login, password string) (string, error)
	CheckToken(ctx context.Context, token string) (bool, error)
}

type AuthService struct {
	userService userService.UserServiceClient
	secret      string
}

func NewAuthService(userService userService.UserServiceClient, secret string) *AuthService {
	return &AuthService{
		userService: userService,
		secret:      secret,
	}
}

// Проверка токена
func (s *AuthService) CheckToken(ctx context.Context, tokenString string) (bool, error) {
	// Секретный ключ (должен быть таким же, как при создании токена)
	secretKey := []byte(s.secret)

	// Парсим токен
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверка алгоритма токена (например, HMAC)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неверный метод подписания: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return false, fmt.Errorf("неверный или просроченный токен: %v", err)
	}

	// Извлекаем claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return false, fmt.Errorf("неверный токен")
	}

	// Проверка срока действия токена
	if exp, ok := claims["exp"].(float64); ok {
		expirationTime := time.Unix(int64(exp), 0)
		if expirationTime.Before(time.Now()) {
			return false, fmt.Errorf("токен истек в %v", expirationTime)
		}
	} else {
		return false, fmt.Errorf("отсутствует поле 'exp' в claims токена")
	}

	// Проверка наличия поля login
	if login, ok := claims["login"].(string); ok {
		fmt.Printf("Токен валиден, данные: login = %s\n", login)
	} else {
		return false, fmt.Errorf("отсутствует требуемое поле login в claims")
	}

	return true, nil
}

func (s *AuthService) Login(ctx context.Context, login, password string) (string, error) {
	// Получаем пользователя по логину
	req := &userService.LoginRequest{Login: login, Password: password}
	response, err := s.userService.Login(ctx, req)

	if err != nil {
		return "", fmt.Errorf("ошибка при проверке пользователя: %v", err)
	}

	if !response.Status {
		return "", errors.New("неверные данные авторизации")
	}

	// Генерируем JWT токен со сроком действия
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login": login,
		"exp":   time.Now().Add(expTime).Unix(),
	})

	// Подписываем токен
	tokenString, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", fmt.Errorf("ошибка при генерации токена: %v", err)
	}

	return tokenString, nil
}

func (s *AuthService) Register(ctx context.Context, user *entity.User) (id uint64, err error) {
	// Создание запроса для регистрации пользователя
	registerReq := &userService.RegisterRequest{
		Login:    user.Login,
		Password: user.Password, // Передаем пароль в открытом виде
	}

	// Выполняем вызов метода Register
	res, err := s.userService.Register(ctx, registerReq)
	if err != nil {
		return 0, fmt.Errorf("ошибка при регистрации пользователя: %v", err)
	}

	fmt.Printf("Пользователь зарегистрирован с ID: %d\n", res.Id)
	return res.Id, nil
}
