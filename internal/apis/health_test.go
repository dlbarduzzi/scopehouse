package apis

import (
	"net/http"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	t.Parallel()

	s := apiTestScenario{
		name:   "health check",
		method: http.MethodGet,
		url:    "/api/v1/health",
		status: http.StatusOK,
		content: []string{
			`"status":200`,
			`"message":"API is healthy."`,
		},
	}

	s.test(t)
}
