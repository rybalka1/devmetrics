package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/rybalka1/devmetrics/internal/storage/memstorage"
)

func UpdateMetricsHandle(store memstorage.Storage) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		metricsInfo, found := strings.CutPrefix(r.URL.Path, "/update/")
		if !found {
			LogAndWriteError(rw, http.StatusNotFound, fmt.Errorf("invalid path"), "Not found")
			return
		}

		pieces := strings.Split(metricsInfo, "/")
		if len(pieces) != 3 {
			LogAndWriteError(rw, http.StatusNotFound, fmt.Errorf("invalid path format"), "Invalid path format")
			return
		}

		mType := strings.TrimSpace(pieces[0])
		mName, err := SanitizeAndValidateMetricName(pieces[1])
		if err != nil {
			LogAndWriteError(rw, http.StatusBadRequest, err, "Invalid metric name")
			return
		}

		mValue := strings.TrimSpace(pieces[2])
		if err := ValidateMetricValue(mValue, mType); err != nil {
			LogAndWriteError(rw, http.StatusBadRequest, err, "Invalid metric value")
			return
		}

		switch mType {
		case "gauge":
			val, err := strconv.ParseFloat(mValue, 64)
			if err != nil {
				LogAndWriteError(rw, http.StatusBadRequest, err, "Invalid gauge value format")
				return
			}
			store.UpdateGauges(mName, val)
		case "counter":
			val, err := strconv.ParseInt(mValue, 10, 64)
			if err != nil {
				LogAndWriteError(rw, http.StatusBadRequest, err, "Invalid counter value format")
				return
			}
			store.UpdateCounters(mName, val)
		default:
			LogAndWriteError(rw, http.StatusBadRequest, fmt.Errorf("invalid metric type: %s", mType), "Invalid metric type")
			return
		}

		rw.Header().Set("Content-Type", "text/plain")
		rw.WriteHeader(http.StatusOK)
	}
}

func UpdateGaugeHandle(store memstorage.Storage) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		value := chi.URLParam(r, "value")

		// Validate and sanitize inputs
		mName, err := SanitizeAndValidateMetricName(name)
		if err != nil {
			LogAndWriteError(rw, http.StatusBadRequest, err, "Invalid metric name")
			return
		}

		if err := ValidateMetricValue(value, "gauge"); err != nil {
			LogAndWriteError(rw, http.StatusBadRequest, err, "Invalid gauge value")
			return
		}

		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			LogAndWriteError(rw, http.StatusBadRequest, err, "Invalid gauge value format")
			return
		}

		store.UpdateGauges(mName, val)
		rw.Header().Set("Content-Type", "text/plain")
		rw.WriteHeader(http.StatusOK)
	}
}

func UpdateCounterHandle(store memstorage.Storage) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		value := chi.URLParam(r, "value")

		// Validate and sanitize inputs
		mName, err := SanitizeAndValidateMetricName(name)
		if err != nil {
			LogAndWriteError(rw, http.StatusBadRequest, err, "Invalid metric name")
			return
		}

		if err := ValidateMetricValue(value, "counter"); err != nil {
			LogAndWriteError(rw, http.StatusBadRequest, err, "Invalid counter value")
			return
		}

		val, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			LogAndWriteError(rw, http.StatusBadRequest, err, "Invalid counter value format")
			return
		}

		store.UpdateCounters(mName, val)
		rw.Header().Set("Content-Type", "text/plain")
		rw.WriteHeader(http.StatusOK)
	}
}
