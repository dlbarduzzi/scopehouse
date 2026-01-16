package apis

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dlbarduzzi/scopehouse/internal/core"
	"github.com/dlbarduzzi/scopehouse/internal/tests"
	"github.com/dlbarduzzi/scopehouse/internal/tools/event"
)

func TestInternalServerError(t *testing.T) {
	s := errorTestScenario{
		name:   "internal server error",
		status: http.StatusInternalServerError,
		content: []string{
			`"status":500`,
			`"message":"Something went wrong while processing your request."`,
		},
		errorTestFunc: func(e *core.EventRequest) {
			internalServerError(e, errors.New("test error"))
		},
	}

	s.test(t)
}

type errorTestScenario struct {
	name    string
	status  int
	content []string

	// errorTestFunc runs error functions test cases.
	errorTestFunc func(*core.EventRequest)
}

func (s *errorTestScenario) test(t *testing.T) {
	t.Run(s.name, func(t *testing.T) {
		s._test(t)
	})
}

func (s *errorTestScenario) _test(t *testing.T) {
	app, err := tests.NewTestApp()
	if err != nil {
		t.Fatalf("failed to initialize test app instance; %v", err)
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	if s.errorTestFunc == nil {
		t.Fatal("you must provide an error function to be tested")
	}

	ev := &core.EventRequest{
		App: app,
		Event: event.Event{
			Request:  req,
			Response: rec,
		},
	}

	s.errorTestFunc(ev)
	res := rec.Result()

	if res.StatusCode != s.status {
		t.Fatalf("expected status code to be %d, got %d", s.status, res.StatusCode)
	}

	testBodyContent(t, rec, s.content)
}
