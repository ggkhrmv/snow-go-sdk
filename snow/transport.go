package snow

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

// NewRequest builds a request relative to the instance base URL.
// `p` should be like "/api/now/table/incident" (leading slash recommended).
func (c *Client) NewRequest(ctx context.Context, method, p string, query url.Values, body any) (*http.Request, error) {
	if c == nil || c.baseURL == nil {
		return nil, ErrMissingInstanceURL
	}

	if ctx == nil {
		ctx = context.Background()
	}

	joined, err := url.JoinPath(c.baseURL.String(), p)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(joined)
	if err != nil {
		return nil, err
	}

	if len(query) > 0 {
		u.RawQuery = query.Encode()
	}

	var r io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		r = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), r)
	if err != nil {
		return nil, err
	}

	// default headers
	for k, vals := range c.defaultHeaders {
		for _, v := range vals {
			req.Header.Add(k, v)
		}
	}
	req.Header.Set("User-Agent", c.userAgent)

	// content-type only when we send a body
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// auth
	if c.auth != nil {
		c.auth.Apply(req)
	}

	return req, nil
}

// Do executes request and decodes JSON into out (if out != nil).
func (c *Client) Do(req *http.Request, out any) error {
	_, err := c.do(req, out, false)
	return err
}

// DoWithResponse performs the request and returns both the raw response and unmarshals the body.
func (c *Client) DoWithResponse(req *http.Request, out any) (*http.Response, error) {
	return c.do(req, out, true)
}

func (c *Client) do(req *http.Request, out any, preserveBody bool) (*http.Response, error) {
	if req == nil {
		return nil, ErrNilRequest
	}
	if c == nil || c.httpClient == nil {
		return nil, ErrMissingHTTPClient
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	raw, err := io.ReadAll(resp.Body)
	closeErr := resp.Body.Close()
	if err != nil {
		return resp, err
	}
	if closeErr != nil {
		return resp, closeErr
	}

	if preserveBody {
		resp.Body = io.NopCloser(bytes.NewReader(raw))
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return resp, parseAPIError(resp.StatusCode, raw)
	}

	if out == nil || len(raw) == 0 {
		return resp, nil
	}

	if err := json.Unmarshal(raw, out); err != nil {
		return resp, err
	}

	return resp, nil
}
