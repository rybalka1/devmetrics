package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/rybalka1/devmetrics/internal/agent"
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
)

func main() {
	runtime.GOMAXPROCS(3)

	addr := "localhost:8080"

	mAgent, err := agent.NewAgent(addr, pollInterval, reportInterval)
	if err != nil {
		fmt.Println(err)
		return
	}
	mAgent.Start()
}
