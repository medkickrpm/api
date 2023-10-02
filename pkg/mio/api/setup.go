package mio_api

const baseURL = "https://api.connect.mio-labs.com/v1"

type Client struct {
	APIKey string
}

// NewClient initializes a new API client
func NewClient(apiKey string) *Client {
	return &Client{
		APIKey: apiKey,
	}
}
