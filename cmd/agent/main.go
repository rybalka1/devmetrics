package main

import (
	"fmt"
	"log"

	"github.com/rybalka1/devmetrics/internal/agent"
	"github.com/rybalka1/devmetrics/internal/config"
)

func main() {
	cfg, err := config.LoadUnifiedConfig()
	if err != nil {
		log.Fatal(err)
	}
	agentConfig := cfg.GetAgentConfig()
	mAgent, err := agent.NewAgent(agentConfig)

	if err != nil {
		fmt.Println(err)
		return
	}
	log.Fatal(mAgent.Start())
}
