package test

import (
	"net/http"
	"os"

	"github.com/conductorone/baton-fluid-topics/pkg/client"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
)

var (
	Users = []map[string]interface{}{
		{
			"id":           "a061ccd9-3b8d-4f73-8d21-d045b3680a9d",
			"displayName":  "Kevin Dickson",
			"emailAddress": "kevin.dickson@powin.com",
		},
		{
			"id":           "e0067f57-46f9-4820-b51b-40898e087167",
			"displayName":  "A mateo",
			"emailAddress": "mateovespasiano.job@gmail.com",
		},
		{
			"id":           "1b70fa74-59ae-47a5-bf30-6e9d925f22e7",
			"displayName":  "Bjorn Tipling",
			"emailAddress": "bjorn.tipling@conductorone.com",
		},
		{
			"id":           "f0775b15-ccc6-4244-b050-64c2b9f965d6",
			"displayName":  "Desmond Cole",
			"emailAddress": "desmond.cole@powin.com",
		},
		{
			"id":           "3864d625-dce4-4c36-8fa1-394deb27c9e3",
			"displayName":  "AAAAA TEST",
			"emailAddress": "aaaaatestsandox@gmail.com",
		},
	}
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
func NewTestClient(response *http.Response, err error) *client.FluidTopicsClient {
	transport := &TestRoundTripper{response: response, err: err}
	httpClient := &http.Client{Transport: transport}
	baseHttpClient := uhttp.NewBaseHttpClient(httpClient)

	bearerToken := ""

	newClientT, _ := client.NewClient(bearerToken, baseHttpClient)

	return newClientT
}

func ReadFile(fileName string) (string, error) {
	data, err := os.ReadFile("../../test/mockResponses/" + fileName)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
