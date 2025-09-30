package errors

import "net/http"

func BadRequest(message string) *AppError {
	return New(http.StatusBadRequest, message)
}
