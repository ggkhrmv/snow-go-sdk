package table

// API request/response types that match ServiceNow wire format exactly

// QueryParams represents all possible query parameters for Table API
type QueryParams struct {
	// Filtering
	Query   string            `url:"sysparm_query,omitempty"`
	Filters map[string]string `url:"-"` // Handled separately as name-value pairs

	// Field selection
	Fields string `url:"sysparm_fields,omitempty"` // Comma-separated

	// Pagination
	Limit  int `url:"sysparm_limit,omitempty"`
	Offset int `url:"sysparm_offset,omitempty"`

	// Display options
	DisplayValue         string `url:"sysparm_display_value,omitempty"` // "true", "false", "all"
	ExcludeReferenceLink bool   `url:"sysparm_exclude_reference_link,omitempty"`
	InputDisplayValue    bool   `url:"sysparm_input_display_value,omitempty"`

	// Domain control
	QueryNoDomain bool `url:"sysparm_query_no_domain,omitempty"`

	// Response formatting
	SuppressPaginationHeader bool `url:"sysparm_suppress_pagination_header,omitempty"`
}

// DisplayValueOption represents valid display_value parameter values
type DisplayValueOption string

const (
	DisplayValueTrue  DisplayValueOption = "true"
	DisplayValueFalse DisplayValueOption = "false"
	DisplayValueAll   DisplayValueOption = "all"
)

func (d DisplayValueOption) Validate() error {
	switch d {
	case DisplayValueTrue, DisplayValueFalse, DisplayValueAll, "":
		return nil
	default:
		return ErrInvalidDisplayValue
	}
}
