package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/rybalka1/devmetrics/internal/storage/memstorage"
)

func GetMetric(store memstorage.Storage) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		mType := strings.TrimSpace(chi.URLParam(r, "mType"))
		mName := strings.TrimSpace(chi.URLParam(r, "mName"))

		// Validate metric type
		if mType != "gauge" && mType != "counter" {
			LogAndWriteError(rw, http.StatusBadRequest, fmt.Errorf("invalid metric type: %s", mType), "Invalid metric type")
			return
		}

		// Validate and sanitize metric name
		_, err := SanitizeAndValidateMetricName(mName)
		if err != nil {
			LogAndWriteError(rw, http.StatusBadRequest, err, "Invalid metric name")
			return
		}

		value := store.GetMetricString(mType, mName)
		if value == "" {
			LogAndWriteError(rw, http.StatusNotFound, fmt.Errorf("metric %s/%s not found", mType, mName), "Metric not found")
			return
		}

		rw.Header().Set("Content-Type", "text/plain")
		_, err = rw.Write([]byte(value))
		if err != nil {
			LogAndWriteError(rw, http.StatusInternalServerError, err, "Failed to write response")
			return
		}
		rw.WriteHeader(http.StatusOK)
	}
}

func GetAllMetrics(store memstorage.Storage) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		metricsStr := store.String()
		if metricsStr == "" {
			LogAndWriteError(rw, http.StatusInternalServerError, fmt.Errorf("no metrics available"), "No metrics found")
			return
		}

		page := fmt.Sprintf("<html><body>%s</body></html>", metricsStr)
		rw.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, err := fmt.Fprint(rw, page)
		if err != nil {
			LogAndWriteError(rw, http.StatusInternalServerError, err, "Failed to write response")
			return
		}
		rw.WriteHeader(http.StatusOK)
	}
}
