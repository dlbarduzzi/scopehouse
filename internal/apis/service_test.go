package apis

import (
	"errors"
	"fmt"
	"testing"

	"github.com/dlbarduzzi/scopehouse/internal/tests"
)

func TestNewService(t *testing.T) {
	svc, err := newTestService(t)
	if err != nil {
		t.Fatal(err)
	}

	if svc == nil {
		t.Fatal("expected service to be initialized")
	}

	if svc.app == nil {
		t.Fatal("expected service app to be initialized")
	}
}

func newTestService(t *testing.T) (*service, error) {
	t.Helper()

	app, err := tests.NewTestApp()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize test app instance; %v", err)
	}

	svc := newService(app)
	if svc == nil {
		return nil, errors.New("expected service to be initialized")
	}

	return svc, nil
}
