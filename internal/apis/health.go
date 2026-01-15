package apis

import (
	"net/http"

	"github.com/dlbarduzzi/scopehouse/internal/tools/event"
)

func (s *service) healthCheck(w http.ResponseWriter, r *http.Request) {
	resp := struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}{
		Status:  http.StatusOK,
		Message: "API is healthy.",
	}

	if err := event.WriteJson(w, resp, resp.Status); err != nil {
		s.internalServerError(w, r, err)
		return
	}
}
