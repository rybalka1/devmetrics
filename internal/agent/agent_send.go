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
		return err
	}
	address := fmt.Sprintf("http://%s/%s/", agent.addr.String(), URL)
	body, err := json.Marshal(metric)
	if err != nil {
		return err
	}
	log.Info().
		RawJSON("body", body).
		Msg("Send: ")
	var buffer = bytes.NewBuffer(body)
	compressionStatus := false
	if agent.compression {
		compressed, err := gzip.Compress(body)
		if err == nil {
			buffer = bytes.NewBuffer(compressed)
			compressionStatus = true
		}
	}

	request, err := http.NewRequest(http.MethodPost, address, buffer)
	if err != nil {
		return err
	}
	request.Header.Set("content-type", "application/json")
	if agent.compression && compressionStatus {
		request.Header.Set("content-encoding", "gzip")
	}
	client := http.Client{Timeout: time.Second * 30}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Info().Str("url", URL).RawJSON("body", body).Msg("Receive: ")
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
