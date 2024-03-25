package config

import (
	"bufio"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
)

var SettingsObj *Settings

type Settings struct {
	RelayerUrl string `json:"RelayerUrl"`
	RelayerId  string `json:"RelayerId"`
}

func LoadConfig() {
	time.Sleep(10 * time.Second)
	file, err := os.Open(strings.TrimSuffix(os.Getenv("CONFIG_PATH"), "/") + "/config/settings.json")
	//file, err := os.Open("/Users/mukundrawat/powerloom/proto-snapshot-server/config/settings.json")
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

	file, err = os.Open("/shared_data/relayer_id.txt")
	if err != nil {
		log.Debugf("Error opening relayer info file: %v", err)
	}
	defer file.Close()

	// Initialize variables to hold relayer URL and ID

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		config.RelayerId = line
	}
	file, err = os.Open("/shared_data/relayer_url.txt")
	if err != nil {
		log.Debugf("Error opening relayer info file: %v", err)
	}
	defer file.Close()

	// Initialize variables to hold relayer URL and ID

	scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		config.RelayerUrl = line
	}
	SettingsObj = &config
}
