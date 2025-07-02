package scraper

import (
	"context"
	"log"
	"strconv"
	"time"

	"ru/mvideo/com/gitlab/token-exporter/internal/gitlab"
	"ru/mvideo/com/gitlab/token-exporter/internal/metrics"
)

type TokenScraper struct {
	gitlabClient         gitlab.GitLabClientInterface
	metrics              *metrics.Handler
	projectIDs           []int
	groupIDs             []int
	knownTokens          map[string]bool // для отслеживания токенов между скрейпингами
	knownUserTokens      map[string]bool // для отслеживания пользовательских токенов между скрейпингами
	knownGroupTokens     map[string]bool // для отслеживания групповых токенов между скрейпингами
	currentProjectTokens map[string]bool
	currentUserTokens    map[string]bool
	currentGroupTokens   map[string]bool
}

func NewTokenScraper(gitlabClient gitlab.GitLabClientInterface, metrics *metrics.Handler, projectIDs []int, groupIDs []int) *TokenScraper {
	return &TokenScraper{
		gitlabClient:         gitlabClient,
		metrics:              metrics,
		projectIDs:           projectIDs,
		groupIDs:             groupIDs,
		knownTokens:          make(map[string]bool),
		knownUserTokens:      make(map[string]bool),
		knownGroupTokens:     make(map[string]bool),
		currentProjectTokens: make(map[string]bool),
		currentUserTokens:    make(map[string]bool),
		currentGroupTokens:   make(map[string]bool),
	}
}

func (s *TokenScraper) Start(ctx context.Context, interval time.Duration) {
	log.Printf("Starting token scraper with interval: %v", interval)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	s.scrape()

	for {
		select {
		case <-ctx.Done():
			log.Println("Token scraper stopped")
			return
		case <-ticker.C:
			s.scrape()
		}
	}
}

func (s *TokenScraper) scrape() {
	start := time.Now()
	log.Println("Starting token scrape...")

	now := time.Now()

	totalProjectTokens := s.scrapeProjectTokens(now)

	totalUserTokens := s.scrapeUserTokens(now)

	totalGroupTokens := s.scrapeGroupTokens(now)

	for knownToken := range s.knownTokens {
		if !s.currentProjectTokens[knownToken] {
			s.metrics.DeleteTokenMetrics(knownToken)
			log.Printf("Removed metrics for deleted project token: %s", knownToken)
		}
	}

	for knownUserToken := range s.knownUserTokens {
		if !s.currentUserTokens[knownUserToken] {
			s.metrics.DeleteUserTokenMetrics(knownUserToken)
			log.Printf("Removed metrics for deleted user token: %s", knownUserToken)
		}
	}

	for knownGroupToken := range s.knownGroupTokens {
		if !s.currentGroupTokens[knownGroupToken] {
			s.metrics.DeleteGroupTokenMetrics(knownGroupToken)
			log.Printf("Removed metrics for deleted group token: %s", knownGroupToken)
		}
	}

	s.knownTokens = s.currentProjectTokens
	s.knownUserTokens = s.currentUserTokens
	s.knownGroupTokens = s.currentGroupTokens

	s.metrics.SetTotalTokens(totalProjectTokens)
	s.metrics.SetTotalUserTokens(totalUserTokens)
	s.metrics.SetTotalGroupTokens(totalGroupTokens)
	s.metrics.SetLastScrapeTime(now)
	duration := time.Since(start)
	s.metrics.RecordScrapeDuration(duration)
	log.Printf("Token scrape completed in %v, found %d project tokens, %d user tokens, %d group tokens", duration, totalProjectTokens, totalUserTokens, totalGroupTokens)
}

