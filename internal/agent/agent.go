package agent

import (
	"fmt"
	"net"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/rybalka1/devmetrics/internal/config"
	"github.com/rybalka1/devmetrics/internal/logger"
	"github.com/rybalka1/devmetrics/internal/metrics"
	"github.com/rybalka1/devmetrics/internal/storage/memstorage"
)

var usedMemStats = []string{
	"Alloc",
	"BuckHashSys",
	"Frees",
	"GCCPUFraction",
	"GCSys",
	"HeapAlloc",
	"HeapIdle",
	"HeapInuse",
	"HeapObjects",
	"HeapReleased",
	"HeapSys",
	"LastGC",
	"Lookups",
	"MCacheInuse",
	"MCacheSys",
	"MSpanInuse",
	"MSpanSys",
	"Mallocs",
	"NextGC",
	"NumForcedGC",
	"NumGC",
	"OtherSys",
	"PauseTotalNs",
	"StackInuse",
	"StackSys",
	"Sys",
	"TotalAlloc",
}

type Agent struct {
	addr           net.Addr
	metricsPoint   string
	metrics        map[string]metrics.MyMetrics
	store          memstorage.Storage
	pollInterval   time.Duration
	reportInterval time.Duration
	compression    bool
	logLevel       string
}

func (agent Agent) InitLogger(loggerLevel string) error {
	return logger.Initialize(loggerLevel)
}

func NewAgent(cfg config.AgentConfig) (*Agent, error) {
	netAddr, err := net.ResolveTCPAddr("tcp", cfg.Address)

	if err != nil {
		return nil, err
	}

	agent := &Agent{
		addr:           netAddr,
		metricsPoint:   "update",
		metrics:        make(map[string]metrics.MyMetrics),
		pollInterval:   time.Duration(cfg.PollInterval) * time.Second,
		reportInterval: time.Duration(cfg.ReportInterval) * time.Second,
		logLevel:       "debug",
	}
	err = agent.InitLogger(agent.logLevel)

	if err != nil {
		return nil, err
	}
	return agent, nil
}

func (agent Agent) Start() error {
	const (
		maxErrCount     = 3
		maxRetryBackoff = 30 * time.Second
	)

	var curErrCount int
	pollTicker := time.NewTicker(agent.pollInterval)
	reportTicker := time.NewTicker(agent.reportInterval)

	defer func() {
		pollTicker.Stop()
		reportTicker.Stop()
	}()

	log.Info().
		Dur("poll_interval", agent.pollInterval).
		Dur("report_interval", agent.reportInterval).
		Msg("Starting agent")

	for {
		select {
		case t1 := <-pollTicker.C:
			log.Debug().Time("poll_time", t1).Msg("Collecting metrics")
			agent.GetMetrics()
		case t2 := <-reportTicker.C:
			log.Debug().Time("report_time", t2).Msg("Sending metrics batch")

			err := agent.SendMetricsJSON()
			if err != nil {
				curErrCount++
				log.Error().
					Err(err).
					Int("error_count", curErrCount).
					Int("max_errors", maxErrCount).
					Msg("Failed to send metrics")

				// Implement exponential backoff for retry interval
				if curErrCount >= maxErrCount {
					log.Error().
						Int("error_count", curErrCount).
						Msg("Maximum error count reached, stopping agent")
					return fmt.Errorf("agent stopped due to repeated errors: %w", err)
				}

				// Reset error count on successful send
				log.Info().Msg("Will retry sending metrics on next interval")
			} else {
				// Reset error count on success
				if curErrCount > 0 {
					log.Info().Int("previous_errors", curErrCount).Msg("Metrics sent successfully, error count reset")
					curErrCount = 0
				}
			}
		}
	}
}
