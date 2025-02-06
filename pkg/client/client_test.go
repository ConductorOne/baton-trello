package client

import (
	"testing"
)

func TestNewTrelloClient(t *testing.T) {
	t.Run("Client URL", func(t *testing.T) {
		client := NewClient()
		if client.baseDomain != domain {
			t.Errorf("Expected baseURL to be %s, got %s", domain, client.baseDomain)
		}
	})
}

func TestTrelloClient_AddCredentials(t *testing.T) {
	client := NewClient()

	mockApiKey := "api-key"
	mockApiToken := "api-token"

	client.WithApiKey(mockApiKey).WithBearerToken(mockApiToken)

	if client.getApiKey() != mockApiKey {
		t.Errorf("Set API key failed. Expected %s, got %s", mockApiKey, client.getApiKey())
	}
	if client.getApiToken() != mockApiToken {
		t.Errorf("Set API token failed. Expected %s, got %s", mockApiToken, client.getApiToken())
	}
}
