# Этап сборки
FROM golang:1.23.3-alpine AS builder

WORKDIR /app/auth-service

RUN apk add --no-cache gcc musl-dev

# Копируем файлы зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Сборка приложения
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o authService-service main.go

# Финальный образ
FROM alpine

WORKDIR /app/auth-service

# Копируем скомпилированное приложение и конфиги
COPY --from=builder /app/auth-service/authService-service .
COPY ./configs /app/auth-service/configs
COPY .env /app/auth-service/

# порт
EXPOSE 20202

# Запуск приложения
CMD ["./authService-service"]