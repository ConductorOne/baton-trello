package connector

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"github.com/conductorone/baton-trello/pkg/client"
	"github.com/conductorone/baton-trello/test"
)

// Tests that the client can fetch users based on the documented API below.
// https://developer.atlassian.com/cloud/trello/rest/api-group-organizations/#api-organizations-id-members-get
func TestTrelloClient_GetUsers(t *testing.T) {
	// Create a mock response.
	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(test.ReadFile("usersMock.json"))),
	}
	mockResponse.Header.Set("Content-Type", "application/json")

	// Create a test client with the mock response.
	testClient := test.NewTestClient(mockResponse, nil)

	// Call GetUsers
	ctx := context.Background()
	result, nextOptions, err := testClient.ListUsers(ctx)

	// Check for errors.
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify the result.
	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	// Check count.
	if len(result) != 2 {
		t.Errorf("Expected Count to be 2, got %d", len(result))
	}

	for index, user := range result {
		expectedUser := client.User{
			ID:       test.UserIDs[index],
			Username: fmt.Sprintf("tester%d", index+1),
			Name:     fmt.Sprintf("Test User %d", index+1),
		}

		if !reflect.DeepEqual(user, expectedUser) {
			t.Errorf("Unexpected user: got %+v, want %+v", user, expectedUser)
		}
	}

	// Check next options.
	if nextOptions == nil {
		t.Fatal("Expected non-nil nextOptions")
	}
}

func TestTrelloClient_GetUsers_RequestDetails(t *testing.T) {
	// Create a custom RoundTripper to capture the request.
	var capturedRequest *http.Request
	mockTransport := &test.MockRoundTripper{
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`[]`)),
			Header:     make(http.Header),
		},
		Err: nil,
	}
	mockTransport.Response.Header.Set("Content-Type", "application/json")

	mockRoundTrip := func(req *http.Request) (*http.Response, error) {
		capturedRequest = req
		return mockTransport.Response, mockTransport.Err
	}
	mockTransport.SetRoundTrip(mockRoundTrip)

	// Create a test client with the mock transport.
	httpClient := &http.Client{Transport: mockTransport}
	baseHttpClient := uhttp.NewBaseHttpClient(httpClient)
	testClient := client.NewClient("api-key", "api-token", test.OrganizationIDs, baseHttpClient)

	// Call GetUsers.
	ctx := context.Background()
	_, _, err := testClient.ListUsers(ctx)

	// Check for errors.
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify the request details.
	if capturedRequest == nil {
		t.Fatal("No request was captured")
	}

	// Check URL components.
	expectedURL := "https://api.trello.com/1/organizations/organizationTest/members?key=api-key&token=api-token"
	if capturedRequest.URL.String() != expectedURL {
		t.Errorf("Expected URL %s, got %s", expectedURL, capturedRequest.URL.String())
	}

	// Check headers.
	expectedHeaders := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	for key, expectedValue := range expectedHeaders {
		if value := capturedRequest.Header.Get(key); value != expectedValue {
			t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, value)
		}
	}
}

// Test that the client can fetch users based on the documented API below.
// https://developer.atlassian.com/cloud/trello/rest/api-group-organizations/#api-organizations-id-memberships-get
func TestTrelloClient_GetUserDetails(t *testing.T) {
	// Create a mock response.
	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body: io.NopCloser(strings.NewReader(`
			{
				"id": "ea960e6c-f613-4bed-8852-ab012603915b",
				"idMember": "ea960e6c-f613-4bed-8852-ab012603915b",
				"memberType": "normal",
				"unconfirmed": false,
				"deactivated": false,
				"lastActive": "2025-02-05T17:34:03.386Z"
			}
		`)),
	}
	mockResponse.Header.Set("Content-Type", "application/json")

	// Create a test client with the mock response.
	testClient := test.NewTestClient(mockResponse, nil)

	// Call GetUsers
	ctx := context.Background()

	userID := test.UserIDs[0]
	user, annotations, err := testClient.GetMemberDetails(ctx, userID)

	// Check for errors.
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify the result.
	if user == nil {
		t.Fatal("Expected non-nil result")
	}

	expectedUser := client.User{
		ID:         userID,
		MemberID:   userID,
		MemberType: "normal",
	}

	if !reflect.DeepEqual(*user, expectedUser) {
		t.Errorf("Unexpected user: got %+v, want %+v", user, expectedUser)
	}

	// Check annotations.
	if annotations == nil {
		t.Fatal("Expected non-nil annotations")
	}
}
