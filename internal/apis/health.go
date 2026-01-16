package apis

import (
	"net/http"

	"github.com/dlbarduzzi/scopehouse/internal/core"
)

func healthCheck(e *core.EventRequest) {
	resp := struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}{
		Status:  http.StatusOK,
		Message: "API is healthy.",
	}

	if err := e.Json(resp, resp.Status); err != nil {
		internalServerError(e, err)
		return
	}
}
