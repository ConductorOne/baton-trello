package connector

import (
	"context"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"github.com/conductorone/baton-trello/pkg/client"
	"github.com/conductorone/baton-trello/test"
)

// Tests that the client can fetch organizations based on the documented API below.
// https://developer.atlassian.com/cloud/trello/rest/api-group-organizations/#api-organizations-id-get
func TestTrelloClient_GetOrganizations(t *testing.T) {
	// Create a mock response.
	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body: io.NopCloser(strings.NewReader(`
			{
				"id": "1ed53893-6225-4d74-9806-3eedcbb402dd",
				"name": "organizationTest",
				"displayName": "Trello Workspace Test",
				"desc": "",
				"descData": {
					"emoji": {}
			  	},
			  	"url": "https://trello.com/w/organizationTest",
			  	"website": null,
			  	"teamType": null,
			  	"logoHash": null,
			  	"logoUrl": null,
			  	"offering": "trello.business_class",
			  	"products": [110],
			  	"powerUps": [110]
			}
		`)),
	}
	mockResponse.Header.Set("Content-Type", "application/json")

	// Create a test client with the mock response.
	testClient := test.NewTestClient(mockResponse, nil)
	testClient.WithOrganizationIDs(test.OrganizationIDs)

	// Call GetOrganizations
	ctx := context.Background()
	result, nextOptions, err := testClient.ListOrganizations(ctx)

	// Check for errors.
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify the result.
	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	// Check count.
	if len(*result) != 1 {
		t.Errorf("Expected Count to be 1, got %d", len(*result))
	}

	for index, organization := range *result {
		expectedOrg := client.Organization{
			ID:          "1ed53893-6225-4d74-9806-3eedcbb402dd",
			Name:        test.OrganizationIDs[index],
			DisplayName: "Trello Workspace Test",
			Url:         "https://trello.com/w/organizationTest",
			Offering:    "trello.business_class",
			Products:    []int{110},
			PowerUps:    []int{110},
		}

		if !reflect.DeepEqual(organization, expectedOrg) {
			t.Errorf("Unexpected organization: got %+v, want %+v", organization, expectedOrg)
		}
	}

	// Check next options.
	if nextOptions == nil {
		t.Fatal("Expected non-nil nextOptions")
	}
}

func TestTrelloClient_GetOrganizations_RequestDetails(t *testing.T) {
	// Create a custom RoundTripper to capture the request.
	var capturedRequest *http.Request
	mockTransport := &test.MockRoundTripper{
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{}`)),
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
	testClient := client.NewClient(baseHttpClient).WithApiKey("api-key").WithBearerToken("api-token").WithOrganizationIDs(test.OrganizationIDs)

	// Call GetOrganizations.
	ctx := context.Background()
	_, _, err := testClient.ListOrganizations(ctx)

	// Check for errors.
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify the request details.
	if capturedRequest == nil {
		t.Fatal("No request was captured")
	}

	// Check URL components.
	expectedURL := "https://api.trello.com/1/organizations/organizationTest?key=api-key&token=api-token"
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
