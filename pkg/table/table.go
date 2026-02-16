package table

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/ggkhrmv/snow-go-sdk/pkg"
)

var (
	ErrInvalidTableName = errors.New("invalid table name")
	ErrInvalidSysID     = errors.New("invalid sys_id")
	ErrNilInput         = errors.New("input body is nil")
)

// Client is a Table API client bound to a specific table, returning records of type T.
type Client[T any] struct {
	r     pkg.Requester
	table string
}

func New[T any](r pkg.Requester, tableName string) *Client[T] {
	return &Client[T]{
		r:     r,
		table: strings.TrimSpace(tableName),
	}
}

// Convenience: default to dynamic records
func NewMap(r pkg.Requester, tableName string) *Client[map[string]any] {
	return New[map[string]any](r, tableName)
}

func (c *Client[T]) basePath() (string, error) {
	if c.table == "" {
		return "", ErrInvalidTableName
	}
	return path.Join("/api/now/table", c.table), nil
}

func (c *Client[T]) List(ctx context.Context, opts *ListOptions) ([]T, error) {
	base, err := c.basePath()
	if err != nil {
		return nil, err
	}

	q := url.Values{}
	if opts != nil {
		opts.apply(q)
	}

	req, err := c.r.NewRequest(ctx, http.MethodGet, base, q, nil)
	if err != nil {
		return nil, err
	}

	var out resultList[T]
	if err := c.r.Do(req, &out); err != nil {
		return nil, err
	}
	return out.Result, nil
}

func (c *Client[T]) Get(ctx context.Context, sysID string, opts *GetOptions) (T, error) {
	var zero T

	base, err := c.basePath()
	if err != nil {
		return zero, err
	}
	if strings.TrimSpace(sysID) == "" {
		return zero, ErrInvalidSysID
	}

	q := url.Values{}
	if opts != nil {
		opts.apply(q)
	}

	req, err := c.r.NewRequest(ctx, http.MethodGet, path.Join(base, sysID), q, nil)
	if err != nil {
		return zero, err
	}

	var out resultOne[T]
	if err := c.r.Do(req, &out); err != nil {
		return zero, err
	}
	return out.Result, nil
}

func (c *Client[T]) Create(ctx context.Context, in any, opts *WriteOptions) (T, error) {
	var zero T

	base, err := c.basePath()
	if err != nil {
		return zero, err
	}
	if in == nil {
		return zero, ErrNilInput
	}

	q := url.Values{}
	if opts != nil {
		opts.apply(q)
	}

	req, err := c.r.NewRequest(ctx, http.MethodPost, base, q, in)
	if err != nil {
		return zero, err
	}

	var out resultOne[T]
	if err := c.r.Do(req, &out); err != nil {
		return zero, err
	}
	return out.Result, nil
}

func (c *Client[T]) Update(ctx context.Context, sysID string, in any, opts *WriteOptions) (T, error) {
	var zero T

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
		opts.apply(q)
	}

	req, err := c.r.NewRequest(ctx, http.MethodPatch, path.Join(base, sysID), q, in)
	if err != nil {
		return zero, err
	}

	var out resultOne[T]
	if err := c.r.Do(req, &out); err != nil {
		return zero, err
	}
	return out.Result, nil
}

func (c *Client[T]) Replace(ctx context.Context, sysID string, in any, opts *WriteOptions) (T, error) {
	var zero T

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
		opts.apply(q)
	}

	req, err := c.r.NewRequest(ctx, http.MethodPut, path.Join(base, sysID), q, in)
	if err != nil {
		return zero, err
	}

	var out resultOne[T]
	if err := c.r.Do(req, &out); err != nil {
		return zero, err
	}
	return out.Result, nil
}

func (c *Client[T]) Delete(ctx context.Context, sysID string) error {
	base, err := c.basePath()
	if err != nil {
		return err
	}
	if strings.TrimSpace(sysID) == "" {
		return ErrInvalidSysID
	}

	req, err := c.r.NewRequest(ctx, http.MethodDelete, path.Join(base, sysID), nil, nil)
	if err != nil {
		return err
	}
	return c.r.Do(req, nil)
}