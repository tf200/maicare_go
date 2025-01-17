package api

// Response represents a generic API response
type Response[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

// SuccessResponse creates a successful response with data
func SuccessResponse[T any](data T, message string) Response[T] {
	return Response[T]{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// errorResponse represents the standard error response
func errorResponse(message error) Response[struct{}] {
	return Response[struct{}]{
		Success: false,
		Message: message.Error(),
	}
}
