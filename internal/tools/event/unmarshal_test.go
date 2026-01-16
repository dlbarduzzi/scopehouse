package event

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEventUnmarshalSuccess(t *testing.T) {
	type person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	var data person

	d := `{"name":"test","age":1}`
	b := strings.NewReader(d)

	e := Event{
		Request:  httptest.NewRequest(http.MethodPost, "/", b),
		Response: httptest.NewRecorder(),
	}

	if err := e.Unmarshal(&data, nil); err != nil {
		t.Fatalf("expected error to be nil, got %v", err)
	}

	if data.Age != 1 || data.Name != "test" {
		t.Fatalf("expected result to contain `%v`, got `%+v`", d, data)
	}
}

func TestEventUnmarshalUnknownFieldsAllowedSuccess(t *testing.T) {
	type person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	var data person

	d := `{"name":"test","address":"111 Main St."}`
	b := strings.NewReader(d)

	e := Event{
		Request:  httptest.NewRequest(http.MethodPost, "/", b),
		Response: httptest.NewRecorder(),
	}

	opts := &UnmarshalOptions{
		DisallowUnknownFields: false,
	}

	if err := e.Unmarshal(&data, opts); err != nil {
		t.Fatalf("expected error to be nil, got %v", err)
	}

	if data.Age != 0 || data.Name != "test" {
		t.Fatalf("expected result to contain `%v`, got `%+v`", d, data)
	}
}

func TestEventWithWhiteSpaceSuccess(t *testing.T) {
	type person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	var data person

	// Valid single json object followed by whitespace.
	b := bytes.NewBufferString(`{"name":"test","age":1}   `)

	e := Event{
		Request:  httptest.NewRequest(http.MethodPost, "/", b),
		Response: httptest.NewRecorder(),
	}

	if err := e.Unmarshal(&data, nil); err != nil {
		t.Fatalf("expected error to be nil, got %v", err)
	}

	if data.Age != 1 || data.Name != "test" {
		t.Fatalf("expected result to contain `%v`, got `%+v`", b, data)
	}
}

func TestEventUnmarshalServerError(t *testing.T) {
	type person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	testCases := []struct {
		name    string
		req     *http.Request
		opts    *UnmarshalOptions
		data    any
		errStr  string
		message string
	}{
		{
			name:    "request is nil",
			req:     nil,
			opts:    nil,
			data:    &person{},
			errStr:  "request or request body cannot be nil",
			message: "Internal configuration error.",
		},
		{
			name:    "request body is nil",
			req:     &http.Request{Body: nil},
			opts:    nil,
			data:    &person{},
			errStr:  "request or request body cannot be nil",
			message: "Internal configuration error.",
		},
		{
			name:    "data destination is nil",
			req:     httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{}`)),
			opts:    nil,
			data:    nil,
			errStr:  "data destination cannot be nil",
			message: "Internal configuration error.",
		},
		{
			name: "invalid unmarshal",
			req:  httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{}`)),
			opts: nil,
			// pass non-pointer to trigger json.InvalidUnmarshalError
			data:    map[string]any{},
			errStr:  "json invalid unmarshal error; json: Unmarshal(non-pointer map[string]interface {})",
			message: "Internal configuration error.",
		},
	}

	for _, tc := range testCases {
		e := Event{
			Request:  tc.req,
			Response: httptest.NewRecorder(),
		}

		var err *UnmarshalError

		if tc.data == nil {
			err = e.Unmarshal(nil, tc.opts)
		} else {
			err = e.Unmarshal(tc.data, tc.opts)
		}
		if err == nil {
			t.Fatal("expected error not to be nil")
		}

		if err.IsClientError {
			t.Fatal("expected error to be `server` error")
		}

		if err.Err.Error() != tc.errStr {
			t.Fatalf("expected error to be `%v`, got `%v`", tc.errStr, err.Err.Error())
		}

		if err.Message != tc.message {
			t.Fatalf("expected error message to be %s, got %s", tc.message, err.Message)
		}
	}
}

type errUnmarshalCloser struct{}

func (e errUnmarshalCloser) Read(_ []byte) (int, error) {
	return 0, errors.New("test-error")
}

func (e errUnmarshalCloser) Close() error {
	return nil
}

