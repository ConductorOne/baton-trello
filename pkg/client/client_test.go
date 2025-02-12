package client

import (
	"testing"
)

func TestNewTrelloClient(t *testing.T) {
	t.Run("Client URL", func(t *testing.T) {
		client := NewClient("", "", []string{})
		if client.BaseDomain != domain {
			t.Errorf("Expected baseURL to be %s, got %s", domain, client.BaseDomain)
		}
	})
}

func TestTrelloClient_AddCredentials(t *testing.T) {
	mockApiKey := "api-key"
	mockApiToken := "api-token"

	client := NewClient(mockApiKey, mockApiToken, []string{})

	if client.ApiKey != mockApiKey {
		t.Errorf("Set API key failed. Expected %s, got %s", mockApiKey, client.ApiKey)
	}
	if client.ApiToken != mockApiToken {
		t.Errorf("Set API token failed. Expected %s, got %s", mockApiToken, client.ApiToken)
	}
}
