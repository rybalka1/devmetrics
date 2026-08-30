package middlewares

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/rs/zerolog/log"
)

type responseData struct {
	respSize   int
	statusCode int
}
type loggerRW struct {
	http.ResponseWriter
	*responseData
}

// Write
func (lrw *loggerRW) Write(buf []byte) (int, error) {
	size, err := lrw.ResponseWriter.Write(buf)
	lrw.responseData.respSize = size
	return size, err
}

func (lrw *loggerRW) WriteHeader(statusCode int) {
	lrw.ResponseWriter.WriteHeader(statusCode)
	lrw.responseData.statusCode = statusCode
}

func LoggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Enhanced logging with more details
		log.Info().
			Str("method", r.Method).
			Str("uri", r.RequestURI).
			Str("remote_addr", r.RemoteAddr).
			Str("user_agent", r.UserAgent()).
			Str("host", r.Host).
			Msg("HTTP request started")

		// Panic recovery
		defer func() {
			if err := recover(); err != nil {
				log.Error().
					Interface("panic", err).
					Str("stack", string(debug.Stack())).
					Str("method", r.Method).
					Str("uri", r.RequestURI).
					Msg("Panic recovered in HTTP handler")

				// Return 500 Internal Server Error
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()

		// Wrap response writer for logging
		var (
			data = &responseData{}
			lrw  = &loggerRW{ResponseWriter: w, responseData: data}
		)

		// Serve the request
		h.ServeHTTP(lrw, r)

		// Calculate duration
		duration := time.Since(start)

		// Enhanced response logging
		logger := log.Info().
			Str("method", r.Method).
			Str("uri", r.RequestURI).
			Int("status", data.statusCode).
			Int("response_size", data.respSize).
			Dur("duration", duration)

		// Log different levels based on status code
		if data.statusCode >= 500 {
			logger.Msg("HTTP request completed with server error")
		} else if data.statusCode >= 400 {
			logger.Msg("HTTP request completed with client error")
		} else if data.statusCode >= 300 {
			logger.Msg("HTTP request completed with redirection")
		} else {
			logger.Msg("HTTP request completed successfully")
		}
	})
}
