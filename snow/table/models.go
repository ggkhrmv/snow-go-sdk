package table

// ServiceNow Table API wraps responses in {"result": ...}
type resultList[T any] struct {
	Result []T `json:"result"`
}

type resultOne[T any] struct {
	Result T `json:"result"`
}
