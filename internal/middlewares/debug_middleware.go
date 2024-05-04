package middlewares

import "net/http"

func DebugStore(h http.Handler) http.Handler {
	debugFn := func(w http.ResponseWriter, r *http.Request) {
		// Do nothing because
	}
	return http.HandlerFunc(debugFn)
}
