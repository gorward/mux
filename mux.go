package mux

import (
	"github.com/gorilla/mux"
	"net/http"
)

// Router Interface
type Router interface {
	Get(string, http.Handler) Router
	Post(string, http.Handler) Router
	Put(string, http.Handler) Router
	Patch(string, http.Handler) Router
	Delete(string, http.Handler) Router
	Group(string, func(Router)) Router
	Methods(...string) Router
	Path(string) Router
	Host(string) Router
	Middlewares(...MiddlewareConstructor) Router
	Handler(http.Handler) Router
	ServeHTTP(http.ResponseWriter, *http.Request)
}

// GowardRouter holds the data of the route
type GowardRouter struct {
	FParentRouter *mux.Router
	FRouter       *mux.Router
	FPath         string
	FPrefix       string
	FHost         string
	FMethods      []string
	FMiddleware   []MiddlewareConstructor
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

	return GowardRouter{FRouter: router}
}

// HANDLER create route using GorwardRouter values
func (r GowardRouter) Handler(handler http.Handler) Router {
	var route *mux.Route

	if r.FMiddleware != nil {
		handler = bindMiddlewares(handler, r.FMiddleware...)
	}

	if r.FPath == "/" && r.FParentRouter != nil {
		route = r.FParentRouter.NewRoute()
		route = route.Path(r.FPrefix)
	} else {
		route = r.FRouter.NewRoute()
		route.Path(r.FPath)
	}
	if r.FHost != "" {
		route = route.Host(r.FHost)
	}
	if r.Methods != nil {
		route = route.Methods(r.FMethods...)
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
func (r GowardRouter) Middlewares(middlewares ...MiddlewareConstructor) Router {
	n := &GowardRouter{
		FParentRouter: r.FParentRouter,
		FRouter:       r.FRouter,
		FHost:         r.FHost,
		FPath:         r.FPath,
		FPrefix:       r.FPrefix,
		FMethods:      r.FMethods,
		FMiddleware:   middlewares,
	}

	return n
}

// METHODS set GowardRoute's Methods
func (r GowardRouter) Methods(methods ...string) Router {
	n := &GowardRouter{
		FParentRouter: r.FParentRouter,
		FRouter:       r.FRouter,
		FHost:         r.FHost,
		FPath:         r.FPath,
		FPrefix:       r.FPrefix,
		FMethods:      methods,
		FMiddleware:   r.FMiddleware,
	}

	return n
}

// HOST set GowardRoute's Host
func (r GowardRouter) Host(host string) Router {
	n := &GowardRouter{
		FParentRouter: r.FParentRouter,
		FRouter:       r.FRouter,
		FHost:         host,
		FPath:         r.FPath,
		FPrefix:       r.FPrefix,
		FMethods:      r.FMethods,
		FMiddleware:   r.FMiddleware,
	}

	return n
}

// PATH set GowardRoute's path
func (r GowardRouter) Path(path string) Router {
	n := &GowardRouter{
		FParentRouter: r.FParentRouter,
		FRouter:       r.FRouter,
		FHost:         r.FHost,
		FPath:         path,
		FPrefix:       r.FPrefix,
		FMethods:      r.FMethods,
		FMiddleware:   r.FMiddleware,
	}

	return n
}

// GROUP creates a Route group
func (r GowardRouter) Group(path string, subRouter func(Router)) Router {
	n := &GowardRouter{
		FParentRouter: r.FRouter,
		FRouter:       r.FRouter.PathPrefix(path).Subrouter(),
		FHost:         r.FHost,
		FPrefix:       path,
		FMethods:      r.FMethods,
		FMiddleware:   r.FMiddleware,
	}

	subRouter(n)

	return r
}

// GET is a shorthand for creating Route with GET method
func (r GowardRouter) Get(path string, handler http.Handler) Router {
	r.Methods("GET").Path(path).Handler(handler)

	return r
}

// Post is a shorthand for creating Route with Post method
func (r GowardRouter) Post(path string, handler http.Handler) Router {
	r.Methods("POST").Path(path).Handler(handler)

	return r
}

// Put is a shorthand for creating Route with Put method
func (r GowardRouter) Put(path string, handler http.Handler) Router {
	r.Methods("PUT").Path(path).Handler(handler)

	return r
}

// Patch is a shorthand for creating Route with Patch method
func (r GowardRouter) Patch(path string, handler http.Handler) Router {
	r.Methods("PATCH").Path(path).Handler(handler)

	return r
}

// Delete is a shorthand for creating Route with Delete method
func (r GowardRouter) Delete(path string, handler http.Handler) Router {
	r.Methods("DELETE").Path(path).Handler(handler)

	return r
}

// OPTIONS is a shorthand for creating Route with OPTIONS method
func (r GowardRouter) OPTIONS(path string, handler http.Handler) Router {
	r.Methods("OPTIONS").Path(path).Handler(handler)

	return r
}

// ServeHTTP to pass checks for type http.Handler
func (r GowardRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.FRouter.ServeHTTP(w, req)
}
