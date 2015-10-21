package main

import (
	"fmt"
	"github.com/gorward/mux"
	"net/http"
	"time"
)

func main() {
	router := mux.NewRouter(nil)
	router.Get("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "adljasldk")
	}))

	router.Get("/x", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "yyyy")
	}))

	router.Middlewares(HowLong, MiddlewareGenerator("one"), MiddlewareGenerator("two")).Group("/api", func(r mux.Router) {
		r.Get("/", HandlerGenerator("api:get"))
		r.Get("/user", HandlerGenerator("user:get"))
		r.Post("/user", HandlerGenerator("user:post"))
		r.Get("/user/{id}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := mux.Vars(r)["id"]

			fmt.Fprintf(w, "Hello, %s!", id)
		}))
	})

	http.ListenAndServe(":8080", router)
}

func HowLong(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		fmt.Printf("[%s]%s STARTED \n", r.Method, r.URL.Path)
		h.ServeHTTP(w, r)
		fmt.Printf("[%s]%s COMPLETED [%s] \n", r.Method, r.URL.Path, time.Since(start))

	})
}

func MiddlewareGenerator(s string) mux.MiddlewareConstructor {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("[%s]%s : %s \n", r.Method, r.URL.Path, s)
			handler.ServeHTTP(w, r)
		})
	}
}

func HandlerGenerator(s string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "path: %s\n output: %s", r.URL.Path, s)
	})
}
