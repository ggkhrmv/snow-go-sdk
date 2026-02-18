package table

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"

	snow "github.com/ggkhrmv/snow-go-sdk/snow"
)

var (
	ErrInvalidTableName = errors.New("invalid table name")
	ErrInvalidSysID     = errors.New("invalid sys_id")
	ErrNilInput         = errors.New("input body is nil")
	ErrNilRequester     = errors.New("requester is nil")
)

// Client is a Table API client bound to a specific table, returning records of type T.
type Client[T any] struct {
	r     snow.Requester
	table string
}

func New[T any](r snow.Requester, tableName string) *Client[T] {
	return &Client[T]{
		r:     r,
		table: strings.TrimSpace(tableName),
	}
}

// NewMap is a convenience constructor returning dynamic records.
func NewMap(r snow.Requester, tableName string) *Client[map[string]any] {
	return New[map[string]any](r, tableName)
}

func (c *Client[T]) basePath() (string, error) {
	if c == nil || c.r == nil {
		return "", ErrNilRequester
	}
	if c.table == "" || strings.ContainsAny(c.table, `/\\`) {
		return "", ErrInvalidTableName
	}
	return path.Join("/api/now/table", c.table), nil
}

func (c *Client[T]) recordPath(sysID string) (string, error) {
	base, err := c.basePath()
	if err != nil {
		return "", err
	}

	sysID = strings.TrimSpace(sysID)
	if sysID == "" || strings.ContainsAny(sysID, `/\\`) {
		return "", ErrInvalidSysID
	}

	return path.Join(base, sysID), nil
}

// List retrieves multiple records with pagination metadata.
func (c *Client[T]) List(ctx context.Context, opts *ListOptions) (*ListResponse[T], error) {
	base, err := c.basePath()
	if err != nil {
		return nil, err
	}

	q := url.Values{}
	if opts != nil {
		if err := opts.apply(q); err != nil {
			return nil, err
		}
	}

	req, err := c.r.NewRequest(ctx, http.MethodGet, base, q, nil)
	if err != nil {
		return nil, err
	}

	var out resultList[T]
	resp, err := c.r.DoWithResponse(req, &out)
	if err != nil {
		return nil, err
	}

	listResp := &ListResponse[T]{
		Result: out.Result,
	}

	if opts == nil || opts.SuppressPaginationHeader == nil || !*opts.SuppressPaginationHeader {
		listResp.Meta = parsePaginationHeaders(resp.Header)
	}

	return listResp, nil
}

func (c *Client[T]) Get(ctx context.Context, sysID string, opts *GetOptions) (*GetResponse[T], error) {
	var zero *GetResponse[T]

	recordPath, err := c.recordPath(sysID)
	if err != nil {
		return zero, err
	}

	q := url.Values{}
	if opts != nil {
		if err := opts.apply(q); err != nil {
			return zero, err
		}
	}

	req, err := c.r.NewRequest(ctx, http.MethodGet, recordPath, q, nil)
	if err != nil {
		return zero, err
	}

	var out resultOne[T]
	if err := c.r.Do(req, &out); err != nil {
		return zero, err
	}

	return &GetResponse[T]{Result: out.Result}, nil
}

func (c *Client[T]) Create(ctx context.Context, in any, opts *WriteOptions) (*WriteResponse[T], error) {
	var zero *WriteResponse[T]

	base, err := c.basePath()
	if err != nil {
		return zero, err
	}
	if in == nil {
		return zero, ErrNilInput
	}

	q := url.Values{}
	if opts != nil {
		if err := opts.apply(q); err != nil {
			return zero, err
		}
	}

	req, err := c.r.NewRequest(ctx, http.MethodPost, base, q, in)
	if err != nil {
		return zero, err
	}

	var out resultOne[T]
	if err := c.r.Do(req, &out); err != nil {
		return zero, err
	}

	return &WriteResponse[T]{Result: out.Result}, nil
}

