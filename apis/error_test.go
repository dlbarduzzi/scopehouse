package apis

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dlbarduzzi/scopehouse/core"
	"github.com/dlbarduzzi/scopehouse/tests"
	"github.com/dlbarduzzi/scopehouse/tools/event"
)

func TestInternalServerError(t *testing.T) {
	app, err := tests.NewTestApp()
	if err != nil {
		t.Fatalf("failed to initialize test app instance - %v", err)
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://127.0.0.1", nil)

	e := &core.EventRequest{
		App: app,
		Event: event.Event{
			Request:  req,
			Response: rec,
		},
	}

	// No results to assert. This test only verifies that the function executes without error.
	// For example, removing `Event` from `core.EventRequest` would cause this test to fail.
	internalServerError(e, nil)
}
