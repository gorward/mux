package mux

import (
	"github.com/gorilla/mux"
	"net/http"
)

// Router Interface
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

// GowardRouter holds the data of the route
type GowardRouter struct {
	R          *mux.Router
	Path       string
	Host       string
	Methods    []string
	Middleware []MiddlewareConstructor
}

// MiddlewareConstructor type Middleware
type MiddlewareConstructor func(http.Handler) http.Handler

// Vars gets path variables from request.
func Vars(r *http.Request) map[string]string {
	return mux.Vars(r)
}

// NewRouter creates an instance of GowardRouter
func NewRouter(router *mux.Router) GowardRouter {
	if router == nil {
		router = mux.NewRouter()
	}

	return GowardRouter{R: router}
}

// HANDLER create route using GorwardRouter values
func (r GowardRouter) HANDLER(handler http.Handler) Router {
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

// bindMiddlewares create handler by chaining middleware and handler
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

// MIDDLEWARES set GowardRoute's middlewares
func (r GowardRouter) MIDDLEWARES(middlewares ...MiddlewareConstructor) Router {
	n := &GowardRouter{
		R:          r.R,
		Host:       r.Host,
		Path:       r.Path,
		Methods:    r.Methods,
		Middleware: middlewares,
	}

	return n
}

// METHODS set GowardRoute's Methods
func (r GowardRouter) METHODS(methods ...string) Router {
	n := &GowardRouter{
		R:          r.R,
		Host:       r.Host,
		Path:       r.Path,
		Methods:    methods,
		Middleware: r.Middleware,
	}

	return n
}

// HOST set GowardRoute's Host
func (r GowardRouter) HOST(host string) Router {
	n := &GowardRouter{
		R:          r.R,
		Host:       host,
		Path:       r.Path,
		Methods:    r.Methods,
		Middleware: r.Middleware,
	}

	return n
}

// PATH set GowardRoute's path
func (r GowardRouter) PATH(path string) Router {
	n := &GowardRouter{
		R:          r.R,
		Host:       r.Host,
		Path:       path,
		Methods:    r.Methods,
		Middleware: r.Middleware,
	}

	return n
}

// GROUP creates a Route group
func (r GowardRouter) GROUP(path string, subRouter func(Router)) Router {
	n := &GowardRouter{
		R:          r.R.PathPrefix(path).Subrouter(),
		Host:       r.Host,
		Methods:    r.Methods,
		Middleware: r.Middleware,
	}

	subRouter(n)

	return r
}

// GET is a shorthand for creating Route with GET method
func (r GowardRouter) GET(path string, handler http.Handler) Router {
	r.METHODS("GET").PATH(path).HANDLER(handler)

	return r
}

// POST is a shorthand for creating Route with POST method
func (r GowardRouter) POST(path string, handler http.Handler) Router {
	r.METHODS("POST").PATH(path).HANDLER(handler)

	return r
}

// PUT is a shorthand for creating Route with PUT method
func (r GowardRouter) PUT(path string, handler http.Handler) Router {
	r.METHODS("PUT").PATH(path).HANDLER(handler)

	return r
}

// PATCH is a shorthand for creating Route with PATCH method
func (r GowardRouter) PATCH(path string, handler http.Handler) Router {
	r.METHODS("PATCH").PATH(path).HANDLER(handler)

	return r
}

// DELETE is a shorthand for creating Route with DELETE method
func (r GowardRouter) DELETE(path string, handler http.Handler) Router {
	r.METHODS("DELETE").PATH(path).HANDLER(handler)

	return r
}

// OPTIONS is a shorthand for creating Route with OPTIONS method
func (r GowardRouter) OPTIONS(path string, handler http.Handler) Router {
	r.METHODS("OPTIONS").PATH(path).HANDLER(handler)

	return r
}

// ServeHTTP to pass checks for type http.Handler
func (r GowardRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.R.ServeHTTP(w, req)
}
