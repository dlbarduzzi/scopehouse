package event

import (
	"net/http"
	"strings"

	"github.com/dlbarduzzi/scopehouse/internal/tools/inflector"
)

type ApiError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// Error makes it compatible with the `error` interface.
func (e *ApiError) Error() string {
	return e.Message
}

func NewApiError(status int, message string) *ApiError {
	msg := strings.TrimSpace(message)
	if msg == "" {
		msg = http.StatusText(status)
	}

	return &ApiError{
		Status:  status,
		Message: inflector.FormatSentence(msg),
	}
}

func NewInternalServerError(message string) *ApiError {
	msg := strings.TrimSpace(message)
	if msg == "" {
		msg = "Something went wrong while processing your request."
	}

	return NewApiError(http.StatusInternalServerError, msg)
}
