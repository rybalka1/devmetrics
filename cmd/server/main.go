package main

import (
	"log"

	"github.com/rybalka1/devmetrics/internal/server"
)

func main() {

	addr := "localhost:8080"

	srv, err := server.NewMetricServer(addr)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(srv.Start())
}
