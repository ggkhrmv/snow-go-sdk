package pkg

import (
	"context"
	"net/http"
	"net/url"
)

// Requester is the minimal surface area subpackages need (table, attachment, etc.).
type Requester interface {
	NewRequest(ctx context.Context, method, path string, query url.Values, body any) (*http.Request, error)
	Do(req *http.Request, out any) error
}
