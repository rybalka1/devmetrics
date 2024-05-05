package http

import (
	"net"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/rybalka1/devmetrics/internal/config"
)

type Server interface {
	Start() error
	Stop() error
	AddMux(mux http.Handler)
}

type MetricServer struct {
	addr *net.TCPAddr
	http.Server
}

func NewServer(args config.Args) (Server, error) {
	return NewMetricServerWithParams(args.Addr)
}

func NewMetricServerWithParams(addr string) (*MetricServer, error) {
	netAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	srv := MetricServer{
		addr: netAddr,
		Server: http.Server{
			Addr: netAddr.String(),
		},
	}
	return &srv, nil
}

func (srv *MetricServer) AddMux(mux http.Handler) {
	srv.Server.Handler = mux
}

func (srv *MetricServer) Stop() error {
	return srv.Server.Close()
}

func (srv *MetricServer) Start() error {
	log.Info().Msgf("[+] Started on: %s", srv.Addr)
	return srv.ListenAndServe()
}
