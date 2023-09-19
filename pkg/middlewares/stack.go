package middlewares

import "net/http"

type Middleware func(next http.Handler) http.Handler

func WrapHandler(h http.Handler, chain ...Middleware) http.Handler {
	for _, f := range chain {
		h = f(h)
	}

	return h
}
