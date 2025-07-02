package config

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// ProjectIDsSlice - кастомный тип для парсинга списка ID
type ProjectIDsSlice []int

// GroupIDsSlice - кастомный тип для парсинга списка ID групп
type GroupIDsSlice []int

func (p *ProjectIDsSlice) Decode(value string) error {
	if value == "" {
		return fmt.Errorf("GITLAB_PROJECT_IDS is empty")
	}

	var projectIDs []int
	idStrs := strings.Split(value, ",")

	for _, idStr := range idStrs {
		idStr = strings.TrimSpace(idStr)
		if idStr == "" {
			continue
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			return fmt.Errorf("invalid project ID %q: %w", idStr, err)
		}

		projectIDs = append(projectIDs, id)
	}

	if len(projectIDs) == 0 {
		return fmt.Errorf("no valid project IDs found in GITLAB_PROJECT_IDS")
	}

	*p = projectIDs
	return nil
}

func (g *GroupIDsSlice) Decode(value string) error {
	if value == "" {
		// Группы не обязательны, поэтому возвращаем пустой слайс
		*g = []int{}
		return nil
	}

	var groupIDs []int
	idStrs := strings.Split(value, ",")

	for _, idStr := range idStrs {
		idStr = strings.TrimSpace(idStr)
		if idStr == "" {
			continue
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			return fmt.Errorf("invalid group ID %q: %w", idStr, err)
		}

		groupIDs = append(groupIDs, id)
	}

	*g = groupIDs
	return nil
}

type Config struct {
	Server struct {
		Port int `envconfig:"SERVER_PORT" default:"8080"`
	} `envconfig:"SERVER"`
	Gitlab struct {
		Token      string          `envconfig:"GITLAB_TOKEN" required:"true"`
		BaseURL    string          `envconfig:"GITLAB_BASE_URL" required:"true"`
		ProjectIDs ProjectIDsSlice `envconfig:"GITLAB_PROJECT_IDS" required:"true"`
		GroupIDs   GroupIDsSlice   `envconfig:"GITLAB_GROUP_IDS"`
	} `envconfig:"GITLAB"`
	Scraper struct {
		Interval time.Duration `envconfig:"SCRAPER_INTERVAL" default:"10s"`
	} `envconfig:"SCRAPER"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to process environment variables: %w", err)
	}

	return &cfg, nil
}
