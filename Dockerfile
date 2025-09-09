# Этап сборки приложения
FROM golang:1.24.5 AS builder

# Создание рабочей директории
WORKDIR /app

# Копирование go-модуля и сумм
COPY go.mod go.sum ./

# Загрузка зависимостей
RUN go mod download

# Копирование всего исходного кода
COPY . .

# Сборка бинарного файла
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/main.go

# Финальный этап развертывания
FROM alpine:latest

# Установка сертификатов
RUN apk add --update ca-certificates tzdata

# Копирование бинарного файла и ресурсов
COPY --from=builder /app/app /app/app
COPY --from=builder /app/config.yaml /app/config.yaml
COPY --from=builder /app/migrations /app/migrations

# Установка рабочего каталога
WORKDIR /app

# Экспозиция порта
EXPOSE 8081

# Запуск приложения
CMD ["./app"]