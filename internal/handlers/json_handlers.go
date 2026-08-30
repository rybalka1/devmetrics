package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rybalka1/devmetrics/internal/metrics"
	"github.com/rybalka1/devmetrics/internal/storage/memstorage"
)

func JSONUpdateOneMetricHandler(store memstorage.Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Validate HTTP method
		if r.Method != http.MethodPost {
			LogAndWriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("wrong method: %s", r.Method), "Method not allowed")
			return
		}

		// Validate content type
		if err := ValidateContentType(r); err != nil {
			LogAndWriteError(w, http.StatusBadRequest, err, "Invalid content type")
			return
		}

		// Read and validate request body
		body, err := ValidateRequestBody(r)
		if err != nil {
			LogAndWriteError(w, http.StatusBadRequest, err, "Invalid request body")
			return
		}
		defer r.Body.Close()

		// Parse JSON
		var metric = new(metrics.Metrics)
		if err := json.Unmarshal(body, &metric); err != nil {
			LogAndWriteError(w, http.StatusBadRequest, err, "Invalid JSON format")
			return
		}

		// Validate metric structure
		if metric == nil {
			LogAndWriteError(w, http.StatusBadRequest, fmt.Errorf("metric is nil"), "Metric is nil")
			return
		}

		// Validate metric name
		if metric.ID == "" {
			LogAndWriteError(w, http.StatusBadRequest, fmt.Errorf("metric ID is empty"), "Metric ID is required")
			return
		}

		if _, err := SanitizeAndValidateMetricName(metric.ID); err != nil {
			LogAndWriteError(w, http.StatusBadRequest, err, "Invalid metric name")
			return
		}

		// Validate metric type and value
		if metric.MType != "gauge" && metric.MType != "counter" {
			LogAndWriteError(w, http.StatusBadRequest, fmt.Errorf("invalid metric type: %s", metric.MType), "Invalid metric type")
			return
		}

		if metric.Delta == nil && metric.Value == nil {
			LogAndWriteError(w, http.StatusBadRequest, fmt.Errorf("metric must have either delta or value"), "Metric must have a value")
			return
		}

		// Update metric in storage
		store.UpdateMetric(metric)
		retMetric := store.GetMetric(metric.ID, metric.MType)
		if retMetric == nil {
			LogAndWriteError(w, http.StatusNotFound, fmt.Errorf("metric %s not found after update", metric.ID), "Metric not found")
			return
		}

		// Return updated metric
		responseBody, err := json.Marshal(retMetric)
		if err != nil {
			LogAndWriteError(w, http.StatusInternalServerError, err, "Failed to marshal response")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(responseBody)
		if err != nil {
			LogAndWriteError(w, http.StatusInternalServerError, err, "Failed to write response")
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func JSONGetMetricHandler(store memstorage.Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Validate HTTP method
		if r.Method != http.MethodPost {
			LogAndWriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("wrong method: %s", r.Method), "Method not allowed")
			return
		}

		// Validate content type
		if err := ValidateContentType(r); err != nil {
			LogAndWriteError(w, http.StatusBadRequest, err, "Invalid content type")
			return
		}

		// Read and validate request body
		body, err := ValidateRequestBody(r)
		if err != nil {
			LogAndWriteError(w, http.StatusBadRequest, err, "Invalid request body")
			return
		}
		defer r.Body.Close()

		// Parse JSON
		var metric = new(metrics.Metrics)
		if err := json.Unmarshal(body, &metric); err != nil {
			LogAndWriteError(w, http.StatusBadRequest, err, "Invalid JSON format")
			return
		}

		// Validate metric structure
		if metric == nil {
			LogAndWriteError(w, http.StatusBadRequest, fmt.Errorf("metric is nil"), "Metric is nil")
			return
		}

		// Validate metric name
		if metric.ID == "" {
			LogAndWriteError(w, http.StatusBadRequest, fmt.Errorf("metric ID is empty"), "Metric ID is required")
			return
		}

		if _, err := SanitizeAndValidateMetricName(metric.ID); err != nil {
			LogAndWriteError(w, http.StatusBadRequest, err, "Invalid metric name")
			return
		}

		// Validate metric type if provided
		if metric.MType != "" && metric.MType != "gauge" && metric.MType != "counter" {
			LogAndWriteError(w, http.StatusBadRequest, fmt.Errorf("invalid metric type: %s", metric.MType), "Invalid metric type")
			return
		}

		// Get metric from storage
		retMetric := store.GetMetric(metric.ID, metric.MType)
		if retMetric == nil {
			LogAndWriteError(w, http.StatusNotFound, fmt.Errorf("metric %s not found", metric.ID), "Metric not found")
			return
		}

		// Return metric
		responseBody, err := json.Marshal(retMetric)
		if err != nil {
			LogAndWriteError(w, http.StatusInternalServerError, err, "Failed to marshal response")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(responseBody)
		if err != nil {
			LogAndWriteError(w, http.StatusInternalServerError, err, "Failed to write response")
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
