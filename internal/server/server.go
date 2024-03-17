package server

import (
	"fmt"
	"net"
	"net/http"

	"github.com/rybalka1/devmetrics/internal/handlers"
	"github.com/rybalka1/devmetrics/internal/memstorage"
)

type MetricServer struct {
	addr  *net.TCPAddr
	Store memstorage.Storage
	http.Server
}

func NewMetricServer(addr string) (*MetricServer, error) {
	store := memstorage.NewMemStorage()
	mux := handlers.NewRouter(store)
	return NewMetricServerWithParams(addr, store, mux)
}
func NewMetricServerWithParams(addr string, store memstorage.Storage, mux http.Handler) (*MetricServer, error) {
	netAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	if store == nil {
		store = memstorage.NewMemStorage()
	}
	if mux == nil {
		mux = handlers.NewRouter(store)

	}

	srv := MetricServer{
		addr:  netAddr,
		Store: store,
		Server: http.Server{
			Addr:    netAddr.String(),
			Handler: mux,
		},
	}
	return &srv, nil
}

func (srv *MetricServer) AddMux(mux *http.ServeMux) {
	srv.Server.Handler = mux
}

func (srv *MetricServer) Start() error {
	fmt.Println("[+] Started on:", srv.Addr)
	return srv.ListenAndServe()
}
