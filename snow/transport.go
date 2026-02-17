package snow

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
)

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

// newRequest builds a request relative to the instance base URL.
// `p` should be like "/api/now/table/incident" (leading slash recommended).
func (c *Client) NewRequest(ctx context.Context, method, p string, query url.Values, body any) (*http.Request, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	u := *c.baseURL // copy
	// TODO: use url.JoinPath instead of path.Join
	u.Path = path.Join(c.baseURL.Path, p)

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

// do executes request and decodes JSON into out (if out != nil).
func (c *Client) Do(req *http.Request, out any) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read body once; used for both error parsing and decode
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return parseAPIError(resp.StatusCode, raw)
	}

	if out == nil {
		return nil
	}
	if len(raw) == 0 {
		return nil
	}

	if err := json.Unmarshal(raw, out); err != nil {
		return err
	}
	return nil
}

// DoWithResponse performs the request and returns both the raw response and unmarshals the body
func (c *Client) DoWithResponse(req *http.Request, out any) (*http.Response, error) {
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    // Read the body
    data, err := io.ReadAll(resp.Body)
    if err != nil {
        return resp, err
    }

    // Check for HTTP errors
    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        return resp, parseAPIError(resp.StatusCode, data)
    }

    // Unmarshal if output is provided
    if out != nil {
        if err := json.Unmarshal(data, out); err != nil {
            return resp, err
        }
    }

    return resp, nil
}