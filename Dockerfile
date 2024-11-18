FROM golang:1.23.3-alpine AS builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники приложения в рабочую директорию
COPY . .

# Компилируем приложение
RUN go build -o auth-service ./cmd/main.go

FROM alpine

WORKDIR /root/

# Копируем скомпилированное приложение из образа builder
COPY --from=builder /app/auth-service .

COPY .env .

# Открываем порт 20202
EXPOSE 20202

# Запускаем приложение
CMD ["./auth-service"]
