# GitLab Exporter

Экспортер метрик для мониторинга токенов доступа GitLab проектов. Предоставляет метрики Prometheus для отслеживания состояния и срока действия токенов доступа.

## Возможности

- 🔍 Мониторинг токенов доступа GitLab проектов (поддержка нескольких проектов)
- 👤 Мониторинг пользовательских токенов доступа GitLab
- 👥 Мониторинг групповых токенов доступа GitLab
- 📊 Экспорт метрик в формате Prometheus
- ⏰ Отслеживание времени истечения токенов
- 🚨 Обнаружение просроченных токенов
- 🐳 Docker контейнеризация
- 🔄 Graceful shutdown
- 🏥 Health checks

## Метрики

### Основные метрики

- `gitlab_token_expires_at` - Часы до истечения токена проекта
- `gitlab_token_is_expired` - Статус истечения токена проекта (1 - истек, 0 - активен)
- `gitlab_tokens_total` - Общее количество токенов проектов

### Метрики пользовательских токенов

- `gitlab_user_token_expires_at` - Часы до истечения пользовательского токена
- `gitlab_user_token_is_expired` - Статус истечения пользовательского токена (1 - истек, 0 - активен)
- `gitlab_user_tokens_total` - Общее количество пользовательских токенов

### Метрики групповых токенов

- `gitlab_group_token_expires_at` - Часы до истечения группового токена
- `gitlab_group_token_is_expired` - Статус истечения группового токена (1 - истек, 0 - активен)
- `gitlab_group_tokens_total` - Общее количество групповых токенов

### Метрики мониторинга

- `gitlab_token_scrape_duration_seconds` - Время выполнения scrape
- `gitlab_token_scrape_errors_total` - Количество ошибок scrape
- `gitlab_token_last_scrape_timestamp` - Время последнего успешного scrape

## Быстрый старт

### С использованием Docker Compose

1. Скопируйте файл с переменными окружения:
```bash
cp env.example .env
```

2. Отредактируйте `.env` файл:
```bash
GITLAB_TOKEN=your_gitlab_token_here
GITLAB_BASE_URL=https://gitlab.com
GITLAB_PROJECT_IDS=12345,67890
```

3. Запустите приложение:
```bash
docker-compose up -d
```

### Локальная сборка

1. Установите Go 1.23 или выше

2. Клонируйте репозиторий:
```bash
git clone <repository-url>
cd gitlab-token-exporter
```

3. Установите зависимости:
```bash
go mod download
```

4. Установите переменные окружения:
```bash
export GITLAB_TOKEN=your_gitlab_token_here
export GITLAB_BASE_URL=https://gitlab.com
export GITLAB_PROJECT_IDS=12345,67890
```

5. Запустите приложение:
```bash
go run ./cmd/server
```

## Конфигурация

### Переменные окружения

| Переменная | Описание | Обязательная | По умолчанию |
|------------|----------|--------------|--------------|
| `GITLAB_TOKEN` | GitLab API токен | Да | - |
| `GITLAB_BASE_URL` | URL GitLab сервера | Да | - |
| `GITLAB_PROJECT_IDS` | Список ID проектов через запятую | Да | - |
| `SERVER_PORT` | Порт HTTP сервера | Нет | 8080 |
| `SCRAPER_INTERVAL` | Интервал обновления метрик | Нет | 10s |

### Endpoints

- `/metrics` - Метрики Prometheus
- `/health` - Health check endpoint

## Мониторинг

### Prometheus

Приложение автоматически экспортирует метрики в формате Prometheus. Добавьте в конфигурацию Prometheus:

```yaml
scrape_configs:
  - job_name: 'gitlab-token-exporter'
    static_configs:
      - targets: ['localhost:8080']
    scrape_interval: 10s
    metrics_path: /metrics
```

### Grafana

Импортируйте готовый дашборд или создайте свой с использованием следующих запросов:

#### Количество активных токенов проектов
```
gitlab_tokens_total
```

#### Количество активных пользовательских токенов
```
gitlab_user_tokens_total
```

#### Количество активных групповых токенов
```
gitlab_group_tokens_total
```

#### Токены проектов с истекающим сроком (менее 24 часов)
```
gitlab_token_expires_at < 24
```

#### Пользовательские токены с истекающим сроком (менее 24 часов)
```
gitlab_user_token_expires_at < 24
```

#### Групповые токены с истекающим сроком (менее 24 часов)
```
gitlab_group_token_expires_at < 24
```

#### Просроченные токены проектов
```
gitlab_token_is_expired == 1
```

#### Просроченные пользовательские токены
```
gitlab_user_token_is_expired == 1
```

#### Просроченные групповые токены
```
gitlab_group_token_is_expired == 1
```

## Разработка

### Структура проекта

```
gitlab-token-exporter/
├── cmd/
│   └── server/          # Точка входа приложения
├── internal/
│   ├── config/          # Конфигурация
│   ├── gitlab/          # GitLab клиент
│   ├── metrics/         # Обработка метрик
│   └── scraper/         # Логика сбора данных
├── configs/             # Конфигурационные файлы
├── Dockerfile           # Docker образ
├── docker-compose.yml   # Docker Compose
└── README.md           # Документация
```

### Добавление новых метрик

1. Добавьте новую метрику в `internal/metrics/handler.go`
2. Зарегистрируйте метрику в конструкторе
3. Добавьте методы для установки значений
4. Используйте метрику в scraper

### Тестирование

```bash
# Запуск тестов
go test ./...

# Запуск тестов с покрытием
go test -cover ./...

# Запуск бенчмарков
go test -bench=. ./...
```

## Логирование

Приложение использует стандартный Go logger. Для продакшена рекомендуется настроить структурированное логирование.

## Безопасность

- Приложение запускается под непривилегированным пользователем в Docker
- Все секреты передаются через переменные окружения
- Health checks для мониторинга состояния

## Troubleshooting

### Проблемы с подключением к GitLab

1. Проверьте правильность `GITLAB_BASE_URL`
2. Убедитесь, что токен имеет необходимые права доступа
3. Проверьте доступность GitLab сервера

### Проблемы с метриками

1. Проверьте endpoint `/metrics`
2. Убедитесь, что Prometheus может подключиться к приложению
3. Проверьте логи на наличие ошибок

## Лицензия

MIT License
