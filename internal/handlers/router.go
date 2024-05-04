package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rybalka1/devmetrics/internal/memstorage"
	"github.com/rybalka1/devmetrics/internal/middlewares"
)

func NewRouter(store memstorage.Storage) http.Handler {
	r := chi.NewRouter()
	r.Use(middlewares.LoggingMiddleware)
	r.Route("/update/", func(r chi.Router) {
		r.Use(printMidl)
		r.Post("/counter/{name}/{value}", UpdateCounterHandle(store))
		r.Post("/gauge/{name}/{value}", UpdateGaugeHandle(store))
		r.Post("/", JSONUpdateOneMetricHandler(store))
		r.Post("/gauge/", NotFoundHandle)
		r.Post("/counter/", NotFoundHandle)
		r.Post("/gauge/{name}", NotFoundHandle)
		r.Post("/counter/{name}", NotFoundHandle)
		r.NotFound(BadRequest)
	})

	r.Route("/value/", func(r chi.Router) {
		r.Use(middlewares.LoggingMiddleware)
		r.Post("/", JSONGetMetricHandler(store))
		r.Get("/{mType}/{mName}", GetMetric(store))
	})

	r.Route("/ping", func(r chi.Router) {
		r.Get("/", pong)
	})

	return r
}

func printMidl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
