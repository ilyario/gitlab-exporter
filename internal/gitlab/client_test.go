package gitlab

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		baseURL string
		wantErr bool
	}{
		{
			name:    "valid client",
			token:   "test-token",
			baseURL: "https://gitlab.com",
			wantErr: false,
		},
		{
			name:    "empty token",
			token:   "",
			baseURL: "https://gitlab.com",
			wantErr: true,
		},
		{
			name:    "empty base URL",
			token:   "test-token",
			baseURL: "",
			wantErr: true,
		},
		{
			name:    "both empty",
			token:   "",
			baseURL: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.token, tt.baseURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && client == nil {
				t.Error("NewClient() returned nil client when no error expected")
			}

			if !tt.wantErr && client.client == nil {
				t.Error("NewClient() returned client with nil gitlab client")
			}
		})
	}
}

func TestClient_GetClient(t *testing.T) {
	client, err := NewClient("test-token", "https://gitlab.com")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	gitlabClient := client.GetClient()
	if gitlabClient == nil {
		t.Error("GetClient() returned nil")
	}
}
