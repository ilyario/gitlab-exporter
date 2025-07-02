package scraper

import (
	"testing"
)

func TestTokenScraper_KnownTokensTracking(t *testing.T) {
	// Простой тест для проверки логики отслеживания известных токенов
	scraper := &TokenScraper{
		knownTokens:     make(map[string]bool),
		knownUserTokens: make(map[string]bool),
	}

	// Имитируем первый скрейпинг с двумя токенами проектов
	currentTokens := map[string]bool{
		"Project1 1 token1": true,
		"Project1 1 token2": true,
	}

	// Проверяем, что известные токены обновляются
	scraper.knownTokens = currentTokens
	if len(scraper.knownTokens) != 2 {
		t.Errorf("Expected 2 known project tokens, got %d", len(scraper.knownTokens))
	}

	// Имитируем второй скрейпинг с одним токеном (один удален)
	currentTokens = map[string]bool{
		"Project1 1 token1": true,
	}

	// Проверяем, что удаленный токен больше не в списке
	if currentTokens["Project1 1 token2"] {
		t.Error("Expected token2 to be removed from current project tokens")
	}

	// Имитируем скрейпинг пользовательских токенов
	currentUserTokens := map[string]bool{
		"User 123 token1": true,
		"User 456 token2": true,
	}

	scraper.knownUserTokens = currentUserTokens
	if len(scraper.knownUserTokens) != 2 {
		t.Errorf("Expected 2 known user tokens, got %d", len(scraper.knownUserTokens))
	}

	// Имитируем второй скрейпинг с одним пользовательским токеном (один удален)
	currentUserTokens = map[string]bool{
		"User 123 token1": true,
	}

	// Проверяем, что удаленный пользовательский токен больше не в списке
	if currentUserTokens["User 456 token2"] {
		t.Error("Expected user token2 to be removed from current user tokens")
	}
}
