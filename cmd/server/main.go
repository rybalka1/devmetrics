package main

import (
	"log"
	"net/http"

	"github.com/rybalka1/devmetrics/internal/handlers"
	"github.com/rybalka1/devmetrics/internal/memstorage"
	"github.com/rybalka1/devmetrics/internal/server"
)

func main() {
	var store memstorage.Storage
	mux := http.NewServeMux()

	addr := "localhost:8080"

	store = memstorage.NewMemStorage()
	srv, err := server.NewMetricServer(addr, store, mux)
	if err != nil {
		log.Fatal(err)
	}

	mux.HandleFunc("/update/", handlers.UpdateMetricsHandle(srv.Store))
	//srv.AddMux(mux)

	log.Fatal(srv.Start())
}
