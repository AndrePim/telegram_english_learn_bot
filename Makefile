# Makefile для Telegram English Bot

.PHONY: help build run test lint clean docker-build docker-run docker-stop dev-setup

# Переменные
APP_NAME=tg-english-bot
DOCKER_IMAGE=your-username/$(APP_NAME)
DOCKER_TAG=latest

# По умолчанию показываем справку
help: ## Показать справку
	@echo "Доступные команды:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Разработка
build: ## Собрать приложение
	go build -o $(APP_NAME) ./cmd/app

run: ## Запустить приложение
	go run ./cmd/app

test: ## Запустить тесты
	go test -v ./...

test-cover: ## Запустить тесты с покрытием
	go test -cover ./...

lint: ## Запустить линтер
	golangci-lint run

clean: ## Очистить сборочные файлы
	rm -f $(APP_NAME)
	go clean

# Docker команды
docker-build: ## Собрать Docker образ
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-run: ## Запустить через Docker Compose (development)
	docker-compose up -d

docker-stop: ## Остановить Docker Compose
	docker-compose down

docker-logs: ## Показать логи Docker контейнеров
	docker-compose logs -f

docker-clean: ## Очистить Docker ресурсы
	docker-compose down -v
	docker system prune -f

# Production команды
prod-deploy: ## Деплоймент в production
	docker-compose -f docker-compose.prod.yml --env-file .env.prod up -d

prod-stop: ## Остановить production
	docker-compose -f docker-compose.prod.yml down

prod-logs: ## Логи production
	docker-compose -f docker-compose.prod.yml logs -f

prod-backup: ## Создать бэкап базы данных
	docker-compose -f docker-compose.prod.yml exec db pg_dump -U $$(grep DB_USER .env.prod | cut -d '=' -f2) $$(grep DB_NAME .env.prod | cut -d '=' -f2) > backup_$$(date +%Y%m%d_%H%M%S).sql

# Настройка среды разработки
dev-setup: ## Настроить среду разработки
	@echo "Настройка среды разработки..."
	@if [ ! -f .env ]; then cp .env.example .env; echo "Создан .env файл. Отредактируйте его перед запуском."; fi
	go mod download
	@echo "Среда разработки готова!"

# Форматирование кода
fmt: ## Форматировать код
	go fmt ./...
	goimports -w .

# Обновление зависимостей
deps-update: ## Обновить зависимости
	go get -u ./...
	go mod tidy

# Генерация документации
docs: ## Генерировать документацию
	godoc -http=:6060
	@echo "Документация доступна по адресу: http://localhost:6060"

