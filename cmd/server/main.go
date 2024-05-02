package main

import (
	"fmt"
	"log"

	"github.com/rybalka1/devmetrics/internal/config"
	"github.com/rybalka1/devmetrics/internal/server"
)

func main() {
	var (
		addr     string
		logLevel string
	)
	config.ServerArgsParse(&addr, &logLevel)
	fmt.Println(addr, logLevel)
	srv, err := server.NewServer(addr, logLevel)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(srv.Start())
}
