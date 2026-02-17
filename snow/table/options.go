package table

import (
	"net/url"
	"strconv"
	"strings"
)

type ListOptions struct {
	Query  string
	Fields []string
	Limit  *int
	Offset *int

	// Often "true", "false", or "all"
	DisplayValue *string

	ExcludeReferenceLink     *bool
	SuppressPaginationHeader *bool
}

type GetOptions struct {
	Fields               []string
	DisplayValue         *string
	ExcludeReferenceLink *bool
}

type WriteOptions struct {
	Fields               []string
	DisplayValue         *string
	ExcludeReferenceLink *bool
	InputDisplayValue    *bool
}

func (o *ListOptions) apply(q url.Values) {
	if o == nil {
		return
	}
	if strings.TrimSpace(o.Query) != "" {
		q.Set("sysparm_query", o.Query)
	}
	if len(o.Fields) > 0 {
		q.Set("sysparm_fields", strings.Join(o.Fields, ","))
	}
	if o.Limit != nil {
		q.Set("sysparm_limit", strconv.Itoa(*o.Limit))
	}
	if o.Offset != nil {
		q.Set("sysparm_offset", strconv.Itoa(*o.Offset))
	}
	if o.DisplayValue != nil {
		q.Set("sysparm_display_value", *o.DisplayValue)
	}
	if o.ExcludeReferenceLink != nil {
		q.Set("sysparm_exclude_reference_link", strconv.FormatBool(*o.ExcludeReferenceLink))
	}
	if o.SuppressPaginationHeader != nil {
		q.Set("sysparm_suppress_pagination_header", strconv.FormatBool(*o.SuppressPaginationHeader))
	}
}

func (o *GetOptions) apply(q url.Values) {
	if o == nil {
		return
	}
	if len(o.Fields) > 0 {
		q.Set("sysparm_fields", strings.Join(o.Fields, ","))
	}
	if o.DisplayValue != nil {
		q.Set("sysparm_display_value", *o.DisplayValue)
	}
	if o.ExcludeReferenceLink != nil {
		q.Set("sysparm_exclude_reference_link", strconv.FormatBool(*o.ExcludeReferenceLink))
	}
}

func (o *WriteOptions) apply(q url.Values) {
	if o == nil {
		return
	}
	if len(o.Fields) > 0 {
		q.Set("sysparm_fields", strings.Join(o.Fields, ","))
	}
	if o.DisplayValue != nil {
		q.Set("sysparm_display_value", *o.DisplayValue)
	}
	if o.ExcludeReferenceLink != nil {
		q.Set("sysparm_exclude_reference_link", strconv.FormatBool(*o.ExcludeReferenceLink))
	}
	if o.InputDisplayValue != nil {
		q.Set("sysparm_input_display_value", strconv.FormatBool(*o.InputDisplayValue))
	}
}

// pointer helpers (handy in callers)
func Int(v int) *int          { return &v }
func Bool(v bool) *bool       { return &v }
func String(v string) *string { return &v }
