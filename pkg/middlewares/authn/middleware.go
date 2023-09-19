package authn

import "net/http"

// a dummy authentication middleware
func WithAPIKey(apiKey string) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if apiKey != "" {
				val := r.Header.Get("Authorization")
				if val != apiKey {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Not Authenticated"))
					return
				}
			}

			h.ServeHTTP(w, r)
		})
	}
}
