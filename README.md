# GitLab Exporter

A metrics exporter for monitoring GitLab project access tokens. Provides Prometheus metrics to track the status and expiration of access tokens.

## Features

- 🔍 Monitoring of GitLab project access tokens (supports multiple projects)
- 👤 Monitoring of GitLab user access tokens
- 👥 Monitoring of GitLab group access tokens
- 📊 Export of metrics in Prometheus format
- ⏰ Tracking token expiration time
- 🚨 Detection of expired tokens
- 🐳 Docker containerization
- 🔄 Graceful shutdown
- 🏥 Health checks

## Metrics

### Main Metrics

- `gitlab_token_expires_at` - Hours until project token expiration
- `gitlab_token_is_expired` - Project token expiration status (1 - expired, 0 - active)
- `gitlab_tokens_total` - Total number of project tokens

### User Token Metrics

- `gitlab_user_token_expires_at` - Hours until user token expiration
- `gitlab_user_token_is_expired` - User token expiration status (1 - expired, 0 - active)
- `gitlab_user_tokens_total` - Total number of user tokens

### Group Token Metrics

- `gitlab_group_token_expires_at` - Hours until group token expiration
- `gitlab_group_token_is_expired` - Group token expiration status (1 - expired, 0 - active)
- `gitlab_group_tokens_total` - Total number of group tokens

### Monitoring Metrics

- `gitlab_token_scrape_duration_seconds` - Scrape execution time
- `gitlab_token_scrape_errors_total` - Number of scrape errors
- `gitlab_token_last_scrape_timestamp` - Timestamp of the last successful scrape

## Quick Start

### Using Docker Compose

1. Copy the environment variables file:
```bash
cp env.example .env
```

2. Edit the `.env` file:
```bash
GITLAB_TOKEN=your_gitlab_token_here
GITLAB_BASE_URL=https://gitlab.com
GITLAB_PROJECT_IDS=12345,67890
```

3. Start the application:
```bash
docker-compose up -d
```

### Local Build

1. Install Go 1.23 or higher

2. Clone the repository:
```bash
git clone <repository-url>
cd gitlab-token-exporter
```

3. Install dependencies:
```bash
go mod download
```

4. Set environment variables:
```bash
export GITLAB_TOKEN=your_gitlab_token_here
export GITLAB_BASE_URL=https://gitlab.com
export GITLAB_PROJECT_IDS=12345,67890
```

5. Run the application:
```bash
go run ./cmd/server
```

## Configuration

### Environment Variables

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `GITLAB_TOKEN` | GitLab API token | Yes | - |
| `GITLAB_BASE_URL` | GitLab server URL | Yes | - |
| `GITLAB_PROJECT_IDS` | Comma-separated list of project IDs | Yes | - |
| `SERVER_PORT` | HTTP server port | No | 8080 |
| `SCRAPER_INTERVAL` | Metrics update interval | No | 10s |

### Endpoints

- `/metrics` - Prometheus metrics
- `/health` - Health check endpoint

## Monitoring

### Prometheus

The application automatically exports metrics in Prometheus format. Add the following to your Prometheus configuration:

```yaml
scrape_configs:
  - job_name: 'gitlab-token-exporter'
    static_configs:
      - targets: ['localhost:8080']
    scrape_interval: 10s
    metrics_path: /metrics
```

### Grafana

Import a ready-made dashboard or create your own using the following queries:

#### Number of active project tokens
```
gitlab_tokens_total
```

#### Number of active user tokens
```
gitlab_user_tokens_total
```

#### Number of active group tokens
```
gitlab_group_tokens_total
```

#### Project tokens expiring in less than 24 hours
```
gitlab_token_expires_at < 24
```

#### User tokens expiring in less than 24 hours
```
gitlab_user_token_expires_at < 24
```

#### Group tokens expiring in less than 24 hours
```
gitlab_group_token_expires_at < 24
```

#### Expired project tokens
```
gitlab_token_is_expired == 1
```

#### Expired user tokens
```
gitlab_user_token_is_expired == 1
```

#### Expired group tokens
```
gitlab_group_token_is_expired == 1
```

## Development

### Project Structure

```
gitlab-token-exporter/
├── cmd/
│   └── server/          # Application entry point
├── internal/
│   ├── config/          # Configuration
│   ├── gitlab/          # GitLab client
│   ├── metrics/         # Metrics handling
│   └── scraper/         # Data scraping logic
├── configs/             # Configuration files
├── Dockerfile           # Docker image
├── docker-compose.yml   # Docker Compose
└── README.md            # Documentation
```

### Adding New Metrics

1. Add the new metric to `internal/metrics/handler.go`
2. Register the metric in the constructor
3. Add methods to set values
4. Use the metric in the scraper

### Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run benchmarks
go test -bench=. ./...
```

## Logging

The application uses the standard Go logger. For production, it is recommended to configure structured logging.

## Security

- The application runs as a non-privileged user in Docker
- All secrets are passed via environment variables
- Health checks for monitoring status

## Troubleshooting

### Issues connecting to GitLab

1. Check the correctness of `GITLAB_BASE_URL`
2. Make sure the token has the necessary permissions
3. Check the availability of the GitLab server

### Issues with metrics

1. Check the `/metrics` endpoint
2. Make sure Prometheus can connect to the application
3. Check the logs for errors

## License

MIT License
