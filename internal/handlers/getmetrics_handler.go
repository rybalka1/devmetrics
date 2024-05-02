package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/rybalka1/devmetrics/internal/memstorage"
	"net/http"
)

func GetMetric(store memstorage.Storage) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		mType := chi.URLParam(r, "mType")
		mName := chi.URLParam(r, "mName")
		value := store.GetMetric(mType, mName)
		if value == "" {
			rw.WriteHeader(http.StatusNotFound)
			rw.Write([]byte(NotFound))
			return
		}
		rw.Header().Add("Content-type", "text/plain")
		_, err := rw.Write([]byte(value))
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		rw.WriteHeader(http.StatusOK)
	}
}