func (c *Client[T]) Update(ctx context.Context, sysID string, in any, opts *WriteOptions) (*WriteResponse[T], error) {
	var zero *WriteResponse[T]

	recordPath, err := c.recordPath(sysID)
	if err != nil {
		return zero, err
	}
	if in == nil {
		return zero, ErrNilInput
	}

	q := url.Values{}
	if opts != nil {
		if err := opts.apply(q); err != nil {
			return zero, err
		}
	}

	req, err := c.r.NewRequest(ctx, http.MethodPatch, recordPath, q, in)
	if err != nil {
		return zero, err
	}

	var out resultOne[T]
	if err := c.r.Do(req, &out); err != nil {
		return zero, err
	}

	return &WriteResponse[T]{Result: out.Result}, nil
}

func (c *Client[T]) Replace(ctx context.Context, sysID string, in any, opts *WriteOptions) (*WriteResponse[T], error) {
	var zero *WriteResponse[T]

	recordPath, err := c.recordPath(sysID)
	if err != nil {
		return zero, err
	}
	if in == nil {
		return zero, ErrNilInput
	}

	q := url.Values{}
	if opts != nil {
		if err := opts.apply(q); err != nil {
			return zero, err
		}
	}

	req, err := c.r.NewRequest(ctx, http.MethodPut, recordPath, q, in)
	if err != nil {
		return zero, err
	}

	var out resultOne[T]
	if err := c.r.Do(req, &out); err != nil {
		return zero, err
	}

	return &WriteResponse[T]{Result: out.Result}, nil
}

func (c *Client[T]) Delete(ctx context.Context, sysID string, opts *DeleteOptions) error {
	recordPath, err := c.recordPath(sysID)
	if err != nil {
		return err
	}

	q := url.Values{}
	if opts != nil {
		if err := opts.apply(q); err != nil {
			return err
		}
	}

	req, err := c.r.NewRequest(ctx, http.MethodDelete, recordPath, q, nil)
	if err != nil {
		return err
	}

	return c.r.Do(req, nil)
}

func parsePaginationHeaders(headers http.Header) *PaginationMeta {
	if headers == nil {
		return nil
	}

	meta := &PaginationMeta{}
	hasMeta := false

	if link := headers.Get("Link"); link != "" {
		parts := strings.Split(link, ",")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			url := extractURL(part)
			switch {
			case strings.Contains(part, `rel="first"`):
				meta.First = url
			case strings.Contains(part, `rel="prev"`):
				meta.Prev = url
			case strings.Contains(part, `rel="next"`):
				meta.Next = url
			case strings.Contains(part, `rel="last"`):
				meta.Last = url
			}
			if url != "" {
				hasMeta = true
			}
		}
	}

	if count := headers.Get("X-Total-Count"); count != "" {
		if n, err := strconv.Atoi(count); err == nil {
			meta.TotalCount = n
			hasMeta = true
		}
	}

	for _, rawURL := range []string{meta.Next, meta.Prev, meta.First, meta.Last} {
		if rawURL == "" {
			continue
		}
		if limit, offset, ok := parseLimitOffset(rawURL); ok {
			meta.Limit = limit
			meta.Offset = offset
			break
		}
	}

	if !hasMeta {
		return nil
	}

	return meta
}

func parseLimitOffset(rawURL string) (limit int, offset int, ok bool) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return 0, 0, false
	}

	q := u.Query()
	limitStr := q.Get("sysparm_limit")
	offsetStr := q.Get("sysparm_offset")
	if limitStr == "" && offsetStr == "" {
		return 0, 0, false
	}

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 0 {
			return 0, 0, false
		}
	}

	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			return 0, 0, false
		}
	}

	return limit, offset, true
}

func extractURL(linkPart string) string {
	start := strings.Index(linkPart, "<")
	end := strings.Index(linkPart, ">")
	if start >= 0 && end > start {
		return linkPart[start+1 : end]
	}
	return ""
}
