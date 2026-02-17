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

// Convenience: default to dynamic records
func NewMap(r snow.Requester, tableName string) *Client[map[string]any] {
    return New[map[string]any](r, tableName)
}

func (c *Client[T]) basePath() (string, error) {
    if c.table == "" {
        return "", ErrInvalidTableName
    }
    return path.Join("/api/now/table", c.table), nil
}

// List retrieves multiple records with pagination metadata
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

    // Parse pagination headers if not suppressed
    if opts == nil || opts.SuppressPaginationHeader == nil || !*opts.SuppressPaginationHeader {
        listResp.Meta = parsePaginationHeaders(resp.Header)
    }

    return listResp, nil
}

// ListSimple returns just the records without metadata (backward compatible)
func (c *Client[T]) ListSimple(ctx context.Context, opts *ListOptions) ([]T, error) {
    resp, err := c.List(ctx, opts)
    if err != nil {
        return nil, err
    }
    return resp.Result, nil
}

func (c *Client[T]) Get(ctx context.Context, sysID string, opts *GetOptions) (*GetResponse[T], error) {
    var zero *GetResponse[T]

    base, err := c.basePath()
    if err != nil {
        return zero, err
    }
    if strings.TrimSpace(sysID) == "" {
        return zero, ErrInvalidSysID
    }

    q := url.Values{}
    if opts != nil {
        if err := opts.apply(q); err != nil {
            return zero, err
        }
    }

    req, err := c.r.NewRequest(ctx, http.MethodGet, path.Join(base, sysID), q, nil)
    if err != nil {
        return zero, err
    }

    var out resultOne[T]
    if err := c.r.Do(req, &out); err != nil {
        return zero, err
    }

    return &GetResponse[T]{Result: out.Result}, nil
}

// GetSimple returns just the record (backward compatible)
func (c *Client[T]) GetSimple(ctx context.Context, sysID string, opts *GetOptions) (T, error) {
    var zero T
    resp, err := c.Get(ctx, sysID, opts)
    if err != nil {
        return zero, err
    }
    return resp.Result, nil
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

    base, err := c.basePath()
    if err != nil {
        return zero, err
    }
    if strings.TrimSpace(sysID) == "" {
        return zero, ErrInvalidSysID
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

    req, err := c.r.NewRequest(ctx, http.MethodPatch, path.Join(base, sysID), q, in)
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

    base, err := c.basePath()
    if err != nil {
        return zero, err
    }
    if strings.TrimSpace(sysID) == "" {
        return zero, ErrInvalidSysID
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

    req, err := c.r.NewRequest(ctx, http.MethodPut, path.Join(base, sysID), q, in)
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
    base, err := c.basePath()
    if err != nil {
        return err
    }
    if strings.TrimSpace(sysID) == "" {
        return ErrInvalidSysID
    }

    q := url.Values{}
    if opts != nil {
        if err := opts.apply(q); err != nil {
            return err
        }
    }

    req, err := c.r.NewRequest(ctx, http.MethodDelete, path.Join(base, sysID), q, nil)
    if err != nil {
        return err
    }
    return c.r.Do(req, nil)
}

// parsePaginationHeaders extracts pagination metadata from response headers
func parsePaginationHeaders(headers http.Header) *PaginationMeta {
    meta := &PaginationMeta{}

    // Parse Link header: <url>; rel="first", <url>; rel="next", etc.
    if link := headers.Get("Link"); link != "" {
        parts := strings.Split(link, ",")
        for _, part := range parts {
            part = strings.TrimSpace(part)
            if strings.Contains(part, `rel="first"`) {
                meta.First = extractURL(part)
            } else if strings.Contains(part, `rel="prev"`) {
                meta.Prev = extractURL(part)
            } else if strings.Contains(part, `rel="next"`) {
                meta.Next = extractURL(part)
            } else if strings.Contains(part, `rel="last"`) {
                meta.Last = extractURL(part)
            }
        }
    }

    // Parse X-Total-Count header
    if count := headers.Get("X-Total-Count"); count != "" {
        if n, err := strconv.Atoi(count); err == nil {
            meta.TotalCount = n
        }
    }

    return meta
}

func extractURL(linkPart string) string {
    start := strings.Index(linkPart, "<")
    end := strings.Index(linkPart, ">")
    if start >= 0 && end > start {
        return linkPart[start+1 : end]
    }
    return ""
}
