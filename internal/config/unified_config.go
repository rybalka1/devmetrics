package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

// AgentConfig holds agent-specific configuration
type AgentConfig struct {
	Address        string
	PollInterval   int
	ReportInterval int
	Compression    bool
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Address       string
	LogLevel      string
	StoreInterval uint64
	StoragePath   string
	Restore       bool
	Profiling     bool
}

// UnifiedConfig holds configuration for both agent and server
type UnifiedConfig struct {
	// Common settings
	Address string

	// Agent-specific settings
	PollInterval   time.Duration
	ReportInterval time.Duration

	// Server-specific settings
	LogLevel      string
	StoreInterval time.Duration
	StoragePath   string
	Restore       bool

	// Feature flags
	EnableCompression bool
	EnableProfiling   bool
}

// Default configuration values
const (
	DefaultAddress         = "localhost:8080"
	DefaultPollInterval    = 2 * time.Second
	DefaultReportInterval  = 10 * time.Second
	DefaultLogLevel        = "info"
	DefaultStoreInterval   = 300 * time.Second
	DefaultStoragePath     = "/tmp/metrics-db.json"
	DefaultRestore         = true
	DefaultEnableCompression = false
	DefaultEnableProfiling   = false
)

// LoadUnifiedConfig loads configuration from environment variables and command line flags
func LoadUnifiedConfig() (*UnifiedConfig, error) {
	config := &UnifiedConfig{
		Address:           getEnvOrDefault("ADDRESS", DefaultAddress),
		PollInterval:      getDurationFromEnvOrDefault("POLL_INTERVAL", DefaultPollInterval),
		ReportInterval:    getDurationFromEnvOrDefault("REPORT_INTERVAL", DefaultReportInterval),
		LogLevel:          getEnvOrDefault("LOG_LEVEL", DefaultLogLevel),
		StoreInterval:     getDurationFromEnvOrDefault("STORE_INTERVAL", DefaultStoreInterval),
		StoragePath:       getEnvOrDefault("FILE_STORAGE_PATH", DefaultStoragePath),
		Restore:           getBoolFromEnvOrDefault("RESTORE", DefaultRestore),
		EnableCompression: getBoolFromEnvOrDefault("ENABLE_COMPRESSION", DefaultEnableCompression),
		EnableProfiling:   getBoolFromEnvOrDefault("ENABLE_PROFILING", DefaultEnableProfiling),
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	log.Info().
		Str("address", config.Address).
		Dur("poll_interval", config.PollInterval).
		Dur("report_interval", config.ReportInterval).
		Str("log_level", config.LogLevel).
		Dur("store_interval", config.StoreInterval).
		Str("storage_path", config.StoragePath).
		Bool("restore", config.Restore).
		Bool("compression", config.EnableCompression).
		Bool("profiling", config.EnableProfiling).
		Msg("Configuration loaded")

	return config, nil
}

// Validate checks if the configuration is valid
func (c *UnifiedConfig) Validate() error {
	if c.Address == "" {
		return fmt.Errorf("address cannot be empty")
	}

	if c.PollInterval <= 0 {
		return fmt.Errorf("poll interval must be positive")
	}

	if c.ReportInterval <= 0 {
		return fmt.Errorf("report interval must be positive")
	}

	if c.StoreInterval <= 0 {
		return fmt.Errorf("store interval must be positive")
	}

	if c.StoragePath == "" {
		return fmt.Errorf("storage path cannot be empty")
	}

	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLogLevels[c.LogLevel] {
		return fmt.Errorf("invalid log level: %s", c.LogLevel)
	}

	return nil
}

// GetAgentConfig returns agent-specific configuration
func (c *UnifiedConfig) GetAgentConfig() AgentConfig {
	return AgentConfig{
		Address:        c.Address,
		PollInterval:   int(c.PollInterval.Seconds()),
		ReportInterval: int(c.ReportInterval.Seconds()),
		Compression:    c.EnableCompression,
	}
}

// GetServerConfig returns server-specific configuration
func (c *UnifiedConfig) GetServerConfig() ServerConfig {
	return ServerConfig{
		Address:       c.Address,
		LogLevel:      c.LogLevel,
		StoreInterval: uint64(c.StoreInterval.Seconds()),
		StoragePath:   c.StoragePath,
		Restore:       c.Restore,
		Profiling:     c.EnableProfiling,
	}
}

// Helper functions for environment variable parsing
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getDurationFromEnvOrDefault(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return time.Duration(intValue) * time.Second
		} else {
			log.Warn().
				Str("key", key).
				Str("value", value).
				Err(err).
				Dur("default", defaultValue).
				Msg("Failed to parse duration from environment, using default")
		}
	}
	return defaultValue
}

func getBoolFromEnvOrDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		} else {
			log.Warn().
				Str("key", key).
				Str("value", value).
				Err(err).
				Bool("default", defaultValue).
				Msg("Failed to parse bool from environment, using default")
		}
	}
	return defaultValue
}
