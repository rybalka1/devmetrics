package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateMetricName(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
		errorMsg  string
	}{
		{
			name:      "valid metric name",
			input:     "test_metric_123",
			wantError: false,
		},
		{
			name:      "empty name",
			input:     "",
			wantError: true,
			errorMsg:  "metric name cannot be empty",
		},
		{
			name:      "name too long",
			input:     strings.Repeat("a", MaxMetricNameLength+1),
			wantError: true,
			errorMsg:  "metric name too long",
		},
		{
			name:      "invalid characters",
			input:     "test@metric#123",
			wantError: true,
			errorMsg:  "metric name contains invalid characters",
		},
		{
			name:      "valid name with hyphen",
			input:     "test-metric-123",
			wantError: false,
		},
		{
			name:      "valid name with underscore",
			input:     "test_metric_123",
			wantError: false,
		},
		{
			name:      "path traversal attempt",
			input:     "../../../etc/passwd",
			wantError: true,
			errorMsg:  "metric name contains invalid characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMetricName(tt.input)
			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateMetricValue(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		mType     string
		wantError bool
		errorMsg  string
	}{
		{
			name:      "valid gauge value",
			value:     "123.456",
			mType:     "gauge",
			wantError: false,
		},
		{
			name:      "valid counter value",
			value:     "123",
			mType:     "counter",
			wantError: false,
		},
		{
			name:      "empty value",
			value:     "",
			mType:     "gauge",
			wantError: true,
			errorMsg:  "metric value cannot be empty",
		},
		{
			name:      "value too long",
			value:     strings.Repeat("1", MaxMetricValueLength+1),
			mType:     "gauge",
			wantError: true,
			errorMsg:  "metric value too long",
		},
		{
			name:      "invalid gauge with letters",
			value:     "abc123",
			mType:     "gauge",
			wantError: true,
			errorMsg:  "invalid gauge value format",
		},
		{
			name:      "invalid counter with letters",
			value:     "123abc",
			mType:     "counter",
			wantError: true,
			errorMsg:  "invalid counter value format",
		},
		{
			name:      "valid gauge with scientific notation",
			value:     "1.23e10",
			mType:     "gauge",
			wantError: false,
		},
		{
			name:      "valid gauge with capital E",
			value:     "2.5E-5",
			mType:     "gauge",
			wantError: false,
		},
		{
			name:      "invalid gauge with letters but no decimal",
			value:     "abc123",
			mType:     "gauge",
			wantError: true,
			errorMsg:  "invalid gauge value format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMetricValue(tt.value, tt.mType)
			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateContentType(t *testing.T) {
	tests := []struct {
		name      string
		contentType string
		wantError   bool
	}{
		{
			name:        "valid json content type",
			contentType: "application/json",
			wantError:   false,
		},
		{
			name:        "valid json with charset",
			contentType: "application/json; charset=utf-8",
			wantError:   false,
		},
		{
			name:        "invalid content type",
			contentType: "text/plain",
			wantError:   true,
		},
		{
			name:        "empty content type",
			contentType: "",
			wantError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/", nil)
			if tt.contentType != "" {
				req.Header.Set("Content-Type", tt.contentType)
			}

			err := ValidateContentType(req)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSanitizeAndValidateMetricName(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      string
		wantError bool
	}{
		{
			name:      "trim whitespace",
			input:     "  test_metric  ",
			want:      "test_metric",
			wantError: false,
		},
		{
			name:      "valid name",
			input:     "valid_name_123",
			want:      "valid_name_123",
			wantError: false,
		},
		{
			name:      "invalid name",
			input:     "invalid@name",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SanitizeAndValidateMetricName(tt.input)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, result)
			}
		})
	}
}

func TestLogAndWriteError(t *testing.T) {
	// Test that LogAndWriteError writes proper JSON error response
	w := httptest.NewRecorder()

	testError := assert.AnError
	testMessage := "Test error message"

	LogAndWriteError(w, http.StatusBadRequest, testError, testMessage)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	body := w.Body.String()
	assert.Contains(t, body, `"error": "Test error message"`)
}