func TestEventUnmarshalUnhandledDecodeError(t *testing.T) {
	// Use a reader that returns a non-standard error to hit default branch.
	b := io.NopCloser(errUnmarshalCloser{})

	e := Event{
		Request:  httptest.NewRequest(http.MethodPost, "/", b),
		Response: httptest.NewRecorder(),
	}

	e.Request.Body = errUnmarshalCloser{}

	opts := &UnmarshalOptions{
		MaxBodyBytes: 1024,
	}

	var data map[string]any

	err := e.Unmarshal(&data, opts)
	if err == nil {
		t.Fatal("expected error not to be nil")
	}

	if err.IsClientError {
		t.Fatal("expected error to be `server` error")
	}

	errStr := "json decoded error; test-error"
	if err.Err.Error() != errStr {
		t.Fatalf("expected error to be `%v`, got `%v`", errStr, err.Err.Error())
	}

	message := "Internal configuration error."
	if err.Message != message {
		t.Fatalf("expected error message to be %s, got %s", message, err.Message)
	}
}

func TestEventUnmarshalClientError(t *testing.T) {
	type person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	testRequest := func(body io.Reader) *http.Request {
		return httptest.NewRequest(http.MethodPost, "/", body)
	}

	testCases := []struct {
		name    string
		req     *http.Request
		opts    *UnmarshalOptions
		data    any
		errStr  string
		message string
	}{
		{
			name:    "body is empty",
			req:     testRequest(http.NoBody),
			opts:    nil,
			data:    map[string]any{},
			errStr:  "EOF",
			message: "Request body must not be empty.",
		},
		{
			name:    "malformed syntax",
			req:     testRequest(strings.NewReader(`{`)),
			opts:    nil,
			data:    map[string]any{},
			errStr:  "unexpected EOF",
			message: "Malformed json content in request body.",
		},
		{
			name:    "malformed syntax array",
			req:     testRequest(strings.NewReader(`[1, 2`)),
			opts:    nil,
			data:    map[string]any{},
			errStr:  "unexpected EOF",
			message: "Malformed json content in request body.",
		},
		{
			name: "syntax error with position",
			// Illegal character `!` - should trigger error.
			req:     testRequest(strings.NewReader(`{"name":!}`)),
			opts:    nil,
			data:    map[string]any{},
			errStr:  "invalid character '!' looking for beginning of value",
			message: "Malformed json content at position 9.",
		},
		{
			name:    "unmarshal type error with field",
			req:     testRequest(strings.NewReader(`{"age":"not-a-number"}`)),
			opts:    nil,
			data:    &person{},
			errStr:  "json: cannot unmarshal string into Go struct field person.age of type int",
			message: "Invalid value type for field \"age\".",
		},
		{
			name:    "unknown field disallowed",
			req:     testRequest(strings.NewReader(`{"name":"test","extra":123}`)),
			opts:    &UnmarshalOptions{DisallowUnknownFields: true},
			data:    &person{},
			errStr:  "json: unknown field \"extra\"",
			message: "Unknown field '\"extra\"' in request body.",
		},
		{
			name:    "multiple objects",
			req:     testRequest(strings.NewReader(`{}{}`)),
			opts:    nil,
			data:    &person{},
			errStr:  "request body must contain a single json object",
			message: "Request body must contain a single json object.",
		},
		{
			name:    "body too large",
			req:     testRequest(strings.NewReader(`{"x":` + strings.Repeat("1", 100) + `}`)),
			opts:    &UnmarshalOptions{MaxBodyBytes: 10},
			data:    &person{},
			errStr:  "http: request body too large",
			message: "Request body must not be larger than 10 bytes.",
		},
	}

	for _, tc := range testCases {
		e := Event{
			Request:  tc.req,
			Response: httptest.NewRecorder(),
		}

		var err *UnmarshalError

		if tc.data == nil {
			err = e.Unmarshal(nil, tc.opts)
		} else {
			err = e.Unmarshal(tc.data, tc.opts)
		}
		if err == nil {
			t.Fatal("expected error not to be nil")
		}

		if !err.IsClientError {
			t.Fatal("expected error to be `client` error")
		}

		if err.Err.Error() != tc.errStr {
			t.Fatalf("expected error to be `%v`, got `%v`", tc.errStr, err.Err.Error())
		}

		if err.Message != tc.message {
			t.Fatalf("expected error message to be %s, got %s", tc.message, err.Message)
		}
	}
}
