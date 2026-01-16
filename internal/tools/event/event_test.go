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

func TestEventJson(t *testing.T) {
	testCases := []eventTestScenario{
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
		tc.eventFunc = func(e *Event) error {
			return e.Json(tc.data, tc.status)
		}
		tc.test(t)
	}
}

func TestEventText(t *testing.T) {
	testCases := []eventTestScenario{
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
		tc.eventFunc = func(e *Event) error {
			return e.Text(tc.status, tc.message)
		}
		tc.test(t)
	}
}

func TestEventStatus(t *testing.T) {
	testCases := []eventTestScenario{
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
		tc.eventFunc = func(e *Event) error {
			return e.Status(tc.status)
		}
		tc.test(t)
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

	message := "Something went wrong while processing your request."
	content := fmt.Sprintf(`{"status":500,"message":"%s"}`, message)

	if resStr != content {
		t.Fatalf("expected content to be \n%v \ngot \n%v", content, resStr)
	}

	if apiErr.Error() != message {
		t.Fatalf("expected error message to be %q, got %q", message, apiErr.Error())
	}
}

type eventTestScenario struct {
	name            string
	data            any
	status          int
	message         string
	headers         map[string]string
	eventFunc       func(e *Event) error
	expectedError   error
	expectedStatus  int
	expectedContent []string
	expectedHeaders map[string]string
}

func (s *eventTestScenario) test(t *testing.T) {
	t.Run(s.name, func(t *testing.T) {
		s._test(t)
	})
}

func (s *eventTestScenario) _test(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()

	for k, v := range s.headers {
		rec.Header().Add(k, v)
	}

	ev := &Event{
		Request:  req,
		Response: rec,
	}

	if s.eventFunc == nil {
		t.Fatal("you must provide an event to be tested")
	}

	err = s.eventFunc(ev)

	if s.expectedError != nil || err != nil {
		if !errors.Is(err, s.expectedError) {
			t.Fatalf("expected error %v, got %v", s.expectedError, err)
		}
	}

	result := rec.Result()

	if result.StatusCode != s.expectedStatus {
		t.Fatalf(
			"expected status code %d, got %d",
			s.expectedStatus, result.StatusCode,
		)
	}

	if err := result.Body.Close(); err != nil {
		t.Fatalf("failed to read response body - %v", err)
	}

	if len(s.expectedContent) == 0 {
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

		for _, content := range s.expectedContent {
			if !strings.Contains(body, content) {
				t.Errorf(
					"expected content %v in response body \n%v",
					content, body,
				)
			}
		}
	}

	for k, v := range s.expectedHeaders {
		if value := result.Header.Get(k); value != v {
			t.Fatalf("expected %q header to be %q, got %q", k, v, value)
		}
	}
}
