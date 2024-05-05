package config

import (
	"flag"
	"os"
	"strconv"
)

type Args struct {
	Addr          string
	LogLevel      string
	StoreInterval uint64
	StoragePath   string
	Restore       bool
}

func ServerArgsParse(args *Args) {
	args.Addr = os.Getenv("ADDRESS")
	args.LogLevel = os.Getenv("LOGLEVEL")
	args.StoreInterval, _ = strconv.ParseUint(os.Getenv("STORE_INTERVAL"), 10, 64)
	args.StoragePath = os.Getenv("FILE_STORAGE_PATH")
	args.Restore, _ = strconv.ParseBool(os.Getenv("RESTORE"))

	var (
		flagAddr          string
		flagLogLevel      string
		flagStoreInterval uint64
		flagStoragePath   string
		flagRestore       bool
	)

	flag.StringVar(&flagAddr, "a", "localhost:8080", "address for server")
	flag.StringVar(&flagLogLevel, "l", "info", "log level")
	flag.Uint64Var(&flagStoreInterval, "i", 300, "interval in seconds for write store data to file")
	flag.StringVar(&flagStoragePath, "f", "/tmp/metrics-db.json", "file for saving metrics")
	flag.BoolVar(&flagRestore, "r", true, "restore from file-db on startup")
	flag.Parse()
	if args.Addr == "" {
		args.Addr = flagAddr
	}
	if args.LogLevel == "" {
		args.LogLevel = flagLogLevel
	}
	if args.StoreInterval == 0 {
		args.StoreInterval = flagStoreInterval
	}
	if args.StoragePath == "" {
		args.StoragePath = flagStoragePath
	}
	if !args.Restore {
		args.Restore = flagRestore
	}
}
