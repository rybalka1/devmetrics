package service

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rybalka1/devmetrics/internal/handlers"
	"github.com/rybalka1/devmetrics/internal/middlewares"
)

func NewRouter(service Service) http.Handler {
	r := chi.NewRouter()
	r.Use(middlewares.LoggingMiddleware)
	r.Use(middlewares.GzipMiddleware)
	r.Route("/", func(r chi.Router) {
		r.Get("/", handlers.GetAllMetrics(service))
	})
	r.Route("/update/", func(r chi.Router) {
		r.Post("/counter/{name}/{value}", handlers.UpdateCounterHandle(service))
		r.Post("/gauge/{name}/{value}", handlers.UpdateGaugeHandle(service))
		r.Post("/", handlers.JSONUpdateOneMetricHandler(service))
		r.Post("/gauge/", handlers.NotFoundHandle)
		r.Post("/counter/", handlers.NotFoundHandle)
		r.Post("/gauge/{name}", handlers.NotFoundHandle)
		r.Post("/counter/{name}", handlers.NotFoundHandle)
		r.NotFound(handlers.BadRequest)
	})
	r.Route("/value/", func(r chi.Router) {
		r.Post("/", handlers.JSONGetMetricHandler(service))
		r.Get("/{mType}/{mName}", handlers.GetMetric(service))
	})
	r.Route("/ping", func(r chi.Router) {
		r.Get("/", handlers.Pong)
	})
	return r
}

func printMidl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
