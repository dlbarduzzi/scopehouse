package apis

import (
	"testing"

	"github.com/dlbarduzzi/scopehouse/internal/tests"
)

func TestServiceHandler(t *testing.T) {
	app, err := tests.NewTestApp()
	if err != nil {
		t.Fatalf("failed to initialize test app instance; %v", err)
	}

	svc := newService(app)

	mux := svc.routes()
	if mux == nil {
		t.Fatal("expected service mux to be initialized")
	}

	handler := svc.handler(mux)
	if handler == nil {
		t.Fatal("expected service handler to be initialized")
	}
}
