package handlers

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/rs/zerolog/log"
)

// Validation constants
const (
	MaxMetricNameLength = 100
	MaxMetricValueLength = 50
	MaxRequestBodySize = 4096
)

// Metric name validation regex - only allow alphanumeric, underscore, and hyphen
var validMetricNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// ValidateMetricName validates metric name for security and format
func ValidateMetricName(name string) error {
	if name == "" {
		return fmt.Errorf("metric name cannot be empty")
	}

	if utf8.RuneCountInString(name) > MaxMetricNameLength {
		return fmt.Errorf("metric name too long")
	}

	if !validMetricNameRegex.MatchString(name) {
		return fmt.Errorf("metric name contains invalid characters (allowed: letters, numbers, underscore, hyphen)")
	}

	return nil
}

// ValidateMetricValue validates metric value for format and length
func ValidateMetricValue(value string, mType string) error {
	if value == "" {
		return fmt.Errorf("metric value cannot be empty")
	}

	if utf8.RuneCountInString(value) > MaxMetricValueLength {
		return fmt.Errorf("metric value too long")
	}

	switch mType {
	case "gauge":
		// Basic float validation - allow scientific notation like "1.23e10"
		if strings.ContainsAny(value, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ") && !strings.Contains(value, ".") && !strings.Contains(strings.ToLower(value), "e") {
			return fmt.Errorf("invalid gauge value format")
		}
	case "counter":
		// Basic integer validation - no letters or decimal points allowed
		if strings.ContainsAny(value, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ.") {
			return fmt.Errorf("invalid counter value format")
		}
	}

	return nil
}

// ValidateContentType checks if content-type is application/json
func ValidateContentType(r *http.Request) error {
	valid := false
	vals := r.Header.Values("content-type")
	for _, val := range vals {
		if strings.Contains(val, "application/json") {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid content-type: %s (expected application/json)", strings.Join(vals, "; "))
	}
	return nil
}

// ValidateRequestBody reads and validates request body
func ValidateRequestBody(r *http.Request) ([]byte, error) {
	// Limit body size to prevent DoS
	r.Body = http.MaxBytesReader(nil, r.Body, MaxRequestBodySize)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read request body")
		return nil, fmt.Errorf("failed to read request body: %w", err)
	}

	if len(body) == 0 {
		return nil, fmt.Errorf("request body is empty")
	}

	return body, nil
}

// SanitizeAndValidateMetricName sanitizes and validates metric name
func SanitizeAndValidateMetricName(name string) (string, error) {
	// Trim whitespace
	name = strings.TrimSpace(name)

	// Validate the name
	if err := ValidateMetricName(name); err != nil {
		return "", err
	}

	return name, nil
}

// LogAndWriteError logs error and writes appropriate HTTP response
func LogAndWriteError(w http.ResponseWriter, statusCode int, err error, message string) {
	log.Error().Err(err).Int("status_code", statusCode).Msg(message)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := fmt.Sprintf(`{"error": "%s"}`, message)
	w.Write([]byte(response))
}
