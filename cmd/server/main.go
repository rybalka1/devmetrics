package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/rybalka1/devmetrics/internal/server"
)

func selectArgs(addr *string) {
	*addr = os.Getenv("ADDRESS")
	var flagAddr string
	flag.StringVar(&flagAddr, "a", "localhost:8080", "address for server")
	flag.Parse()
	if *addr == "" {
		*addr = flagAddr
	}

}

func main() {
	var (
		addr string
	)
	selectArgs(&addr)
	fmt.Println(addr)
	srv, err := server.NewMetricServer(addr)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(srv.Start())
}
