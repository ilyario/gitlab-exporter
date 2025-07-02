# Переменные
BINARY_NAME=gitlab-token-exporter
BUILD_DIR=build
DOCKER_IMAGE=gitlab-token-exporter
DOCKER_TAG=latest

# Go переменные
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Цели
.PHONY: all build clean test deps docker-build docker-run help

# По умолчанию
all: clean deps test build

# Сборка приложения
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/server

# Очистка
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)

# Запуск тестов
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Запуск тестов с покрытием
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Установка зависимостей
deps:
	@echo "Installing dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Запуск приложения локально
run:
	@echo "Running application..."
	$(GOCMD) run ./cmd/server

# Сборка Docker образа
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

# Запуск Docker контейнера
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 \
		-e GITLAB_TOKEN=$$GITLAB_TOKEN \
		-e GITLAB_BASE_URL=$$GITLAB_BASE_URL \
		-e GITLAB_PROJECT_ID=$$GITLAB_PROJECT_ID \
		$(DOCKER_IMAGE):$(DOCKER_TAG)

# Запуск с Docker Compose
docker-compose-up:
	@echo "Starting with Docker Compose..."
	docker-compose up -d

# Остановка Docker Compose
docker-compose-down:
	@echo "Stopping Docker Compose..."
	docker-compose down

# Линтинг кода
lint:
	@echo "Running linter..."
	golangci-lint run

# Форматирование кода
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

# Проверка безопасности зависимостей
security-check:
	@echo "Checking dependencies for security vulnerabilities..."
	$(GOCMD) list -json -deps ./... | nancy sleuth

# Бенчмарки
bench:
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...

# Генерация документации
docs:
	@echo "Generating documentation..."
	godoc -http=:6060

# Установка инструментов разработки
install-tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/sonatype-nexus-community/nancy@latest

# Проверка кода перед коммитом
pre-commit: fmt lint test

# Помощь
help:
	@echo "Available targets:"
	@echo "  build              - Build the application"
	@echo "  clean              - Clean build artifacts"
	@echo "  test               - Run tests"
	@echo "  test-coverage      - Run tests with coverage report"
	@echo "  deps               - Install dependencies"
	@echo "  run                - Run application locally"
	@echo "  docker-build       - Build Docker image"
	@echo "  docker-run         - Run Docker container"
	@echo "  docker-compose-up  - Start with Docker Compose"
	@echo "  docker-compose-down- Stop Docker Compose"
	@echo "  lint               - Run linter"
	@echo "  fmt                - Format code"
	@echo "  security-check     - Check dependencies for vulnerabilities"
	@echo "  bench              - Run benchmarks"
	@echo "  docs               - Generate documentation"
	@echo "  install-tools      - Install development tools"
	@echo "  pre-commit         - Run pre-commit checks"
	@echo "  help               - Show this help"
