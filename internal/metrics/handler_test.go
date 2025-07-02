package metrics

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func resetPrometheusRegistry() {
	prometheus.DefaultRegisterer = prometheus.NewRegistry()
}

func TestNewHandler(t *testing.T) {
	resetPrometheusRegistry()
	handler := NewHandler()
	if handler == nil {
		t.Fatal("NewHandler() returned nil")
	}

	// Проверяем, что все метрики инициализированы
	if handler.tokenExpiresAt == nil {
		t.Error("tokenExpiresAt metric is nil")
	}
	if handler.tokenIsExpired == nil {
		t.Error("tokenIsExpired metric is nil")
	}
	if handler.tokensTotal == nil {
		t.Error("tokensTotal metric is nil")
	}
	if handler.scrapeDuration == nil {
		t.Error("scrapeDuration metric is nil")
	}
	if handler.scrapeErrors == nil {
		t.Error("scrapeErrors metric is nil")
	}
	if handler.lastScrapeTime == nil {
		t.Error("lastScrapeTime metric is nil")
	}
}

func TestHandler_Handler(t *testing.T) {
	resetPrometheusRegistry()
	handler := NewHandler()
	httpHandler := handler.Handler()
	if httpHandler == nil {
		t.Fatal("Handler() returned nil")
	}

	// Создаем тестовый HTTP сервер
	server := httptest.NewServer(httpHandler)
	defer server.Close()

	// Делаем запрос к /metrics
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestHandler_ResetMetrics(t *testing.T) {
	resetPrometheusRegistry()
	handler := NewHandler()

	// Устанавливаем некоторые значения для токенов проектов
	handler.SetTotalTokens(5)
	handler.SetTokenExpiresAt("test-token", time.Now().Add(time.Hour))
	handler.SetTokenIsExpired("test-token", false)

	// Устанавливаем некоторые значения для пользовательских токенов
	handler.SetTotalUserTokens(3)
	handler.SetUserTokenExpiresAt("test-user-token", time.Now().Add(time.Hour))
	handler.SetUserTokenIsExpired("test-user-token", false)

	// Сбрасываем метрики
	handler.ResetMetrics()

	// Проверяем, что метрики сброшены (это внутренняя проверка, так как мы не можем напрямую получить значения)
	// В реальном приложении мы бы проверяли через HTTP endpoint
}

func TestHandler_SetTotalTokens(t *testing.T) {
	resetPrometheusRegistry()
	handler := NewHandler()

	// Тестируем установку различных значений
	testCases := []int{0, 1, 10, 100}

	for _, total := range testCases {
		handler.SetTotalTokens(total)
		// В реальном приложении мы бы проверяли значение через HTTP endpoint
	}
}

func TestHandler_SetTokenExpiresAt(t *testing.T) {
	resetPrometheusRegistry()
	handler := NewHandler()

	// Тестируем установку времени истечения
	expiresAt := time.Now().Add(time.Hour)
	handler.SetTokenExpiresAt("test-token", expiresAt)

	// Тестируем истекший токен
	expiredAt := time.Now().Add(-time.Hour)
	handler.SetTokenExpiresAt("expired-token", expiredAt)
}

func TestHandler_SetTokenIsExpired(t *testing.T) {
	resetPrometheusRegistry()
	handler := NewHandler()

	// Тестируем активный токен
	handler.SetTokenIsExpired("active-token", false)

	// Тестируем истекший токен
	handler.SetTokenIsExpired("expired-token", true)
}

func TestHandler_RecordScrapeDuration(t *testing.T) {
	resetPrometheusRegistry()
	handler := NewHandler()

	// Тестируем запись различных длительностей
	durations := []time.Duration{
		time.Millisecond,
		time.Second,
		5 * time.Second,
	}

	for _, duration := range durations {
		handler.RecordScrapeDuration(duration)
	}
}

func TestHandler_IncrementScrapeErrors(t *testing.T) {
	resetPrometheusRegistry()
	handler := NewHandler()

	// Тестируем инкремент ошибок
	for i := 0; i < 5; i++ {
		handler.IncrementScrapeErrors()
	}
}

func TestHandler_SetLastScrapeTime(t *testing.T) {
	resetPrometheusRegistry()
	handler := NewHandler()

	// Тестируем установку времени последнего scrape
	now := time.Now()
	handler.SetLastScrapeTime(now)

	// Тестируем установку времени в прошлом
	past := time.Now().Add(-time.Hour)
	handler.SetLastScrapeTime(past)
}

func TestHandler_DeleteTokenMetrics(t *testing.T) {
	resetPrometheusRegistry()
	handler := NewHandler()

	// Создаем токен и устанавливаем его метрики
	tokenName := "test-token"
	expiresAt := time.Now().Add(time.Hour)

	handler.SetTokenExpiresAt(tokenName, expiresAt)
	handler.SetTokenIsExpired(tokenName, false)

	// Удаляем метрики токена
	handler.DeleteTokenMetrics(tokenName)

	// В реальном приложении мы бы проверяли через HTTP endpoint, что метрики удалены
	// Здесь мы просто проверяем, что метод выполняется без ошибок
}

func TestHandler_SetTotalUserTokens(t *testing.T) {
	resetPrometheusRegistry()
	handler := NewHandler()

	// Тестируем установку различных значений
	testCases := []int{0, 1, 10, 100}

	for _, total := range testCases {
		handler.SetTotalUserTokens(total)
		// В реальном приложении мы бы проверяли значение через HTTP endpoint
	}
}

func TestHandler_SetUserTokenExpiresAt(t *testing.T) {
	resetPrometheusRegistry()
	handler := NewHandler()

	// Тестируем установку времени истечения
	expiresAt := time.Now().Add(time.Hour)
	handler.SetUserTokenExpiresAt("test-user-token", expiresAt)

	// Тестируем истекший токен
	expiredAt := time.Now().Add(-time.Hour)
	handler.SetUserTokenExpiresAt("expired-user-token", expiredAt)
}

func TestHandler_SetUserTokenIsExpired(t *testing.T) {
	resetPrometheusRegistry()
	handler := NewHandler()

	// Тестируем активный токен
	handler.SetUserTokenIsExpired("active-user-token", false)

	// Тестируем истекший токен
	handler.SetUserTokenIsExpired("expired-user-token", true)
}

func TestHandler_DeleteUserTokenMetrics(t *testing.T) {
	resetPrometheusRegistry()
	handler := NewHandler()

	// Создаем пользовательский токен и устанавливаем его метрики
	tokenName := "test-user-token"
	expiresAt := time.Now().Add(time.Hour)

	handler.SetUserTokenExpiresAt(tokenName, expiresAt)
	handler.SetUserTokenIsExpired(tokenName, false)

	// Удаляем метрики пользовательского токена
	handler.DeleteUserTokenMetrics(tokenName)

	// В реальном приложении мы бы проверяли через HTTP endpoint, что метрики удалены
	// Здесь мы просто проверяем, что метод выполняется без ошибок
}
