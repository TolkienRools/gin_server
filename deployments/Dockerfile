FROM golang:1.22.2

WORKDIR /app

# Копируем модули и загружаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код, сохраняя структуру проекта
COPY cmd ./cmd
COPY internal ./internal
COPY config ./config
COPY templates ./templates
# COPY uploads ./uploads

# Указываем путь к конфигурации
ENV CONFIG_PATH=/app/config/local.yaml

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -o ./weather-server ./cmd/weather_app/main.go

# Запускаем приложение
CMD ["./weather-server"]