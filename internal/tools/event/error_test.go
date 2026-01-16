package event

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestNewApiError(t *testing.T) {
	t.Parallel()

	testCases := []apiErrorTestScenario{
		{
			name:    "empty message",
			apiErr:  NewApiError(400, ""),
			content: []string{`"status":400`, `"message":"Bad Request."`},
			message: "Bad Request.",
		},
		{
			name:    "custom message",
			apiErr:  NewApiError(400, "Test - Bad Request."),
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
			apiErr: NewInternalServerError(""),
			content: []string{
				`"status":500`,
				`"message":"Something went wrong while processing your request."`,
			},
			message: "Something went wrong while processing your request.",
		},
		{
			name:   "custom message",
			apiErr: NewInternalServerError("Test - NewInternalServerError"),
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
	apiErr  *ApiError
	content []string
	message string
}

func (s *apiErrorTestScenario) test(t *testing.T) {
	t.Run(s.name, func(t *testing.T) {
		s._test(t)
	})
}

func (s *apiErrorTestScenario) _test(t *testing.T) {
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
