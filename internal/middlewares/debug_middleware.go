package middlewares

import "net/http"

func DebugStore(h http.Handler) http.Handler {
	debugFn := func(w http.ResponseWriter, r *http.Request) {
		// Do nothing because
		// You can add some logging or debugging statements here if needed
		h.ServeHTTP(w, r) // Call the original handler
	}
	return http.HandlerFunc(debugFn)
}
