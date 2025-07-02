package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Handler struct {
	tokenExpiresAt *prometheus.GaugeVec
	tokenIsExpired *prometheus.GaugeVec
	tokensTotal    prometheus.Gauge
	scrapeDuration prometheus.Histogram
	scrapeErrors   prometheus.Counter
	lastScrapeTime prometheus.Gauge
	// Метрики для пользовательских токенов
	userTokenExpiresAt *prometheus.GaugeVec
	userTokenIsExpired *prometheus.GaugeVec
	userTokensTotal    prometheus.Gauge
	// Метрики для групповых токенов
	groupTokenExpiresAt *prometheus.GaugeVec
	groupTokenIsExpired *prometheus.GaugeVec
	groupTokensTotal    prometheus.Gauge
}

func NewHandler() *Handler {
	h := &Handler{
		tokenExpiresAt: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "gitlab_token_expires_at",
				Help: "Hours until token expires",
			},
			[]string{"name"},
		),
		tokenIsExpired: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "gitlab_token_is_expired",
				Help: "Whether token is expired (1) or not (0)",
			},
			[]string{"name"},
		),
		tokensTotal: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "gitlab_tokens_total",
				Help: "Total number of tokens",
			},
		),
		scrapeDuration: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Name:    "gitlab_token_scrape_duration_seconds",
				Help:    "Duration of token scraping",
				Buckets: prometheus.DefBuckets,
			},
		),
		scrapeErrors: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "gitlab_token_scrape_errors_total",
				Help: "Total number of scraping errors",
			},
		),
		lastScrapeTime: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "gitlab_token_last_scrape_timestamp",
				Help: "Timestamp of last successful scrape",
			},
		),
		// Метрики для пользовательских токенов
		userTokenExpiresAt: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "gitlab_user_token_expires_at",
				Help: "Hours until user token expires",
			},
			[]string{"name"},
		),
		userTokenIsExpired: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "gitlab_user_token_is_expired",
				Help: "Whether user token is expired (1) or not (0)",
			},
			[]string{"name"},
		),
		userTokensTotal: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "gitlab_user_tokens_total",
				Help: "Total number of user tokens",
			},
		),
		// Метрики для групповых токенов
		groupTokenExpiresAt: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "gitlab_group_token_expires_at",
				Help: "Hours until group token expires",
			},
			[]string{"name"},
		),
		groupTokenIsExpired: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "gitlab_group_token_is_expired",
				Help: "Whether group token is expired (1) or not (0)",
			},
			[]string{"name"},
		),
		groupTokensTotal: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "gitlab_group_tokens_total",
				Help: "Total number of group tokens",
			},
		),
	}

	prometheus.MustRegister(
		h.tokenExpiresAt,
		h.tokenIsExpired,
		h.tokensTotal,
		h.scrapeDuration,
		h.scrapeErrors,
		h.lastScrapeTime,
		h.userTokenExpiresAt,
		h.userTokenIsExpired,
		h.userTokensTotal,
		h.groupTokenExpiresAt,
		h.groupTokenIsExpired,
		h.groupTokensTotal,
	)

	return h
}

func (h *Handler) Handler() http.Handler {
	return promhttp.Handler()
}

func (h *Handler) ResetMetrics() {
	h.tokenExpiresAt.Reset()
	h.tokenIsExpired.Reset()
	h.userTokenExpiresAt.Reset()
	h.userTokenIsExpired.Reset()
	h.groupTokenExpiresAt.Reset()
	h.groupTokenIsExpired.Reset()
}

func (h *Handler) SetTotalTokens(total int) {
	h.tokensTotal.Set(float64(total))
}

func (h *Handler) SetTokenExpiresAt(name string, expiresAt time.Time) {
	h.tokenExpiresAt.WithLabelValues(name).Set(time.Until(expiresAt).Hours())
}

func (h *Handler) SetTokenIsExpired(name string, isExpired bool) {
	value := 0.0
	if isExpired {
		value = 1.0
	}
	h.tokenIsExpired.WithLabelValues(name).Set(value)
}

func (h *Handler) RecordScrapeDuration(duration time.Duration) {
	h.scrapeDuration.Observe(duration.Seconds())
}

func (h *Handler) IncrementScrapeErrors() {
	h.scrapeErrors.Inc()
}

func (h *Handler) SetLastScrapeTime(timestamp time.Time) {
	h.lastScrapeTime.Set(float64(timestamp.Unix()))
}

func (h *Handler) DeleteTokenMetrics(name string) {
	h.tokenExpiresAt.DeleteLabelValues(name)
	h.tokenIsExpired.DeleteLabelValues(name)
}

// Методы для пользовательских токенов
func (h *Handler) SetTotalUserTokens(total int) {
	h.userTokensTotal.Set(float64(total))
}

func (h *Handler) SetUserTokenExpiresAt(name string, expiresAt time.Time) {
	h.userTokenExpiresAt.WithLabelValues(name).Set(time.Until(expiresAt).Hours())
}

func (h *Handler) SetUserTokenIsExpired(name string, isExpired bool) {
	value := 0.0
	if isExpired {
		value = 1.0
	}
	h.userTokenIsExpired.WithLabelValues(name).Set(value)
}

func (h *Handler) DeleteUserTokenMetrics(name string) {
	h.userTokenExpiresAt.DeleteLabelValues(name)
	h.userTokenIsExpired.DeleteLabelValues(name)
}

// Методы для групповых токенов
func (h *Handler) SetTotalGroupTokens(total int) {
	h.groupTokensTotal.Set(float64(total))
}

func (h *Handler) SetGroupTokenExpiresAt(name string, expiresAt time.Time) {
	h.groupTokenExpiresAt.WithLabelValues(name).Set(time.Until(expiresAt).Hours())
}

func (h *Handler) SetGroupTokenIsExpired(name string, isExpired bool) {
	value := 0.0
	if isExpired {
		value = 1.0
	}
	h.groupTokenIsExpired.WithLabelValues(name).Set(value)
}

func (h *Handler) DeleteGroupTokenMetrics(name string) {
	h.groupTokenExpiresAt.DeleteLabelValues(name)
	h.groupTokenIsExpired.DeleteLabelValues(name)
}
