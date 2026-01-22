package structs

// Standard envelope for all JSON responses
// @x-extension-openapi-ignore
type SuccessResponse[T any] struct {
	Status  bool   `json:"status"`
	Message string `json:"message,omitempty"`
	Data    T      `json:"data,omitempty"`
}

// ErrorResponse handles standard error structures
type ErrorResponse struct {
	Status bool   `json:"status"`
	Error  string `json:"error"`
	Code   int    `json:"code,omitempty"`
}
