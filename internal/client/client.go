package client

import "net/http"

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a new client
func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}
}
