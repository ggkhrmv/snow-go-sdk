package snow

import (
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	baseURL    *url.URL
	httpClient Doer
	auth       Auth

	userAgent string

	defaultHeaders http.Header
}

func NewClient(opts ...Option) (*Client, error) {
	c := &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		userAgent:  "servicenow-go-sdk/0.1",
		defaultHeaders: http.Header{
			"Accept": []string{"application/json"},
		},
	}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	if c.baseURL == nil {
		return nil, ErrMissingInstanceURL
	}
	if c.auth == nil {
		return nil, ErrMissingAuth
	}

	return c, nil
}
