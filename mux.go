package mux

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Router interface {
	GET(string, http.Handler) Router
	POST(string, http.Handler) Router
	PUT(string, http.Handler) Router
	PATCH(string, http.Handler) Router
	DELETE(string, http.Handler) Router
	GROUP(string, func(Router)) Router
	METHODS(...string) Router
	PATH(string) Router
	HOST(string) Router
	MIDDLEWARES(...MiddlewareConstructor) Router
	HANDLER(http.Handler) Router
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type gorwardRouter struct {
	R          *mux.Router
	Path       string
	Host       string
	Methods    []string
	Middleware []MiddlewareConstructor
}

type MiddlewareConstructor func(http.Handler) http.Handler

func Vars(r *http.Request) map[string]string {
	return mux.Vars(r)
}

func NewRouter(router *mux.Router) gorwardRouter {
	if router == nil {
		router = mux.NewRouter()
	}

	return gorwardRouter{R: router}
}

func (r gorwardRouter) HANDLER(handler http.Handler) Router {
	var route = r.R.NewRoute()
	if r.Host != "" {
		route = route.Host(r.Host)
	}
	if r.Methods != nil {
		route = route.Methods(r.Methods...)
	}
	if r.Path != "" {
		route = route.Path(r.Path)
	}
	if r.Middleware != nil {
		handler = bindMiddlewares(handler, r.Middleware...)
	}

	route.Handler(handler)

	return r
}

func bindMiddlewares(handler http.Handler, middlewares ...MiddlewareConstructor) http.Handler {
	var final http.Handler
	if handler != nil {
		final = handler
	} else {
		final = http.DefaultServeMux
	}

	for i := len(middlewares) - 1; i >= 0; i-- {
		final = middlewares[i](final)
	}

	return final
}

func (r gorwardRouter) MIDDLEWARES(middlewares ...MiddlewareConstructor) Router {
	n := &gorwardRouter{
		R:          r.R,
		Host:       r.Host,
		Path:       r.Path,
		Methods:    r.Methods,
		Middleware: middlewares,
	}

	return n
}

func (r gorwardRouter) METHODS(methods ...string) Router {
	n := &gorwardRouter{
		R:          r.R,
		Host:       r.Host,
		Path:       r.Path,
		Methods:    methods,
		Middleware: r.Middleware,
	}

	return n
}

func (r gorwardRouter) HOST(host string) Router {
	n := &gorwardRouter{
		R:          r.R,
		Host:       host,
		Path:       r.Path,
		Methods:    r.Methods,
		Middleware: r.Middleware,
	}

	return n
}

func (r gorwardRouter) PATH(path string) Router {
	n := &gorwardRouter{
		R:          r.R,
		Host:       r.Host,
		Path:       path,
		Methods:    r.Methods,
		Middleware: r.Middleware,
	}

	return n
}

func (r gorwardRouter) GROUP(path string, subRouter func(Router)) Router {
	n := &gorwardRouter{
		R:          r.R.PathPrefix(path).Subrouter(),
		Host:       r.Host,
		Methods:    r.Methods,
		Middleware: r.Middleware,
	}

	subRouter(n)

	return r
}

func (r gorwardRouter) GET(path string, handler http.Handler) Router {
	r.METHODS("GET").PATH(path).HANDLER(handler)

	return r
}

func (r gorwardRouter) POST(path string, handler http.Handler) Router {
	r.METHODS("POST").PATH(path).HANDLER(handler)

	return r
}

func (r gorwardRouter) PUT(path string, handler http.Handler) Router {
	r.METHODS("PUT").PATH(path).HANDLER(handler)

	return r
}

func (r gorwardRouter) PATCH(path string, handler http.Handler) Router {
	r.METHODS("PATCH").PATH(path).HANDLER(handler)

	return r
}

func (r gorwardRouter) DELETE(path string, handler http.Handler) Router {
	r.METHODS("DELETE").PATH(path).HANDLER(handler)

	return r
}

func (r gorwardRouter) OPTIONS(path string, handler http.Handler) Router {
	r.METHODS("OPTIONS").PATH(path).HANDLER(handler)

	return r
}

func (r gorwardRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.R.ServeHTTP(w, req)
}
