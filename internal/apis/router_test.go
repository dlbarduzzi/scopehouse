package apis

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dlbarduzzi/scopehouse/internal/core"
	"github.com/dlbarduzzi/scopehouse/internal/tests"
)

func TestNewRouter(t *testing.T) {
	app, err := tests.NewTestApp()
	if err != nil {
		t.Fatalf("failed to initialize test app instance - %v", err)
	}

	router := newRouter(app)

	// The calls to be made by each registered endpoint.
	calls := ""

	router.get("/a", func(*core.EventRequest) {
		calls += "a"
	})

	router.get("/b", func(*core.EventRequest) {
		calls += "b"
	})

	router.get("/a/b", func(*core.EventRequest) {
		calls += "a_b"
	})

	handler := router.handler()

	server := httptest.NewServer(handler)
	defer server.Close()

	client := server.Client()

	testCases := []struct {
		path   string
		calls  string
		method string
	}{
		{"/a", "a", http.MethodGet},
		{"/b", "b", http.MethodGet},
		{"/a/b", "a_b", http.MethodGet},
	}

	for _, tc := range testCases {
		t.Run(tc.method+"_"+tc.path, func(t *testing.T) {
			calls = "" // reset

			req, err := http.NewRequest(tc.method, server.URL+tc.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			_, err = client.Do(req)
			if err != nil {
				t.Fatal(err)
			}

			if calls != tc.calls {
				t.Fatalf("expected calls to be %q, got %q", tc.calls, calls)
			}
		})
	}
}
