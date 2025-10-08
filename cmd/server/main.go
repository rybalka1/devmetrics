package main

import (
	"github.com/rs/zerolog/log"
	"github.com/rybalka1/devmetrics/internal/config"
	"github.com/rybalka1/devmetrics/internal/service"
)

func main() {
	cfg, err := config.LoadUnifiedConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}
	serverConfig := cfg.GetServerConfig()
	log.Info().Str("addr", serverConfig.Address).Str("log", serverConfig.LogLevel).Send()
	Service, err := service.NewService(serverConfig)

	if err != nil {
		log.Fatal().Err(err).Send()
	}
	log.Fatal().Err(Service.Start()).Send()
}
