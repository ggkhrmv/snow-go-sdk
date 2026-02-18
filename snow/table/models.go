package table

// ListResponse represents a paginated list response from ServiceNow Table API
type ListResponse[T any] struct {
    Result []T               `json:"result"`
    Meta   *PaginationMeta   `json:"-"` // Populated from headers
}

// GetResponse represents a single record response
type GetResponse[T any] struct {
    Result T `json:"result"`
}

// WriteResponse represents a create/update/replace response
type WriteResponse[T any] struct {
    Result T `json:"result"`
}

// PaginationMeta contains pagination information from response headers
// These are returned unless sysparm_suppress_pagination_header=true
type PaginationMeta struct {
    // Link header with first, prev, next, last URLs
    First string
    Prev  string
    Next  string
    Last  string
    
    // Total count of records matching query
    TotalCount int
    
    // Current page info
    Limit  int
    Offset int
}

// ErrorResponse represents ServiceNow error responses
type ErrorResponse struct {
    Error  *ErrorDetail `json:"error,omitempty"`
    Status string       `json:"status,omitempty"`
}

type ErrorDetail struct {
    Message string `json:"message"`
    Detail  string `json:"detail,omitempty"`
}

// Internal envelope types matching the ServiceNow result wrapper.
type resultList[T any] struct {
    Result []T `json:"result"`
}

type resultOne[T any] struct {
    Result T `json:"result"`
}
