package middleware

import "github.com/julienschmidt/httprouter"

// Chain wraps a given http.Handler with middlewares.
func Chain(handler httprouter.Handle, middlewares ...func(httprouter.Handle) httprouter.Handle) httprouter.Handle {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}

	return handler
}
