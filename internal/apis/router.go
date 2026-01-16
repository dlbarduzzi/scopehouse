package apis

import (
	"fmt"
	"net/http"

	"github.com/dlbarduzzi/scopehouse/internal/core"
	"github.com/dlbarduzzi/scopehouse/internal/tools/event"
)

type apiRoute struct {
	pattern string
	handler func(*core.EventRequest)
}

type middleware struct {
	fn func(*core.EventRequest, http.Handler)
}

type router struct {
	app         core.App
	apiRoutes   []apiRoute
	middlewares []middleware
}

func newRouter(app core.App) *router {
	r := &router{
		app:         app,
		apiRoutes:   []apiRoute{},
		middlewares: []middleware{},
	}

	// routes wires all service API endpoints to their handlers.
	r.routes()

	return r
}

func (r *router) add(pattern string, handler func(*core.EventRequest)) {
	r.apiRoutes = append(r.apiRoutes, apiRoute{
		pattern: pattern,
		handler: handler,
	})
}

func (r *router) get(pattern string, handler func(*core.EventRequest)) {
	r.add(fmt.Sprintf("%s %s", http.MethodGet, pattern), handler)
}

func (r *router) handler() http.Handler {
	mux := http.NewServeMux()

	// register API routes
	for _, route := range r.apiRoutes {
		mux.HandleFunc(route.pattern, func(res http.ResponseWriter, req *http.Request) {
			e := &core.EventRequest{
				App: r.app,
				Event: event.Event{
					Request:  req,
					Response: res,
				},
			}
			route.handler(e)
		})
	}

	// Assign the mux to a http.Handler so it can be successively wrapped by
	// the middlewares, since each middleware returns a new http.Handler.
	var handler http.Handler = mux

	// register middlewares
	for _, m := range r.middlewares {
		mf := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				e := &core.EventRequest{
					App: r.app,
					Event: event.Event{
						Request:  req,
						Response: res,
					},
				}
				m.fn(e, next)
			})
		}
		// wrap handler with middleware
		handler = mf(handler)
	}

	return handler
}

func (r *router) routes() {
	r.get("/api/v1/health", healthCheck)
}
