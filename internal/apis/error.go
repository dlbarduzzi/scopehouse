package apis

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/dlbarduzzi/scopehouse/internal/tools/event"
	"github.com/dlbarduzzi/scopehouse/internal/tools/inflector"
)

func (s *service) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	s.app.Logger().Error("internal server error",
		slog.Any("error", err),
		slog.String("method", r.Method),
		slog.String("request", r.RequestURI),
	)

	resp := newInternalServerError("")

	if err := event.WriteJson(w, resp, resp.Status); err != nil {
		_ = event.WriteStatus(w, http.StatusInternalServerError)
		return
	}
}

type apiError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// Error makes it compatible with the `error` interface.
func (e *apiError) Error() string {
	return e.Message
}

func newApiError(status int, message string) *apiError {
	msg := strings.TrimSpace(message)
	if msg == "" {
		msg = http.StatusText(status)
	}

	return &apiError{
		Status:  status,
		Message: inflector.FormatSentence(msg),
	}
}

func newInternalServerError(message string) *apiError {
	msg := strings.TrimSpace(message)
	if msg == "" {
		msg = "Something went wrong while processing your request."
	}

	return newApiError(http.StatusInternalServerError, msg)
}
