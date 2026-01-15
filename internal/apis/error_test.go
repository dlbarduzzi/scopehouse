package apis

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type errorTestScenario struct {
	name    string
	status  int
	content []string

	// errorTestFunc runs error functions test cases.
	errorTestFunc func(http.ResponseWriter, *http.Request)
}

func (s *errorTestScenario) test(t *testing.T) {
	t.Run(s.name, func(t *testing.T) {
		s._test(t)
	})
}

func (s *errorTestScenario) _test(t *testing.T) {
	t.Helper()
	t.Parallel()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	if s.errorTestFunc == nil {
		t.Fatal("you must provide a function to be tested")
	}

	s.errorTestFunc(rec, req)
	res := rec.Result()

	if res.StatusCode != s.status {
		t.Fatalf("expected status code to be %d, got %d", s.status, res.StatusCode)
	}

	testBodyContent(t, rec, s.content)
}

func TestInternalServerError(t *testing.T) {
	svc, err := newTestService(t)
	if err != nil {
		t.Fatal(err)
	}

	s := errorTestScenario{
		name:   "internal server error",
		status: http.StatusInternalServerError,
		content: []string{
			`"status":500`,
			`"message":"Something went wrong while processing your request."`,
		},
		errorTestFunc: func(w http.ResponseWriter, r *http.Request) {
			svc.internalServerError(w, r, errors.New("test error"))
		},
	}

	s.test(t)
}

func TestNewApiError(t *testing.T) {
	t.Parallel()

	testCases := []apiErrorTestScenario{
		{
			name:    "empty message",
			apiErr:  newApiError(400, ""),
			content: []string{`"status":400`, `"message":"Bad Request."`},
			message: "Bad Request.",
		},
		{
			name:    "custom message",
			apiErr:  newApiError(400, "Test - Bad Request."),
			content: []string{`"status":400`, `"message":"Test - Bad Request."`},
			message: "Test - Bad Request.",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.test(t)
		})
	}
}

func TestNewInternalServerError(t *testing.T) {
	t.Parallel()

	testCases := []apiErrorTestScenario{
		{
			name:   "empty message",
			apiErr: newInternalServerError(""),
			content: []string{
				`"status":500`,
				`"message":"Something went wrong while processing your request."`,
			},
			message: "Something went wrong while processing your request.",
		},
		{
			name:   "custom message",
			apiErr: newInternalServerError("Test - NewInternalServerError"),
			content: []string{
				`"status":500`,
				`"message":"Test - NewInternalServerError."`,
			},
			message: "Test - NewInternalServerError.",
		},
	}

	for _, tc := range testCases {
		tc.test(t)
	}
}

type apiErrorTestScenario struct {
	name    string
	apiErr  *apiError
	content []string
	message string
}

func (s *apiErrorTestScenario) test(t *testing.T) {
	t.Run(s.name, func(t *testing.T) {
		s._test(t)
	})
}

func (s *apiErrorTestScenario) _test(t *testing.T) {
	t.Helper()
	t.Parallel()

	e := s.apiErr

	res, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}

	resStr := string(res)

	for _, content := range s.content {
		if !strings.Contains(resStr, content) {
			t.Fatalf(
				"expected content `%v` in response \n%v",
				content, s,
			)
		}
	}

	if e.Error() != s.message {
		t.Fatalf("expected error message to be %q, got %q", s.message, e.Error())
	}
}
