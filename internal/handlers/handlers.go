package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/rybalka1/devmetrics/internal/memstorage"
)

var (
	NotFound = "404 page not found\n"
)

func UpdateMetricsHandle(store memstorage.Storage) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		metricsInfo, found := strings.CutPrefix(r.URL.Path, "/update/")
		if !found {
			rw.Write([]byte(NotFound))
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		pieces := strings.Split(metricsInfo, "/")
		if len(pieces) != 3 {
			rw.Write([]byte(NotFound))
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		mType := pieces[0]
		mName := pieces[1]
		mValue := pieces[2]
		switch mType {
		case "gauge":
			val, err := strconv.ParseFloat(mValue, 64)
			if err != nil {
				rw.WriteHeader(http.StatusBadRequest)
				return
			}
			store.UpdateGauges(mName, val)
		case "counter":
			val, err := strconv.ParseInt(mValue, 10, 64)
			if err != nil {
				rw.WriteHeader(http.StatusBadRequest)
				return
			}
			store.UpdateCounters(mName, val)
		default:
			rw.WriteHeader(http.StatusBadRequest)
		}
		fmt.Println(store)
		rw.WriteHeader(http.StatusOK)
	}
}

func UpdateGaugeHandle(store memstorage.Storage) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		value := chi.URLParam(r, "value")
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		store.UpdateGauges(name, val)
		fmt.Println(store)
		rw.WriteHeader(http.StatusOK)
	}
}

func UpdateCounterHandle(store memstorage.Storage) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		value := chi.URLParam(r, "value")
		val, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		store.UpdateCounters(name, val)
		fmt.Println(store)
		rw.WriteHeader(http.StatusOK)
	}
}

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
		rw.WriteHeader(http.StatusOK)
		_, err := rw.Write([]byte(value))
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func NotFoundHandle(rw http.ResponseWriter, r *http.Request) {
	http.NotFound(rw, r)
}

func BadRequest(rw http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	rw.WriteHeader(http.StatusBadRequest)
}
