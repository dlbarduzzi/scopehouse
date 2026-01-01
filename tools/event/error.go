package event

import (
	"errors"
	"net/http"
	"strings"

	"github.com/dlbarduzzi/scopehouse/tools/inflector"
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
	message = strings.TrimSpace(message)
	if message == "" {
		message = http.StatusText(status)
	}

	return &ApiError{
		Status:  status,
		Message: inflector.FormatSentence(message),
	}
}

func GenericApiError() *ApiError {
	return NewApiError(http.StatusBadRequest, "")
}

func NewInternalServerError(message string) *ApiError {
	message = strings.TrimSpace(message)
	if message == "" {
		message = "Something went wrong while processing this request."
	}

	return NewApiError(http.StatusInternalServerError, message)
}

func ErrorHandler(w http.ResponseWriter, r *http.Request, err error) error {
	if err == nil {
		return nil
	}

	var apiErr *ApiError

	if !errors.As(err, &apiErr) {
		apiErr = GenericApiError()
	}

	e := Event{
		Request:  r,
		Response: w,
	}

	if err := e.Json(apiErr, apiErr.Status); err != nil {
		return e.Status(http.StatusInternalServerError)
	}

	return nil
}
