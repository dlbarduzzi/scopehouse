package event

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Event struct {
	Request  *http.Request
	Response http.ResponseWriter
}

func (e *Event) Json(data any, status int) error {
	res, err := json.Marshal(data)
	if err != nil {
		return nil
	}

	res = append(res, '\n')

	e.Response.Header().Set("Content-Type", "application/json")
	e.Response.WriteHeader(status)

	if _, err := e.Response.Write(res); err != nil {
		return err
	}

	return nil
}

func (e *Event) Text(status int, message string) error {
	message = strings.TrimSpace(message)
	if message == "" {
		message = http.StatusText(status)
	}

	e.Response.Header().Set("Content-Type", "text/plain; charset=utf-8")
	e.Response.WriteHeader(status)

	if _, err := e.Response.Write([]byte(message)); err != nil {
		return err
	}

	return nil
}

func (e *Event) Status(status int) error {
	return e.Text(status, "")
}

func (e *Event) InternalServerError(message string) *ApiError {
	return NewInternalServerError(message)
}
