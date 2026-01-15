package apis

import "net/http"

// routes wires all service API endpoints to their handlers.
func (s *service) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/health", s.healthCheck)
	return mux
}

// handler returns the HTTP handler for the service wrapped by middlewares.
func (s *service) handler(mux *http.ServeMux) http.Handler {
	var handler http.Handler = mux
	return handler
}
