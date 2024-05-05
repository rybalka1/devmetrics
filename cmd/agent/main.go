package main

import (
	"fmt"
	"log"

	"github.com/rybalka1/devmetrics/internal/agent"
	"github.com/rybalka1/devmetrics/internal/config"
)

func main() {
	var (
		addr                         string
		pollInterval, reportInterval int
	)

	config.AgentSelectArgs(&addr, &pollInterval, &reportInterval)
	mAgent, err := agent.NewAgent(addr, pollInterval, reportInterval)

	if err != nil {
		fmt.Println(err)
		return
	}
	log.Fatal(mAgent.Start())
}
