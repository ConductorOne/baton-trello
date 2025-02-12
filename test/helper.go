package test

import (
	"log"
	"net/http"
	"os"

	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"github.com/conductorone/baton-trello/pkg/client"
)

var (
	UserIDs         = []string{"ea960e6c-f613-4bed-8852-ab012603915b", "8b21d0aa-39a4-4c09-86d2-d29dff8d261f"}
	OrganizationIDs = []string{"organizationTest"}
	BoardIDs        = []string{"f7a6a858-ab65-4524-9632-b64a21aa3c79", "eef3dd14-929f-4b85-b601-7cc4a484fa97"}
)

// Custom RoundTripper for testing.
type TestRoundTripper struct {
	response *http.Response
	err      error
}

type MockRoundTripper struct {
	Response  *http.Response
	Err       error
	roundTrip func(*http.Request) (*http.Response, error)
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.roundTrip(req)
}

func (m *MockRoundTripper) SetRoundTrip(roundTrip func(*http.Request) (*http.Response, error)) {
	m.roundTrip = roundTrip
}

func (t *TestRoundTripper) RoundTrip(*http.Request) (*http.Response, error) {
	return t.response, t.err
}

// Helper function to create a test client with custom transport.
func NewTestClient(response *http.Response, err error) *client.TrelloClient {
	transport := &TestRoundTripper{response: response, err: err}
	httpClient := &http.Client{Transport: transport}
	baseHttpClient := uhttp.NewBaseHttpClient(httpClient)
	return client.NewClient("", "", OrganizationIDs, baseHttpClient)
}

func ReadFile(fileName string) string {
	data, err := os.ReadFile("../../test/mockResponses/" + fileName)
	if err != nil {
		log.Fatal(err)
	}

	return string(data)
}
