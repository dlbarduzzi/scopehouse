package apis

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/dlbarduzzi/scopehouse/internal/core"
)

func internalServerError(e *core.EventRequest, err error) {
	e.App.Logger().Error("internal server error",
		slog.String("code", "INTERNAL_SERVER_ERROR"),
		slog.String("error", fmt.Sprintf("%v", err)),
		slog.String("method", e.Request.Method),
		slog.String("request", e.Request.RequestURI),
	)

	resp := e.InternalServerError("")

	if err := e.Json(resp, resp.Status); err != nil {
		_ = e.Status(http.StatusInternalServerError)
		return
	}
}
