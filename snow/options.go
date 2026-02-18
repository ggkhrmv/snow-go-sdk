package snow

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
)

type Option func(*Client) error

func WithInstanceURL(instance string) Option {
	return func(c *Client) error {
		u, err := url.Parse(instance)
		if err != nil {
			return err
		}
		if u.Scheme == "" || u.Host == "" {
			return errors.New("instance URL must include scheme and host, e.g. https://dev12345.service-now.com")
		}

		// normalize: no trailing slash
		u.Path = strings.TrimRight(u.Path, "/")
		u.RawQuery = ""
		u.Fragment = ""
		c.baseURL = u
		return nil
	}
}

func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) error {
		if hc == nil {
			return errors.New("http client is nil")
		}
		c.httpClient = hc
		return nil
	}
}

func WithUserAgent(ua string) Option {
	return func(c *Client) error {
		if strings.TrimSpace(ua) == "" {
			return errors.New("user agent cannot be empty")
		}
		c.userAgent = ua
		return nil
	}
}

func WithDefaultHeader(key, value string) Option {
	return func(c *Client) error {
		if strings.TrimSpace(key) == "" {
			return errors.New("header key cannot be empty")
		}
		c.defaultHeaders.Add(key, value)
		return nil
	}
}