func (s *TokenScraper) scrapeProjectTokens(now time.Time) int {
	totalTokens := 0
	s.currentProjectTokens = make(map[string]bool)

	for _, projectID := range s.projectIDs {
		tokens, err := s.gitlabClient.GetProjectAccessTokens(projectID)

		if err != nil {
			log.Printf("Failed to get project access tokens for project %d: %v", projectID, err)
			s.metrics.IncrementScrapeErrors()
			continue
		}

		projectName, err := s.gitlabClient.GetProjectName(projectID)
		if err != nil {
			log.Printf("Failed to get project name for project %d: %v", projectID, err)
			s.metrics.IncrementScrapeErrors()
			continue
		}

		for _, token := range tokens {
			metrics_name := projectName + " " + strconv.Itoa(projectID) + " " + token.Name

			s.currentProjectTokens[metrics_name] = true

			expiresAt := time.Time(*token.ExpiresAt)
			isExpired := expiresAt.Before(now)

			s.metrics.SetTokenExpiresAt(metrics_name, expiresAt)
			s.metrics.SetTokenIsExpired(metrics_name, isExpired)

			log.Printf("Project: %d, Token: %s, Expires: %s, IsExpired: %t", projectID, token.Name, expiresAt.Format(time.RFC3339), isExpired)
		}
		totalTokens += len(tokens)
	}

	return totalTokens
}

func (s *TokenScraper) scrapeUserTokens(now time.Time) int {
	totalUserTokens := 0
	s.currentUserTokens = make(map[string]bool)

	userTokens, err := s.gitlabClient.GetUserAccessTokens()

	if err != nil {
		log.Printf("Failed to get user access tokens: %v", err)
		s.metrics.IncrementScrapeErrors()
		return 0
	}

	for _, token := range userTokens {
		userName, err := s.gitlabClient.GetUserName(token.UserID)

		if err != nil {
			log.Printf("Failed to get user name for token %d: %v", token.UserID, err)
			s.metrics.IncrementScrapeErrors()
			userName = "Unknown user"
		}

		metrics_name := userName + " " + strconv.Itoa(token.UserID) + " " + token.Name

		s.currentUserTokens[metrics_name] = true

		expiresAt := time.Time(*token.ExpiresAt)
		isExpired := expiresAt.Before(now)

		s.metrics.SetUserTokenExpiresAt(metrics_name, expiresAt)
		s.metrics.SetUserTokenIsExpired(metrics_name, isExpired)

		log.Printf("User: %s, Token: %s, Expires: %s, IsExpired: %t", userName, token.Name, expiresAt.Format(time.RFC3339), isExpired)
	}
	totalUserTokens = len(userTokens)

	return totalUserTokens
}

func (s *TokenScraper) scrapeGroupTokens(now time.Time) int {
	totalGroupTokens := 0
	s.currentGroupTokens = make(map[string]bool)

	for _, groupID := range s.groupIDs {
		tokens, err := s.gitlabClient.GetGroupAccessTokens(groupID)

		if err != nil {
			log.Printf("Failed to get group access tokens for group %d: %v", groupID, err)
			s.metrics.IncrementScrapeErrors()
			continue
		}

		groupName, err := s.gitlabClient.GetGroupName(groupID)
		if err != nil {
			log.Printf("Failed to get group name for group %d: %v", groupID, err)
			s.metrics.IncrementScrapeErrors()
			continue
		}

		for _, token := range tokens {
			metrics_name := groupName + " " + strconv.Itoa(groupID) + " " + token.Name

			s.currentGroupTokens[metrics_name] = true

			expiresAt := time.Time(*token.ExpiresAt)
			isExpired := expiresAt.Before(now)

			s.metrics.SetGroupTokenExpiresAt(metrics_name, expiresAt)
			s.metrics.SetGroupTokenIsExpired(metrics_name, isExpired)

			log.Printf("Group: %d, Token: %s, Expires: %s, IsExpired: %t", groupID, token.Name, expiresAt.Format(time.RFC3339), isExpired)
		}
		totalGroupTokens += len(tokens)
	}

	return totalGroupTokens
}
