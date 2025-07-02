package config

import (
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestLoad(t *testing.T) {
	// Сохраняем оригинальные переменные окружения
	originalToken := os.Getenv("GITLAB_TOKEN")
	originalBaseURL := os.Getenv("GITLAB_BASE_URL")
	originalProjectIDs := os.Getenv("GITLAB_PROJECT_IDS")
	originalPort := os.Getenv("SERVER_PORT")
	originalInterval := os.Getenv("SCRAPER_INTERVAL")

	// Восстанавливаем переменные после теста
	defer func() {
		if originalToken != "" {
			os.Setenv("GITLAB_TOKEN", originalToken)
		} else {
			os.Unsetenv("GITLAB_TOKEN")
		}
		if originalBaseURL != "" {
			os.Setenv("GITLAB_BASE_URL", originalBaseURL)
		} else {
			os.Unsetenv("GITLAB_BASE_URL")
		}
		if originalProjectIDs != "" {
			os.Setenv("GITLAB_PROJECT_IDS", originalProjectIDs)
		} else {
			os.Unsetenv("GITLAB_PROJECT_IDS")
		}
		if originalPort != "" {
			os.Setenv("SERVER_PORT", originalPort)
		} else {
			os.Unsetenv("SERVER_PORT")
		}
		if originalInterval != "" {
			os.Setenv("SCRAPER_INTERVAL", originalInterval)
		} else {
			os.Unsetenv("SCRAPER_INTERVAL")
		}
	}()

	tests := []struct {
		name    string
		env     map[string]string
		wantErr bool
	}{
		{
			name: "valid configuration single project",
			env: map[string]string{
				"GITLAB_TOKEN":       "test-token",
				"GITLAB_BASE_URL":    "https://gitlab.com",
				"GITLAB_PROJECT_IDS": "12345",
				"SERVER_PORT":        "9090",
				"SCRAPER_INTERVAL":   "30s",
			},
			wantErr: false,
		},
		{
			name: "valid configuration multiple projects",
			env: map[string]string{
				"GITLAB_TOKEN":       "test-token",
				"GITLAB_BASE_URL":    "https://gitlab.com",
				"GITLAB_PROJECT_IDS": "12345,67890,13579",
				"SERVER_PORT":        "9090",
				"SCRAPER_INTERVAL":   "30s",
			},
			wantErr: false,
		},
		{
			name: "valid configuration with spaces",
			env: map[string]string{
				"GITLAB_TOKEN":       "test-token",
				"GITLAB_BASE_URL":    "https://gitlab.com",
				"GITLAB_PROJECT_IDS": "12345, 67890, 13579",
				"SERVER_PORT":        "9090",
				"SCRAPER_INTERVAL":   "30s",
			},
			wantErr: false,
		},
		{
			name: "missing gitlab token",
			env: map[string]string{
				"GITLAB_BASE_URL":    "https://gitlab.com",
				"GITLAB_PROJECT_IDS": "12345",
			},
			wantErr: true,
		},
		{
			name: "missing gitlab base url",
			env: map[string]string{
				"GITLAB_TOKEN":       "test-token",
				"GITLAB_PROJECT_IDS": "12345",
			},
			wantErr: true,
		},
		{
			name: "missing project ids",
			env: map[string]string{
				"GITLAB_TOKEN":    "test-token",
				"GITLAB_BASE_URL": "https://gitlab.com",
			},
			wantErr: true,
		},
		{
			name: "empty project ids",
			env: map[string]string{
				"GITLAB_TOKEN":       "test-token",
				"GITLAB_BASE_URL":    "https://gitlab.com",
				"GITLAB_PROJECT_IDS": "",
			},
			wantErr: true,
		},
		{
			name: "invalid project id",
			env: map[string]string{
				"GITLAB_TOKEN":       "test-token",
				"GITLAB_BASE_URL":    "https://gitlab.com",
				"GITLAB_PROJECT_IDS": "invalid",
			},
			wantErr: true,
		},
		{
			name: "mixed valid and invalid project ids",
			env: map[string]string{
				"GITLAB_TOKEN":       "test-token",
				"GITLAB_BASE_URL":    "https://gitlab.com",
				"GITLAB_PROJECT_IDS": "12345,invalid,67890",
			},
			wantErr: true,
		},
		{
			name: "invalid server port",
			env: map[string]string{
				"GITLAB_TOKEN":       "test-token",
				"GITLAB_BASE_URL":    "https://gitlab.com",
				"GITLAB_PROJECT_IDS": "12345",
				"SERVER_PORT":        "invalid",
			},
			wantErr: true,
		},
		{
			name: "invalid scraper interval",
			env: map[string]string{
				"GITLAB_TOKEN":       "test-token",
				"GITLAB_BASE_URL":    "https://gitlab.com",
				"GITLAB_PROJECT_IDS": "12345",
				"SCRAPER_INTERVAL":   "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Очищаем все переменные окружения перед тестом
			os.Unsetenv("GITLAB_TOKEN")
			os.Unsetenv("GITLAB_BASE_URL")
			os.Unsetenv("GITLAB_PROJECT_IDS")
			os.Unsetenv("SERVER_PORT")
			os.Unsetenv("SCRAPER_INTERVAL")

			// Устанавливаем переменные окружения для теста
			for key, value := range tt.env {
				os.Setenv(key, value)
			}

			cfg, err := Load()
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && cfg == nil {
				t.Error("Load() returned nil config when no error expected")
				return
			}

			if !tt.wantErr {
				// Проверяем значения конфигурации
				if cfg.Gitlab.Token != tt.env["GITLAB_TOKEN"] {
					t.Errorf("Gitlab.Token = %v, want %v", cfg.Gitlab.Token, tt.env["GITLAB_TOKEN"])
				}
				if cfg.Gitlab.BaseURL != tt.env["GITLAB_BASE_URL"] {
					t.Errorf("Gitlab.BaseURL = %v, want %v", cfg.Gitlab.BaseURL, tt.env["GITLAB_BASE_URL"])
				}

				// Проверяем ProjectIDs
				expectedIDs := []int{}
				if tt.env["GITLAB_PROJECT_IDS"] != "" {
					// Парсим ожидаемые ID из строки
					idStrs := strings.Split(tt.env["GITLAB_PROJECT_IDS"], ",")
					for _, idStr := range idStrs {
						idStr = strings.TrimSpace(idStr)
						if idStr != "" {
							if id, err := strconv.Atoi(idStr); err == nil {
								expectedIDs = append(expectedIDs, id)
							}
						}
					}
				}

				if len([]int(cfg.Gitlab.ProjectIDs)) != len(expectedIDs) {
					t.Errorf("Gitlab.ProjectIDs length = %v, want %v", len([]int(cfg.Gitlab.ProjectIDs)), len(expectedIDs))
				} else {
					for i, expectedID := range expectedIDs {
						if []int(cfg.Gitlab.ProjectIDs)[i] != expectedID {
							t.Errorf("Gitlab.ProjectIDs[%d] = %v, want %v", i, []int(cfg.Gitlab.ProjectIDs)[i], expectedID)
						}
					}
				}
			}
		})
	}
}
