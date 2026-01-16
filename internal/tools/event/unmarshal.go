package event

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const DefaultMaxBodyBytes = 1 << 20 // 1MB

type UnmarshalOptions struct {
	MaxBodyBytes          int64
	DisallowUnknownFields bool
}

var defaultUnmarshalOptions = &UnmarshalOptions{
	MaxBodyBytes:          DefaultMaxBodyBytes,
	DisallowUnknownFields: false,
}

type UnmarshalError struct {
	Err           error
	Message       string
	IsClientError bool
}

func (e *Event) Unmarshal(data any, opts *UnmarshalOptions) *UnmarshalError {
	if e.Request == nil || e.Request.Body == nil {
		err := errors.New("request or request body cannot be nil")
		return unmarshalServerError(err)
	}

	if data == nil {
		err := errors.New("data destination cannot be nil")
		return unmarshalServerError(err)
	}

	if opts == nil {
		opts = defaultUnmarshalOptions
	}

	if opts.MaxBodyBytes <= 0 {
		opts.MaxBodyBytes = DefaultMaxBodyBytes
	}

	e.Request.Body = http.MaxBytesReader(e.Response, e.Request.Body, opts.MaxBodyBytes)
	defer func() {
		if err := e.Request.Body.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "[error] request body close; %s\n", err)
		}
	}()

	dec := json.NewDecoder(e.Request.Body)
	if opts.DisallowUnknownFields {
		dec.DisallowUnknownFields()
	}

	if err := dec.Decode(data); err != nil {
		return decodedError(err, opts.MaxBodyBytes)
	}

	if dec.More() {
		msg := "Request body must contain a single json object."
		err := errors.New("request body must contain a single json object")
		return unmarshalClientError(err, msg)
	}

	return nil
}

func unmarshalServerError(err error) *UnmarshalError {
	return &UnmarshalError{
		Err:           err,
		Message:       "Internal configuration error.",
		IsClientError: false,
	}
}

func unmarshalClientError(err error, message string) *UnmarshalError {
	return &UnmarshalError{
		Err:           err,
		Message:       message,
		IsClientError: true,
	}
}

func decodedError(err error, maxBodyBytes int64) *UnmarshalError {
	var se *json.SyntaxError
	if errors.As(err, &se) {
		msg := fmt.Sprintf("Malformed json content at position %d.", se.Offset)
		return unmarshalClientError(err, msg)
	}

	if errors.Is(err, io.ErrUnexpectedEOF) {
		return unmarshalClientError(err, "Malformed json content in request body.")
	}

	var ute *json.UnmarshalTypeError
	if errors.As(err, &ute) {
		msg := fmt.Sprintf("Invalid value type at character %q.", ute.Offset)
		if ute.Field != "" {
			msg = fmt.Sprintf("Invalid value type for field %q.", ute.Field)
		}
		return unmarshalClientError(err, msg)
	}

	if errors.Is(err, io.EOF) {
		return unmarshalClientError(err, "Request body must not be empty.")
	}

	if strings.HasPrefix(err.Error(), "json: unknown field") {
		field := strings.TrimPrefix(err.Error(), "json: unknown field ")
		msg := fmt.Sprintf("Unknown field '%s' in request body.", field)
		return unmarshalClientError(err, msg)
	}

	var mbe *http.MaxBytesError
	if errors.As(err, &mbe) {
		msg := fmt.Sprintf("Request body must not be larger than %d bytes.", maxBodyBytes)
		return unmarshalClientError(err, msg)
	}

	var iue *json.InvalidUnmarshalError
	if errors.As(err, &iue) {
		e := fmt.Errorf("json invalid unmarshal error; %v", err)
		return unmarshalServerError(e)
	}

	e := fmt.Errorf("json decoded error; %v", err)
	return unmarshalServerError(e)
}
