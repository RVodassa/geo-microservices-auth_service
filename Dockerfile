# Этап сборки
FROM golang:1.23.3-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app/auth-service

# Устанавливаем необходимые пакеты
RUN apk add --no-cache gcc musl-dev

# Копируем файлы зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Сборка приложения
RUN go build -o auth-service main.go

# Финальный образ
FROM alpine:latest

# Устанавливаем рабочую директорию
WORKDIR /app/auth-service

# Копируем скомпилированное приложение
COPY --from=builder /app/auth-service/auth-service .

# Копируем конфиги и .env файл
COPY ./configs ./configs
COPY .env .

# Открываем порт
EXPOSE 20202

# Запуск приложения
CMD ["./auth-service"]