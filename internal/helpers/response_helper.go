package helpers

import "kora_1/internal/structs"

func NewSuccess[T any](data T, message string) structs.SuccessResponse[T] {
	return structs.SuccessResponse[T]{
		Status:  true,
		Message: message,
		Data:    data,
	}
}

func NewError(err string, code int) structs.ErrorResponse {
	return structs.ErrorResponse{
		Status: false,
		Error:  err,
		Code:   code,
	}
}
