package middleware

type responseEnvelope[T any] struct {
	Success bool   `json:"success"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

func errorResponse(err error) responseEnvelope[struct{}] {
	return responseEnvelope[struct{}]{
		Success: false,
		Message: err.Error(),
	}
}
