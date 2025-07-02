# Alerts for monitoring GitLab tokens

This directory contains alert configurations for monitoring GitLab token expiration.

## Alert files

### 1. `gitlab-token-alerts.yaml`
Basic alerts for token monitoring:

- **TokenExpiresSoon** - triggered when the token expires in less than 2 weeks (336 hours)
- **TokenExpiresCritical** - triggered when the token expires in less than a week (168 hours)
- **TokenExpired** - triggered when the token has already expired
- **TokenScraperErrors** - triggered when there are errors collecting metrics
- **TokenScraperDown** - triggered when the exporter is unavailable

### 2. `gitlab-token-detailed-alerts.yaml`
Detailed alerts with different time intervals:

- **TokenExpiresInOneMonth** - expiration in one month (720 hours)
- **TokenExpiresInThreeDays** - expiration in 3 days (72 hours)
- **TokenExpiresInOneDay** - expiration in one day (24 hours)
- **TokenExpiresInOneHour** - expiration in one hour
- **NoTokenMetrics** - no token metrics

## Setting alerts

### For Kubernetes with Prometheus Operator:

```bash
# Apply basic alerts
kubectl apply -f k8s/gitlab-token-alerts.yaml

# Apply detailed alerts (optional)
kubectl apply -f k8s/gitlab-token-detailed-alerts.yaml
```

### For standalone Prometheus:

Copy the contents of the files to the `rule_files` section of your `prometheus.yml`:

```yaml
rule_files:
- "gitlab-token-alerts.yaml"
- "gitlab-token-detailed-alerts.yaml"
```

## Setting up notifications

To configure notifications in Alertmanager, add to `alertmanager.yml`:

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

## Metrics

Alerts are based on the following metrics:

- `gitlab_token_expires_at` - hours until token expiration
- `gitlab_user_token_expires_at` - hours until user token expiration
- `gitlab_token_is_expired` - token expiration flag (0/1)
- `gitlab_user_token_is_expired` - user token expiration flag (0/1)
- `gitlab_tokens_total` - total number of tokens
- `gitlab_user_tokens_total` - total number of user tokens
- `gitlab_token_scrape_errors_total` - number of metrics scraping errors
- `up{job="gitlab-token-exporter"}` - exporter availability status

## Time intervals

- **336 hours** = 2 weeks
- **168 hours** = 1 week
- **72 hours** = 3 days
- **24 hours** = 1 day
- **1 hour** = 1 hour

## Severity levels

- **info** - informational notifications
- **warning** - warnings that require attention
- **critical** - critical issues that require immediate attention
