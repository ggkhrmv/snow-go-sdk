package snow

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	ErrMissingInstanceURL = errors.New("missing instance URL: use WithInstanceURL")
	ErrMissingAuth        = errors.New("missing auth: use WithBasicAuth (or other auth options later)")
	ErrMissingHTTPClient  = errors.New("missing http client")
	ErrInvalidBasicAuth   = errors.New("basic auth requires non-empty username and password")
	ErrNilRequest         = errors.New("request is nil")
)

type APIError struct {
	Status  int
	Message string
	Detail  string
	Raw     []byte
}

func (e *APIError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("servicenow API error (%d): %s (%s)", e.Status, e.Message, e.Detail)
	}
	if e.Message != "" {
		return fmt.Sprintf("servicenow API error (%d): %s", e.Status, e.Message)
	}
	return fmt.Sprintf("servicenow API error (%d)", e.Status)
}

// ServiceNow often returns errors like: {"error": {"message":"...", "detail":"..."}, "status":"failure"}
// but it can vary across endpoints. We parse the common case and fall back to raw.
func parseAPIError(status int, raw []byte) error {
	type snErr struct {
		Error struct {
			Message string `json:"message"`
			Detail  string `json:"detail"`
		} `json:"error"`
	}
	var e snErr
	if err := json.Unmarshal(raw, &e); err == nil {
		if e.Error.Message != "" || e.Error.Detail != "" {
			return &APIError{
				Status:  status,
				Message: e.Error.Message,
				Detail:  e.Error.Detail,
				Raw:     raw,
			}
		}
	}

	return &APIError{Status: status, Raw: raw}
}
