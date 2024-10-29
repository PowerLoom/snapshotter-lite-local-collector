package config

import (
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

var SettingsObj *Settings

type Settings struct {
	SequencerId            string
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
	config.TrustedRelayersListUrl = getEnvWithDefault("TRUSTED_RELAYERS_LIST_URL",
		"https://raw.githubusercontent.com/PowerLoom/snapshotter-lite-local-collector/feat/trusted-relayers/relayers.json")

	// Load private key from file or env
	config.RelayerPrivateKey = loadPrivateKey()

	// Numeric values with defaults
	config.MaxStreamPoolSize = getEnvAsInt("MAX_STREAM_POOL_SIZE", 2)

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
