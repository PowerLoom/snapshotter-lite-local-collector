package config

import (
	"encoding/json"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

var SettingsObj *Settings

type Settings struct {
	SequencerId            string `json:"SequencerId"`
	RelayerRendezvousPoint string `json:"RelayerRendezvousPoint"`
	ClientRendezvousPoint  string `json:"ClientRendezvousPoint"`
	RelayerPrivateKey      string `json:"RelayerPrivateKey"`
	PowerloomReportingUrl  string `json:"PowerloomReportingUrl"`
	SignerAccountAddress   string `json:"SignerAccountAddress"`
	PortNumber             string `json:"LocalCollectorPort"`
	TrustedRelayersListUrl string `json:"TrustedRelayersListUrl"`
	DataMarketAddress      string `json:"DataMarketAddress"`
}

func LoadConfig() {
	file, err := os.Open(strings.TrimSuffix(os.Getenv("CONFIG_PATH"), "/") + "/config/settings.json")
	if err != nil {
		log.Fatalf("Failed to open config file: %v", err)
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Errorf("Unable to close file: %s", err.Error())
		}
	}(file)

	decoder := json.NewDecoder(file)
	config := Settings{}
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalf("Failed to decode config file: %v", err)
	}

	if config.TrustedRelayersListUrl == "" {
		config.TrustedRelayersListUrl = "https://raw.githubusercontent.com/PowerLoom/snapshotter-lite-local-collector/feat/trusted-relayers/relayers.json"
	}

	SettingsObj = &config
}
