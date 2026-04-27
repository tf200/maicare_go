package domain

// ListResult is a generic wrapper for paginated list results.
type ListResult[T any] struct {
	Items      []T
	TotalCount int64
}
