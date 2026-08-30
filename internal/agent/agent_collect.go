package agent

import (
	"math/rand"
	"reflect"
	"runtime"
	"strconv"
	"time"

	"github.com/rybalka1/devmetrics/internal/metrics"
	"github.com/rybalka1/devmetrics/internal/storage/memstorage"
)

func (agent Agent) GetMetrics() {
	// Pre-allocate map with estimated size to reduce allocations
	if agent.metrics == nil {
		agent.metrics = make(map[string]metrics.MyMetrics, len(usedMemStats)+2)
	}

	// Read memory stats once
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	values := reflect.ValueOf(m)

	for _, name := range usedMemStats {
		field := values.FieldByName(name)
		if !field.IsValid() {
			continue
		}

		var valueStr string

		switch {
		case field.CanInt():
			intVal := field.Int()
			valueStr = strconv.FormatInt(intVal, 10)
		case field.CanUint():
			uintVal := field.Uint()
			valueStr = strconv.FormatUint(uintVal, 10)
		case field.CanFloat():
			floatVal := field.Float()
			valueStr = strconv.FormatFloat(floatVal, 'f', -1, 64)
		default:
			continue
		}

		agent.metrics[name] = metrics.MyMetrics{
			Value:    string(valueStr),
			SendType: metrics.Gauge,
		}
	}

	// Update PollCount counter efficiently
	if metric, exists := agent.metrics["PollCount"]; exists {
		metric.AddVal(1)
		agent.metrics["PollCount"] = metric
	} else {
		agent.metrics["PollCount"] = metrics.MyMetrics{
			Value:    "1",
			SendType: metrics.Counter,
		}
	}

	// Generate RandomValue with better seed management
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randVal := float64(r.Intn(1000)) + r.Float64()

	agent.metrics["RandomValue"] = metrics.MyMetrics{
		Value:    strconv.FormatFloat(randVal, 'f', -1, 64),
		SendType: metrics.Gauge,
	}
}

func (agent Agent) CollectMetrics() {
	agent.store = memstorage.NewMemStorage()
	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)
	values := reflect.ValueOf(m)
	for _, name := range usedMemStats {
		if values.FieldByName(name).IsValid() {
			if values.FieldByName(name).CanInt() {
				val := values.FieldByName(name).Int()
				agent.store.UpdateMetric(&metrics.Metrics{
					ID:    name,
					MType: metrics.Counter,
					Delta: &val,
					Value: nil,
				})
			}
			if values.FieldByName(name).CanFloat() {
				val := values.FieldByName(name).Float()
				agent.store.UpdateMetric(&metrics.Metrics{
					ID:    name,
					MType: metrics.Gauge,
					Delta: nil,
					Value: &val,
				})
			}
		}
	}

	metric, ok := agent.metrics["PollCount"]
	if !ok {
		agent.metrics["PollCount"] = metrics.MyMetrics{
			Value:    "1",
			SendType: metrics.Counter,
		}
	} else {
		metric.AddVal(1)
		agent.metrics["PollCount"] = metric
	}

	r := rand.New(rand.NewSource(time.Now().Unix()))
	randVal := float64(r.Intn(1000)) + r.Float64()
	agent.metrics["RandomValue"] = metrics.MyMetrics{
		Value:    strconv.FormatFloat(randVal, 'f', -1, 64),
		SendType: metrics.Gauge,
	}
}
