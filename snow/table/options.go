package table

import (
    "errors"
    "net/url"
    "strconv"
    "strings"
)

var (
    ErrMutuallyExclusiveFilters = errors.New("Query and Filters are mutually exclusive")
    ErrInvalidDisplayValue      = errors.New("DisplayValue must be 'true', 'false', or 'all'")
)

// ListOptions provides ergonomic API for list queries
type ListOptions struct {
    // Filtering (mutually exclusive)
    Query   string            // Encoded query string
    Filters map[string]string // Name-value pairs: {"active": "true", "state": "closed"}
    
    // Field selection
    Fields []string // List of field names to return
    
    // Pagination
    Limit  *int
    Offset *int
    
    // Display options
    DisplayValue         *DisplayValueOption // Use constants: DisplayValueTrue, DisplayValueFalse, DisplayValueAll
    ExcludeReferenceLink *bool
    
    // Domain control (requires special permissions)
    QueryNoDomain *bool
    
    // Response formatting
    SuppressPaginationHeader *bool
}

type GetOptions struct {
    Fields               []string
    DisplayValue         *DisplayValueOption
    ExcludeReferenceLink *bool
}

type WriteOptions struct {
    Fields               []string
    DisplayValue         *DisplayValueOption
    ExcludeReferenceLink *bool
    InputDisplayValue    *bool // For write operations: interpret input as display values
}

type DeleteOptions struct {
    QueryNoDomain *bool // Restrict to user's domains
}

// Validate checks for configuration errors
func (o *ListOptions) Validate() error {
    if o == nil {
        return nil
    }
    
    // Query and Filters are mutually exclusive
    if strings.TrimSpace(o.Query) != "" && len(o.Filters) > 0 {
        return ErrMutuallyExclusiveFilters
    }
    
    // Validate DisplayValue if set
    if o.DisplayValue != nil {
        if err := o.DisplayValue.Validate(); err != nil {
            return err
        }
    }
    
    return nil
}

func (o *GetOptions) Validate() error {
    if o != nil && o.DisplayValue != nil {
        return o.DisplayValue.Validate()
    }
    return nil
}

func (o *WriteOptions) Validate() error {
    if o != nil && o.DisplayValue != nil {
        return o.DisplayValue.Validate()
    }
    return nil
}

// apply converts SDK options to URL query parameters
func (o *ListOptions) apply(q url.Values) error {
    if o == nil {
        return nil
    }
    
    if err := o.Validate(); err != nil {
        return err
    }
    
    // Filtering: Query takes precedence over Filters
    if strings.TrimSpace(o.Query) != "" {
        q.Set("sysparm_query", o.Query)
    } else if len(o.Filters) > 0 {
        // Add name-value pairs directly
        for k, v := range o.Filters {
            q.Set(k, v)
        }
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
        q.Set("sysparm_display_value", string(*o.DisplayValue))
    }
    if o.ExcludeReferenceLink != nil {
        q.Set("sysparm_exclude_reference_link", strconv.FormatBool(*o.ExcludeReferenceLink))
    }
    if o.QueryNoDomain != nil {
        q.Set("sysparm_query_no_domain", strconv.FormatBool(*o.QueryNoDomain))
    }
    if o.SuppressPaginationHeader != nil {
        q.Set("sysparm_suppress_pagination_header", strconv.FormatBool(*o.SuppressPaginationHeader))
    }
    
    return nil
}

func (o *GetOptions) apply(q url.Values) error {
    if o == nil {
        return nil
    }
    
    if err := o.Validate(); err != nil {
        return err
    }
    
    if len(o.Fields) > 0 {
        q.Set("sysparm_fields", strings.Join(o.Fields, ","))
    }
    if o.DisplayValue != nil {
        q.Set("sysparm_display_value", string(*o.DisplayValue))
    }
    if o.ExcludeReferenceLink != nil {
        q.Set("sysparm_exclude_reference_link", strconv.FormatBool(*o.ExcludeReferenceLink))
    }
    
    return nil
}

func (o *WriteOptions) apply(q url.Values) error {
    if o == nil {
        return nil
    }
    
    if err := o.Validate(); err != nil {
        return err
    }
    
    if len(o.Fields) > 0 {
        q.Set("sysparm_fields", strings.Join(o.Fields, ","))
    }
    if o.DisplayValue != nil {
        q.Set("sysparm_display_value", string(*o.DisplayValue))
    }
    if o.ExcludeReferenceLink != nil {
        q.Set("sysparm_exclude_reference_link", strconv.FormatBool(*o.ExcludeReferenceLink))
    }
    if o.InputDisplayValue != nil {
        q.Set("sysparm_input_display_value", strconv.FormatBool(*o.InputDisplayValue))
    }
    
    return nil
}

func (o *DeleteOptions) apply(q url.Values) error {
    if o == nil {
        return nil
    }
    if o.QueryNoDomain != nil {
        q.Set("sysparm_query_no_domain", strconv.FormatBool(*o.QueryNoDomain))
    }
    return nil
}

// Pointer helpers
func Int(v int) *int       { return &v }
func Bool(v bool) *bool    { return &v }
func String(v string) *string { return &v }
func DisplayValue(v DisplayValueOption) *DisplayValueOption { return &v }