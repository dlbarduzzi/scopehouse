package apis

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime"
	"time"

	"github.com/dlbarduzzi/scopehouse/core"
	"github.com/dlbarduzzi/scopehouse/tools/event"
)

const (
	defaultMiddlewarePriority = 1000

	logRequestMiddlewareId       = "sh_log_request"
	logRequestMiddlewarePriority = defaultMiddlewarePriority - 20

	panicRecoverMiddlewareId       = "sh_panic_recover"
	panicRecoverMiddlewarePriority = defaultMiddlewarePriority - 10
)

type responseWriter struct {
	status  int
	wrapped http.ResponseWriter
	written bool
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		status:  http.StatusOK,
		wrapped: w,
	}
}

func (rw *responseWriter) Header() http.Header {
	return rw.wrapped.Header()
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.wrapped.WriteHeader(status)

	if !rw.written {
		rw.status = status
		rw.written = true
	}
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.written = true
	return rw.wrapped.Write(b)
}

func (rw *responseWriter) Unwrap() http.ResponseWriter {
	return rw.wrapped
}

func logRequest() middleware {
	return middleware{
		id: logRequestMiddlewareId,
		fn: func(e *core.EventRequest, next http.Handler) {
			start := time.Now()

			rw := newResponseWriter(e.Response)
			next.ServeHTTP(rw, e.Request)

			duration := float64(time.Since(start)) / float64(time.Millisecond)

			e.App.Logger().Info("request details",
				slog.String("url", e.Request.RequestURI),
				slog.String("method", e.Request.Method),
				slog.Int("status_code", rw.status),
				slog.Float64("duration", duration),
			)
		},
		priority: logRequestMiddlewarePriority,
	}
}

func panicRecover() middleware {
	return middleware{
		id: panicRecoverMiddlewareId,
		fn: func(e *core.EventRequest, next http.Handler) {
			defer func() {
				if rec := recover(); rec != nil {
					stack := make([]byte, 2<<10) // 2KB
					length := runtime.Stack(stack, true)

					e.App.Logger().Error("panic recover",
						slog.String("error", fmt.Sprintf("%v - %s", rec, stack[:length])),
						slog.String("url", e.Request.RequestURI),
						slog.String("method", e.Request.Method),
					)

					err := e.InternalServerError("")
					_ = event.ErrorHandler(e.Response, e.Request, err)
				}
			}()

			next.ServeHTTP(e.Response, e.Request)
		},
		priority: panicRecoverMiddlewarePriority,
	}
}
