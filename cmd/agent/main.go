package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/rybalka1/devmetrics/internal/agent"
)

const (
	defaultPollInterval   = 2 * time.Second
	defaultReportInterval = 10 * time.Second
)

func main() {
	var (
		addr                         string
		pollInterval, reportInterval int
	)
	flag.StringVar(&addr, "a", "localhost:8080",
		"address for metric server")

	flag.IntVar(&pollInterval, "r", int(defaultPollInterval.Seconds()),
		"frequency of gathering metrics")
	flag.IntVar(&reportInterval, "p", int(defaultReportInterval.Seconds()),
		"frequency of sending metrics")
	flag.Parse()

	mAgent, err := agent.NewAgent(addr, pollInterval, reportInterval)
	if err != nil {
		fmt.Println(err)
		return
	}
	mAgent.Start()
}
