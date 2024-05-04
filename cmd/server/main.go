package main

import (
	"github.com/rs/zerolog/log"
	"github.com/rybalka1/devmetrics/internal/config"
	"github.com/rybalka1/devmetrics/internal/server"
)

func main() {
	var (
		addr     string
		logLevel string
	)
	config.ServerArgsParse(&addr, &logLevel)
	log.Info().
		Str("addr", addr).
		Str("log", logLevel).Send()
	srv, err := server.NewServer(addr, logLevel)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	log.Fatal().Err(srv.Start()).Send()
}
