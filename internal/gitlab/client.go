package gitlab

import (
	"fmt"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// GitLabClientInterface - интерфейс для клиента GitLab
type GitLabClientInterface interface {
	GetProjectAccessTokens(projectID int) ([]*gitlab.ProjectAccessToken, error)
	GetProjectName(projectID int) (string, error)
	GetUserAccessTokens() ([]*gitlab.PersonalAccessToken, error)
	GetUserName(userID int) (string, error)
	GetGroupAccessTokens(groupID int) ([]*gitlab.GroupAccessToken, error)
	GetGroupName(groupID int) (string, error)
	GetClient() *gitlab.Client
}

type Client struct {
	client *gitlab.Client
}

// Убеждаемся, что Client реализует GitLabClientInterface
var _ GitLabClientInterface = (*Client)(nil)

func NewClient(token, baseURL string) (*Client, error) {
	if token == "" {
		return nil, fmt.Errorf("gitlab token is required")
	}
	if baseURL == "" {
		return nil, fmt.Errorf("gitlab base URL is required")
	}

	client, err := gitlab.NewClient(token, gitlab.WithBaseURL(baseURL))
	if err != nil {
		return nil, fmt.Errorf("failed to create gitlab client: %w", err)
	}

	return &Client{client: client}, nil
}

func (c *Client) GetProjectAccessTokens(projectID int) ([]*gitlab.ProjectAccessToken, error) {
	tokens, _, err := c.client.ProjectAccessTokens.ListProjectAccessTokens(projectID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list project access tokens: %w", err)
	}

	return tokens, nil
}

func (c *Client) GetClient() *gitlab.Client {
	return c.client
}

func (c *Client) GetProjectName(projectID int) (string, error) {
	project, _, err := c.client.Projects.GetProject(projectID, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get project name: %w", err)
	}
	return project.Name, nil
}

func (c *Client) GetUserAccessTokens() ([]*gitlab.PersonalAccessToken, error) {
	state := "active"
	revoked := false
	options := &gitlab.ListPersonalAccessTokensOptions{
		State:   &state,
		Revoked: &revoked,
	}
	tokens, _, err := c.client.PersonalAccessTokens.ListPersonalAccessTokens(options)
	if err != nil {
		return nil, fmt.Errorf("failed to list user access tokens: %w", err)
	}

	return tokens, nil
}

func (c *Client) GetUserName(userID int) (string, error) {
	user, _, err := c.client.Users.GetUser(userID, gitlab.GetUsersOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to get user name: %w", err)
	}
	return user.Name, nil
}

func (c *Client) GetGroupAccessTokens(groupID int) ([]*gitlab.GroupAccessToken, error) {
	state := gitlab.AccessTokenStateActive
	revoked := false
	options := &gitlab.ListGroupAccessTokensOptions{
		State:   &state,
		Revoked: &revoked,
	}
	tokens, _, err := c.client.GroupAccessTokens.ListGroupAccessTokens(groupID, options)
	if err != nil {
		return nil, fmt.Errorf("failed to list group access tokens: %w", err)
	}

	return tokens, nil
}

func (c *Client) GetGroupName(groupID int) (string, error) {
	group, _, err := c.client.Groups.GetGroup(groupID, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get group name: %w", err)
	}
	return group.Name, nil
}
