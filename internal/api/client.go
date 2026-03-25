// Package api provides the HTTP client for the Exasol SaaS API.
package api

const DefaultBaseURL = "https://cloud.exasol.com"

// DefaultProvider is the cloud provider used for all Exasol SaaS resources.
const DefaultProvider = "aws"

// Client holds credentials and configuration for communicating with the Exasol SaaS API.
type Client struct {
	token   string
	baseURL string
}

// NewClient creates a Client that authenticates using the provided token.
func NewClient(token, baseURL string) *Client {
	return &Client{
		token:   token,
		baseURL: baseURL,
	}
}
