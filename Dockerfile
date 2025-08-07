# Этап сборки
FROM golang:1.24-alpine AS builder

# Устанавливаем необходимые пакеты
RUN apk add --no-cache git

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/app

# Финальный этап
FROM alpine:latest

# Устанавливаем ca-certificates для HTTPS запросов
RUN apk --no-cache add ca-certificates

# Создаем пользователя для запуска приложения
RUN adduser -D -s /bin/sh appuser

WORKDIR /root/

# Копируем бинарный файл из этапа сборки
COPY --from=builder /app/main .

# Меняем владельца файла
RUN chown appuser:appuser main

# Переключаемся на пользователя appuser
USER appuser

# Команда для запуска приложения
CMD ["./main"]

