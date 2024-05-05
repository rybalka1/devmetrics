package service

import (
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/rybalka1/devmetrics/internal/config"
	"github.com/rybalka1/devmetrics/internal/logger"
	"github.com/rybalka1/devmetrics/internal/server/http"
	"github.com/rybalka1/devmetrics/internal/storage/filestorage"
	"github.com/rybalka1/devmetrics/internal/storage/memstorage"
)

type Service interface {
	Start() error
	Stop() error
	SetupLogger(level string) error
	memstorage.Storage
	filestorage.FileStorageInterface
}

type MetricService struct {
	server http.Server
	memstorage.Storage
	filestorage.FileStorageInterface
	StoreInterval time.Duration
	Restore       bool
}

func (service *MetricService) SetupLogger(level string) error {
	return logger.Initialize(level)
}

func NewService(args config.Args) (Service, error) {
	store := memstorage.NewMemStorage()
	fstore, err := filestorage.NewFileStorage(args.StoragePath, store)
	if err != nil {
		return nil, err
	}
	service := &MetricService{
		server:               nil,
		Storage:              store,
		FileStorageInterface: fstore,
		StoreInterval:        time.Duration(args.StoreInterval) * time.Second,
		Restore:              args.Restore,
	}
	server, err := http.NewServer(args)
	if err != nil {
		return nil, err
	}
	mux := NewRouter(service)
	server.AddMux(mux)
	service.server = server
	return service, nil
}

func (service MetricService) Start() error {
	if service.Restore {
		service.FileStorageInterface.Load()
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go service.BackupLoop(wg)
	err := service.server.Start()
	wg.Wait()
	return err
}

func (service *MetricService) BackupLoop(wg *sync.WaitGroup) error {
	timer := time.NewTicker(service.StoreInterval)
	for t := range timer.C {
		log.Info().Str("time", t.Format(time.DateTime)).Msg("DB saved")
		service.FileStorageInterface.SaveAll()
	}
	wg.Done()
	return nil
}

func (service MetricService) Stop() error {
	log.Error().Err(service.FileStorageInterface.SaveAll()).Msg("fs-storage saved")
	log.Error().Err(service.FileStorageInterface.Close()).Msg("fs-storage closed")
	return service.server.Stop()
}
