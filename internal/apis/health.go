package apis

import (
	"fmt"
	"log/slog"
	"net/http"
)

func (s *service) healthCheck(w http.ResponseWriter, r *http.Request) {
	resp := struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}{
		Status:  http.StatusOK,
		Message: "API is healthy.",
	}

	s.app.Logger().Info("health called", slog.String("url", r.RequestURI))
	_, _ = fmt.Fprintf(w, "%d - %s", resp.Status, resp.Message)
}
