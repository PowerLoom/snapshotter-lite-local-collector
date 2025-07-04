package config

import (
	"os"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

var SettingsObj *Settings

type Settings struct {
	LogLevel               string
	SequencerID            string
	RelayerRendezvousPoint string
	ClientRendezvousPoint  string
	RelayerPrivateKey      string
	PowerloomReportingUrl  string
	SignerAccountAddress   string
	PortNumber             string
	TrustedRelayersListUrl string
	DataMarketAddress      string
	MaxStreamPoolSize      int
	DataMarketInRequest    bool

	// Stream Pool Configuration
	StreamHealthCheckTimeout time.Duration
	StreamWriteTimeout       time.Duration
	MaxWriteRetries          int
	MaxConcurrentWrites      int
	MaxStreamQueueSize       int
	WorkerPoolSize           int

	// Connection management settings
	ConnectionRefreshInterval time.Duration
}

func LoadConfig() {
	config := Settings{}

	// Required fields
	if port := os.Getenv("LOCAL_COLLECTOR_PORT"); port != "" {
		config.PortNumber = port
	} else {
		config.PortNumber = "50051" // Default value
	}

	if contract := os.Getenv("DATA_MARKET_CONTRACT"); contract == "" {
		log.Fatal("DATA_MARKET_CONTRACT environment variable is required")
	} else {
		config.DataMarketAddress = contract
	}
	if value := os.Getenv("DATA_MARKET_IN_REQUEST"); value == "true" {
		config.DataMarketInRequest = true
	} else {
		config.DataMarketInRequest = false
	}
	// Optional fields with defaults
	config.PowerloomReportingUrl = os.Getenv("POWERLOOM_REPORTING_URL")
	config.SignerAccountAddress = os.Getenv("SIGNER_ACCOUNT_ADDRESS")
	config.TrustedRelayersListUrl = getEnvWithDefault("TRUSTED_RELAYERS_LIST_URL", "https://raw.githubusercontent.com/PowerLoom/snapshotter-lite-local-collector/feat/trusted-relayers/relayers.json")

	// Load private key from file or env
	config.RelayerPrivateKey = loadPrivateKey()

	// Numeric values with defaults
	config.MaxStreamPoolSize = getEnvAsInt("MAX_STREAM_POOL_SIZE", 100)
	config.StreamHealthCheckTimeout = time.Duration(getEnvAsInt("STREAM_HEALTH_CHECK_TIMEOUT_MS", 5000)) * time.Millisecond
	config.StreamWriteTimeout = time.Duration(getEnvAsInt("STREAM_WRITE_TIMEOUT_MS", 5000)) * time.Millisecond
	config.MaxWriteRetries = getEnvAsInt("MAX_WRITE_RETRIES", 5)
	config.MaxConcurrentWrites = getEnvAsInt("MAX_CONCURRENT_WRITES", 100)
	config.MaxStreamQueueSize = getEnvAsInt("MAX_STREAM_QUEUE_SIZE", 1000)
	config.WorkerPoolSize = getEnvAsInt("WORKER_POOL_SIZE", 250)

	// Add log level setting (default "info")
	config.LogLevel = getEnvWithDefault("LOG_LEVEL", "info")

	// Add connection refresh interval setting (default 5 minutes)
	config.ConnectionRefreshInterval = time.Duration(getEnvAsInt("CONNECTION_REFRESH_INTERVAL_SEC", 300)) * time.Second

	SettingsObj = &config
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
		log.Warnf("Invalid value for %s, using default: %d", key, defaultValue)
	}
	return defaultValue
}

func loadPrivateKey() string {
	// Try loading from file first
	if keyBytes, err := os.ReadFile("/keys/key.txt"); err == nil {
		return string(keyBytes)
	}
	// Fall back to environment variable
	return os.Getenv("RELAYER_PRIVATE_KEY")
}
