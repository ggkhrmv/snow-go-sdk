package pkg

import "net/http"

type Auth interface {
	Apply(req *http.Request)
}

type basicAuth struct {
	username string
	password string
}

func (a basicAuth) Apply(req *http.Request) {
	req.SetBasicAuth(a.username, a.password)
}

func WithBasicAuth(username, password string) Option {
	return func(c *Client) error {
		if username == "" || password == "" {
			return ErrInvalidBasicAuth
		}
		c.auth = basicAuth{username: username, password: password}
		return nil
	}
}