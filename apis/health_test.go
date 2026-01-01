package apis

import (
	"net/http"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	t.Parallel()

	s := apiTestScenario{
		name:           "health check",
		url:            "/api/v1/health",
		method:         http.MethodGet,
		expectedStatus: 200,
		expectedContent: []string{
			`"status":200`,
			`"message":"API is healthy."`,
		},
	}

	s.Test(t)
}
