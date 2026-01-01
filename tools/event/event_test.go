package event

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type testCase struct {
	name            string
	data            any
	status          int
	message         string
	headers         map[string]string
	expectedError   error
	expectedStatus  int
	expectedContent []string
	expectedHeaders map[string]string
}

func TestEventJson(t *testing.T) {
	testCases := []testCase{
		{
			name:            "no header",
			data:            map[string]any{"foo": "bar", "num": 123},
			status:          200,
			headers:         nil,
			expectedError:   nil,
			expectedStatus:  200,
			expectedContent: []string{`"foo":"bar"`, `"num":123`},
			expectedHeaders: map[string]string{"content-type": "application/json"},
		},
		{
			name:            "custom header",
			data:            map[string]any{"foo": "bar", "num": 123},
			status:          200,
			headers:         map[string]string{"content-type": "application/test"},
			expectedError:   nil,
			expectedStatus:  200,
			expectedContent: []string{`"foo":"bar"`, `"num":123`},
			expectedHeaders: map[string]string{"content-type": "application/json"},
		},
		{
			name:            "status 400",
			data:            map[string]any{"foo": "bar", "num": 123},
			status:          400,
			headers:         map[string]string{"content-type": "application/test"},
			expectedError:   nil,
			expectedStatus:  400,
			expectedContent: []string{`"foo":"bar"`, `"num":123`},
			expectedHeaders: map[string]string{"content-type": "application/json"},
		},
	}

	for _, tc := range testCases {
		testEvent(t, tc, func(e *Event) error {
			return e.Json(tc.data, tc.status)
		})
	}
}

func TestEventText(t *testing.T) {
	testCases := []testCase{
		{
			name:            "status 200",
			data:            nil,
			status:          200,
			message:         "hello world",
			headers:         nil,
			expectedError:   nil,
			expectedStatus:  200,
			expectedContent: []string{"hello world"},
			expectedHeaders: map[string]string{"content-type": "text/plain; charset=utf-8"},
		},
		{
			name:            "status 400",
			data:            nil,
			status:          400,
			message:         "hello world",
			headers:         nil,
			expectedError:   nil,
			expectedStatus:  400,
			expectedContent: []string{"hello world"},
			expectedHeaders: map[string]string{"content-type": "text/plain; charset=utf-8"},
		},
		{
			name:            "empty message",
			data:            nil,
			status:          500,
			message:         "",
			headers:         nil,
			expectedError:   nil,
			expectedStatus:  500,
			expectedContent: []string{"Internal Server Error"},
			expectedHeaders: map[string]string{"content-type": "text/plain; charset=utf-8"},
		},
	}

	for _, tc := range testCases {
		testEvent(t, tc, func(e *Event) error {
			return e.Text(tc.status, tc.message)
		})
	}
}

func TestEventStatus(t *testing.T) {
	testCases := []testCase{
		{
			name:            "status 200",
			data:            nil,
			status:          200,
			headers:         nil,
			expectedError:   nil,
			expectedStatus:  200,
			expectedContent: []string{"OK"},
			expectedHeaders: map[string]string{"content-type": "text/plain; charset=utf-8"},
		},
		{
			name:            "status 400",
			data:            nil,
			status:          400,
			headers:         nil,
			expectedError:   nil,
			expectedStatus:  400,
			expectedContent: []string{"Bad Request"},
			expectedHeaders: map[string]string{"content-type": "text/plain; charset=utf-8"},
		},
		{
			name:            "status 500",
			data:            nil,
			status:          500,
			headers:         nil,
			expectedError:   nil,
			expectedStatus:  500,
			expectedContent: []string{"Internal Server Error"},
			expectedHeaders: map[string]string{"content-type": "text/plain; charset=utf-8"},
		},
	}

	for _, tc := range testCases {
		testEvent(t, tc, func(e *Event) error {
			return e.Status(tc.status)
		})
	}
}

func TestEventInternalServerError(t *testing.T) {
	t.Parallel()

	ev := Event{}
	apiErr := ev.InternalServerError("")

	res, err := json.Marshal(apiErr)
	if err != nil {
		t.Fatal(err)
	}

	resStr := string(res)

	message := "Something went wrong while processing this request."
	content := fmt.Sprintf(`{"status":500,"message":"%s"}`, message)

	if resStr != content {
		t.Fatalf("expected content to be \n%v \ngot \n%v", content, resStr)
	}

	if apiErr.Error() != message {
		t.Fatalf("expected error message to be %q, got %q", message, apiErr.Error())
	}
}

func TestErrorHandler(t *testing.T) {
	t.Parallel()

	tc := testCase{
		name: "no ApiError defined",
		expectedContent: []string{
			`"status":400`,
			`"message":"Bad Request."`,
		},
		expectedStatus: 400,
		expectedHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	testEvent(t, tc, func(e *Event) error {
		return ErrorHandler(e.Response, e.Request, errors.New("test error"))
	})

	tc = testCase{
		name: "with ApiError defined",
		expectedContent: []string{
			`"status":500`,
			`"message":"Internal Server Error."`,
		},
		expectedStatus: 500,
		expectedHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	testEvent(t, tc, func(e *Event) error {
		err := NewApiError(tc.expectedStatus, "")
		return ErrorHandler(e.Response, e.Request, err)
	})
}

func testEvent(t *testing.T, tc testCase, fn func(e *Event) error) {
	t.Helper()

	t.Run(tc.name, func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		rec := httptest.NewRecorder()

		ev := &Event{
			Request:  req,
			Response: rec,
		}

		for k, v := range tc.headers {
			ev.Response.Header().Add(k, v)
		}

		err = fn(ev)
		if tc.expectedError != nil || err != nil {
			if !errors.Is(err, tc.expectedError) {
				t.Fatalf("expected error %v, got %v", tc.expectedError, err)
			}
		}

		result := rec.Result()

		if result.StatusCode != tc.expectedStatus {
			t.Fatalf(
				"expected status code %d, got %d",
				tc.expectedStatus, result.StatusCode,
			)
		}

		if err := result.Body.Close(); err != nil {
			t.Fatalf("failed to read response body - %v", err)
		}

		if len(tc.expectedContent) == 0 {
			if len(rec.Body.Bytes()) != 0 {
				t.Fatalf(
					"expected empty content, got \n%v",
					rec.Body.String(),
				)
			}
		} else {
			var body string
			buf := new(bytes.Buffer)

			err := json.Compact(buf, rec.Body.Bytes())
			if err != nil {
				// Not a json payload.
				body = rec.Body.String()
			} else {
				// A valid json payload.
				body = buf.String()
			}

			for _, content := range tc.expectedContent {
				if !strings.Contains(body, content) {
					t.Errorf(
						"expected content %v in response body \n%v",
						content, body,
					)
				}
			}
		}

		for k, v := range tc.expectedHeaders {
			if value := result.Header.Get(k); value != v {
				t.Fatalf("expected %q header to be %q, got %q", k, v, value)
			}
		}
	})
}
