gorward/mux
=======

don't use yet

This is just a wrapper for gorilla/mux

This is an attempt to make grouping routes easier in go.

Look at example/example.go for some examples. :)

Example
-----
```go
	router := mux.NewRouter(nil)
	router.GET("/", IndexGetHandler)
	router.POST("/", IndexPostHandler)
	router.MIDDLEWARES(Logger, AuthBasic).GROUP("/api", func(subRouter mux.Router) {
    	subRouter.GET("/user", UserGetHandler)
   	 subRouter.POST("/user", UserPostHandler)
    })
	http.ListenAndServe(":8080", router)
```

todo
-----

- [X] grouping of routes
- [X] middleware support 
- [X] accessing variables
- [ ] Tests
