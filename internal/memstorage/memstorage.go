package memstorage

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/rybalka1/devmetrics/internal/metrics"
)

type Storage interface {
	UpdateCounters(name string, value int64)
	UpdateGauges(name string, value float64)
	String() string
	GetMetricString(mType, mName string) string
	GetOneMetric(mName string) *metrics.Metrics
	GetMetric(mName string, mType string) *metrics.Metrics
	UpdateMetrics(newMetrics []*metrics.Metrics) []error
	UpdateMetric(newMetric *metrics.Metrics)
	AddMetric(newMetric *metrics.Metrics)
}

type MemStorage struct {
	dataCounters map[string]int64
	dataGauges   map[string]float64
	metrics      []*metrics.Metrics
	mu           sync.RWMutex
}

type MetricStorage struct {
	metrics []*metrics.Metrics
}

func (ms *MemStorage) String() string {
	var s1 = "[counters]\n"
	var s2 = "[gauges]\n"
	count := 0
	numInLine := 5
	for key, val := range ms.dataCounters {
		s1 += fmt.Sprintf("%s:%d ", key, val)
		count++
		if count == numInLine {
			count = 0
			s1 += "\n"
		}
	}
	s1 += "\n"
	count = 0
	for key, val := range ms.dataGauges {
		s2 += fmt.Sprintf("%s:%f ", key, val)
		count++
		if count == numInLine {
			count = 0
			s2 += "\n"
		}
	}
	s2 += "\n\n"

	return s1 + s2

}

func (ms *MemStorage) UpdateMetric(newMetric *metrics.Metrics) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	switch newMetric.MType {
	case metrics.Gauge:
		ms.dataGauges[newMetric.ID] = *newMetric.Value
		return
	case metrics.Counter:
		ms.dataCounters[newMetric.ID] += *newMetric.Delta
		return
	}
}

func (ms *MemStorage) AddMetric(newMetric *metrics.Metrics) {
	ms.UpdateMetric(newMetric)
}

func (ms *MemStorage) GetMetricString(mType, mName string) string {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	switch mType {
	case metrics.Gauge:
		val, ok := ms.dataGauges[mName]
		if !ok {
			return ""
		}
		return strconv.FormatFloat(val, 'f', -1, 64)
	case metrics.Counter:
		val, ok := ms.dataCounters[mName]
		if !ok {
			return ""
		}
		return strconv.FormatInt(val, 10)
	}
	return ""
}

func (ms *MemStorage) GetMetric(mName string, mType string) *metrics.Metrics {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	switch mType {
	case metrics.Gauge:
		for key, value := range ms.dataGauges {
			if key == mName {
				return &metrics.Metrics{
					ID:    key,
					MType: metrics.Gauge,
					Delta: nil,
					Value: &value,
				}
			}
		}
	case metrics.Counter:
		for key, value := range ms.dataCounters {
			if key == mName {
				return &metrics.Metrics{
					ID:    key,
					MType: metrics.Counter,
					Delta: &value,
					Value: nil,
				}
			}
		}
	}

	return nil
}

func (ms *MemStorage) GetOneMetric(mName string) *metrics.Metrics {
	for key, value := range ms.dataGauges {
		if key == mName {
			return &metrics.Metrics{
				ID:    key,
				MType: metrics.Gauge,
				Delta: nil,
				Value: &value,
			}
		}
	}

	for key, value := range ms.dataCounters {
		if key == mName {
			return &metrics.Metrics{
				ID:    key,
				MType: metrics.Counter,
				Delta: &value,
				Value: nil,
			}
		}
	}
	return nil
}

func NewMemStorage() *MemStorage {
	ms := new(MemStorage)
	ms.dataCounters = make(map[string]int64)
	ms.dataGauges = make(map[string]float64)
	return ms
}
