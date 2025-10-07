package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/rybalka1/devmetrics/internal/metrics"
	"github.com/rybalka1/devmetrics/internal/utils/compression/gzip"
)

func (agent Agent) SendMetrics() {
	if agent.metrics == nil {
		return
	}
	for mName, metric := range agent.metrics {
		url := fmt.Sprintf("http://%s/%s/%s/%s/%s", agent.addr.String(),
			agent.metricsPoint, metric.SendType, mName, metric.Value)
		fmt.Println(url)
		resp, err := http.Post(url, "text/plain", nil)
		if err != nil {
			continue
		}
		err = resp.Body.Close()
		if err != nil {
			continue
		}
	}
}

func (agent Agent) SendOneMetricJSON(name string, mymetric metrics.MyMetrics) error {
	URL := "update"
	metric, err := metrics.ConvertMymetric2Metric(name, mymetric)
	if err != nil {
		return fmt.Errorf("failed to convert metric %s: %w", name, err)
	}

	address := fmt.Sprintf("http://%s/%s/", agent.addr.String(), URL)
	body, err := json.Marshal(metric)
	if err != nil {
		return fmt.Errorf("failed to marshal metric %s: %w", name, err)
	}

	log.Debug().
		RawJSON("body", body).
		Str("metric_name", name).
		Msg("Sending metric")

	var buffer = bytes.NewBuffer(body)
	compressionStatus := false
	if agent.compression {
		compressed, err := gzip.Compress(body)
		if err == nil {
			buffer = bytes.NewBuffer(compressed)
			compressionStatus = true
		} else {
			log.Warn().Err(err).Msg("Failed to compress request body")
		}
	}

	request, err := http.NewRequest(http.MethodPost, address, buffer)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")
	if agent.compression && compressionStatus {
		request.Header.Set("Content-Encoding", "gzip")
	}

	// Use configurable timeout with reasonable defaults
	client := &http.Client{
		Timeout: time.Second * 30,
		Transport: &http.Transport{
			MaxIdleConns:        10,
			IdleConnTimeout:     30 * time.Second,
			DisableCompression:  false,
		},
	}

	resp, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("HTTP request failed for metric %s: %w", name, err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode >= 400 {
		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return fmt.Errorf("HTTP %d for metric %s (failed to read response body: %w)", resp.StatusCode, name, readErr)
		}
		return fmt.Errorf("HTTP %d for metric %s: %s", resp.StatusCode, name, string(body))
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body for metric %s: %w", name, err)
	}

	log.Debug().
		Str("url", URL).
		RawJSON("body", responseBody).
		Str("metric_name", name).
		Int("status_code", resp.StatusCode).
		Msg("Received response")

	return nil
}

func (agent Agent) SendMetricsJSON() error {
	for name, mymetric := range agent.metrics {
		err := agent.SendOneMetricJSON(name, mymetric)
		if err != nil {
			log.Error().
				Err(err).Send()
		}
	}
	return nil
}
