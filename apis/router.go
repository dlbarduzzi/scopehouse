package apis

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/dlbarduzzi/scopehouse/core"
	"github.com/dlbarduzzi/scopehouse/tools/event"
)

type route struct {
	pattern string
	handler func(*core.EventRequest)
}

type middleware struct {
	id       string
	fn       func(*core.EventRequest, http.Handler)
	priority int
}

type router struct {
	app         core.App
	routes      []route
	middlewares []middleware
}

func newRouter(app core.App) *router {
	r := &router{
		app:         app,
		routes:      []route{},
		middlewares: []middleware{},
	}

	r.bindMiddleware(logRequest())
	r.bindMiddleware(panicRecover())

	bindHealthApi(r)

	return r
}

func (r *router) add(pattern string, handler func(*core.EventRequest)) {
	r.routes = append(r.routes, route{
		pattern: pattern,
		handler: handler,
	})
}

func (r *router) get(pattern string, handler func(*core.EventRequest)) {
	r.add(fmt.Sprintf("GET %s", pattern), handler)
}

func (r *router) bindMiddleware(m middleware) {
	var exists bool

	if strings.TrimSpace(m.id) == "" {
		m.id = generateId(20)

	DUPLICATE_CHECK:
		for _, existing := range r.middlewares {
			if existing.id == m.id {
				m.id = generateId(20)
				goto DUPLICATE_CHECK
			}
		}
	} else {
		for i, existing := range r.middlewares {
			if existing.id == m.id {
				r.middlewares[i] = m
				exists = true
				break
			}
		}
	}

	if !exists {
		r.middlewares = append(r.middlewares, m)
	}

	sort.SliceStable(r.middlewares, func(i, j int) bool {
		return r.middlewares[i].priority < r.middlewares[j].priority
	})
}

func (r *router) buildMux() http.Handler {
	mux := http.NewServeMux()
	mws := make([]func(next http.Handler) http.Handler, 0)

	// Register middlewares.
	for _, middleware := range r.middlewares {
		mw := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				middleware.fn(&core.EventRequest{
					App: r.app,
					Event: event.Event{
						Request:  req,
						Response: res,
					},
				}, next)
			})
		}
		mws = append(mws, mw)
	}

	// Create API routes.
	for _, route := range r.routes {
		hf := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			route.handler(&core.EventRequest{
				App: r.app,
				Event: event.Event{
					Request:  req,
					Response: res,
				},
			})
		})

		pattern := route.pattern
		handler := http.Handler(hf)

		// Use middlewares.
		for i := len(mws) - 1; i >= 0; i-- {
			handler = mws[i](handler)
		}

		mux.Handle(pattern, handler)
	}

	return mux
}
