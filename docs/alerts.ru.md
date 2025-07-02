# Алерты для мониторинга токенов GitLab

Этот каталог содержит конфигурации алертов для мониторинга истечения токенов GitLab.

## Файлы алертов

### 1. `gitlab-token-alerts.yaml`
Основные алерты для мониторинга токенов:

- **TokenExpiresSoon** - срабатывает, когда до истечения токена осталось менее 2 недель (336 часов)
- **TokenExpiresCritical** - срабатывает, когда до истечения токена осталось менее недели (168 часов)
- **TokenExpired** - срабатывает, когда токен уже истек
- **TokenScraperErrors** - срабатывает при ошибках сбора метрик
- **TokenScraperDown** - срабатывает, когда экспортер недоступен

### 2. `gitlab-token-detailed-alerts.yaml`
Детальные алерты с различными временными интервалами:

- **TokenExpiresInOneMonth** - истечение через месяц (720 часов)
- **TokenExpiresInThreeDays** - истечение через 3 дня (72 часа)
- **TokenExpiresInOneDay** - истечение через день (24 часа)
- **TokenExpiresInOneHour** - истечение через час
- **NoTokenMetrics** - отсутствие метрик токенов

## Установка алертов

### Для Kubernetes с Prometheus Operator:

```bash
# Применить основные алерты
kubectl apply -f k8s/gitlab-token-alerts.yaml

# Применить детальные алерты (опционально)
kubectl apply -f k8s/gitlab-token-detailed-alerts.yaml
```

### Для standalone Prometheus:

Скопируйте содержимое файлов в секцию `rule_files` вашего `prometheus.yml`:

```yaml
rule_files:
  - "gitlab-token-alerts.yaml"
  - "gitlab-token-detailed-alerts.yaml"
```

## Настройка уведомлений

Для настройки уведомлений в Alertmanager добавьте в `alertmanager.yml`:

```yaml
route:
  group_by: ['alertname']
  group_wait: 10s
  group_interval: 10s
  repeat_interval: 1h
  receiver: 'web.hook'
  routes:
  - match:
      severity: critical
    receiver: 'pager-duty'
    repeat_interval: 30m
  - match:
      severity: warning
    receiver: 'slack'
    repeat_interval: 1h

receivers:
- name: 'web.hook'
  webhook_configs:
  - url: 'http://127.0.0.1:5001/'
- name: 'pager-duty'
  pagerduty_configs:
  - service_key: <your-pagerduty-key>
- name: 'slack'
  slack_configs:
  - api_url: 'https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK'
    channel: '#alerts'
```

## Метрики

Алерты основаны на следующих метриках:

- `gitlab_token_expires_at` - часы до истечения токена
- `gitlab_user_token_expires_at` - часы до истечения пользовательского токена
- `gitlab_token_is_expired` - флаг истечения токена (0/1)
- `gitlab_user_token_is_expired` - флаг истечения пользовательского токена (0/1)
- `gitlab_tokens_total` - общее количество токенов
- `gitlab_user_tokens_total` - общее количество пользовательских токенов
- `gitlab_token_scrape_errors_total` - количество ошибок сбора метрик
- `up{job="gitlab-token-exporter"}` - статус доступности экспортера

## Временные интервалы

- **336 часов** = 2 недели
- **168 часов** = 1 неделя
- **72 часа** = 3 дня
- **24 часа** = 1 день
- **1 час** = 1 час

## Уровни серьезности

- **info** - информационные уведомления
- **warning** - предупреждения, требующие внимания
- **critical** - критические проблемы, требующие немедленного вмешательства
