package http

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

// StartWithGracefulShutdown starts the server with graceful shutdown support
func (srv *MetricServer) StartWithGracefulShutdown() error {
	// Create channel to listen for interrupt or terminate signals
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		log.Info().
			Str("address", srv.Addr).
			Msg("Starting HTTP server with graceful shutdown support")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server failed to start")
		}
	}()

	// Wait for shutdown signal
	<-done
	log.Info().Msg("Shutting down server...")

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Server forced to shutdown")
		return err
	}

	log.Info().Msg("Server exited gracefully")
	return nil
}
