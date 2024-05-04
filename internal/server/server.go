package server

import (
	"net"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/rybalka1/devmetrics/internal/handlers"
	"github.com/rybalka1/devmetrics/internal/logger"
	"github.com/rybalka1/devmetrics/internal/memstorage"
)

type Server interface {
	InitLogger(level string) error
	Start() error
}

func NewServer(addr string, loggerLevel string) (Server, error) {
	store := memstorage.NewMemStorage()
	mux := handlers.NewRouter(store)
	return NewMetricServerWithParams(addr, store, mux, loggerLevel)
}

type MetricServer struct {
	addr  *net.TCPAddr
	Store memstorage.Storage
	http.Server
}

func (srv *MetricServer) InitLogger(level string) error {
	return logger.Initialize(level)
}

func NewMetricServer(addr string) (*MetricServer, error) {
	store := memstorage.NewMemStorage()
	mux := handlers.NewRouter(store)

	return NewMetricServerWithParams(addr, store, mux, "info")
}

func NewMetricServerWithParams(addr string,
	store memstorage.Storage,
	mux http.Handler, loggerLevel string) (*MetricServer, error) {

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

	err = srv.InitLogger(loggerLevel)
	if err != nil {
		return nil, err
	}
	return &srv, nil
}

func (srv *MetricServer) AddMux(mux *http.ServeMux) {
	srv.Server.Handler = mux
}

func (srv *MetricServer) Start() error {
	log.Info().Msgf("[+] Started on: %s", srv.Addr)
	return srv.ListenAndServe()
}
