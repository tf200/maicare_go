package handler

type Response[T any] struct {
	Success bool   `json:"success"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

func successResponse[T any](data T, message string) Response[T] {
	return Response[T]{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func errorResponse(message error) Response[struct{}] {
	return Response[struct{}]{
		Success: false,
		Message: message.Error(),
	}
}
