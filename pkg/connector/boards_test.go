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

var expectedMemberships = [][]client.User{
	{
		{ID: test.UserIDs[0],
			MemberID:   test.UserIDs[0],
			MemberType: "admin"},
		{ID: test.UserIDs[1],
			MemberID:   test.UserIDs[1],
			MemberType: "normal"},
	},
	{
		{ID: test.UserIDs[0],
			MemberID:   test.UserIDs[0],
			MemberType: "admin"},
	},
}

// Tests that the client can fetch boards based on the documented API below.
// https://developer.atlassian.com/cloud/trello/rest/api-group-organizations/#api-organizations-id-boards-get
func TestTrelloClient_GetBoards(t *testing.T) {
	// Create a mock response.
	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(test.ReadFile("boardsMock.json"))),
	}
	mockResponse.Header.Set("Content-Type", "application/json")

	// Create a test client with the mock response.
	testClient := test.NewTestClient(mockResponse, nil)
	testClient.WithOrganizationIDs(test.OrganizationIDs)

	// Call GetBoards
	ctx := context.Background()
	result, nextOptions, err := testClient.ListBoards(ctx)

	// Check for errors.
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify the result.
	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	// Check count.
	if len(*result) != 2 {
		t.Errorf("Expected Count to be 2, got %d", len(*result))
	}

	for index, board := range *result {
		invitations := "members"
		if index == 0 {
			invitations = "admins"
		}
		expectedBoard := client.Board{
			ID:             test.BoardIDs[index],
			Name:           fmt.Sprintf("Test %d", index+1),
			Closed:         false,
			IdOrganization: test.OrganizationIDs[0],
			Pinned:         false,
			Url:            fmt.Sprintf("https://trello.com/b/test/test%d", index+1),
			Preferences: client.Preferences{
				PermissionLevel:     "org",
				HideVotes:           false,
				Voting:              "disabled",
				Comments:            "members",
				Invitations:         invitations,
				SelfJoin:            true,
				CardCovers:          true,
				ShowCompleteStatus:  true,
				CardCounts:          false,
				IsTemplate:          false,
				CardAging:           "regular",
				CalendarFeedEnabled: false,
				CanBePublic:         false,
				CanBeEnterprise:     false,
				CanBeOrg:            false,
				CanBePrivate:        false,
				CanInvite:           true,
			},
			Memberships: expectedMemberships[index],
		}

		if !reflect.DeepEqual(board, expectedBoard) {
			t.Errorf("Unexpected board: got %+v, want %+v", board, expectedBoard)
		}
	}

	// Check next options.
	if nextOptions == nil {
		t.Fatal("Expected non-nil nextOptions")
	}
}

func TestTrelloClient_GetBoards_RequestDetails(t *testing.T) {
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
	testClient := client.NewClient(baseHttpClient).WithApiKey("api-key").WithBearerToken("api-token").WithOrganizationIDs(test.OrganizationIDs)

	// Call GetBoards.
	ctx := context.Background()
	_, _, err := testClient.ListBoards(ctx)

	// Check for errors.
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify the request details.
	if capturedRequest == nil {
		t.Fatal("No request was captured")
	}

	// Check URL components.
	expectedURL := "https://api.trello.com/1/organizations/organizationTest/boards?key=api-key&token=api-token"
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
