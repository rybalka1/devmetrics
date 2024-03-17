package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/rybalka1/devmetrics/internal/server"
)

var (
	addr string
)

func initFlags() {
	flag.StringVar(&addr, "a", "localhost:8080", "address for server")
}

func main() {

	initFlags()
	flag.Parse()
	fmt.Println(addr)
	srv, err := server.NewMetricServer(addr)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(srv.Start())
}
