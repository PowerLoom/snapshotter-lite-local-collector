package relayers

import (
	"encoding/json"
	"fmt"
	"github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
	"os"
	"path/filepath"
	"runtime"
)

type Relayer struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Address         string `json:"address"`
	RendezvousPoint string `json:"rendezvousPoint"`
}

func GetTrustedRelayersWithAddresses(filename ...string) ([]peer.AddrInfo, error) {
	file, err := getFilePath("relayers.json", filename...)
	if err != nil {
		return nil, err
	}

	var relayers []Relayer

	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &relayers)
	if err != nil {
		return nil, err
	}

	var relayPeers []peer.AddrInfo
	for _, relayer := range relayers {
		maddr, err := ma.NewMultiaddr(relayer.Address)
		if err != nil {
			return nil, err
		}

		pid, err := peer.Decode(relayer.ID)
		if err != nil {
			return nil, err
		}

		relayPeers = append(relayPeers, peer.AddrInfo{
			ID:    pid,
			Addrs: []ma.Multiaddr{maddr},
		})
	}

	return relayPeers, nil
}

func GetTrustedRelayerIDs(rendezvousPoint string, filename ...string) ([]string, error) {
	file, err := getFilePath("relayers.json", filename...)
	if err != nil {
		return nil, err
	}

	var relayers []Relayer

	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &relayers)
	if err != nil {
		return nil, err
	}

	var relayIDs []string
	for _, relayer := range relayers {
		if relayer.RendezvousPoint == rendezvousPoint {
			relayIDs = append(relayIDs, relayer.ID)
		}
	}

	return relayIDs, nil
}

func getFilePath(defaultFilename string, filename ...string) (string, error) {
	if len(filename) > 0 {
		return filename[0], nil
	}

	_, currentFile, _, ok := runtime.Caller(1)
	if !ok {
		return "", fmt.Errorf("unable to get current file info")
	}
	dir := filepath.Dir(currentFile)
	return filepath.Join(dir, defaultFilename), nil
}

func IsRelayerTrusted(relayerID string, rendezvousPoint string) (bool, error) {
	trustedRelayerIDs, err := GetTrustedRelayerIDs(rendezvousPoint)
	if err != nil {
		return false, fmt.Errorf("failed to get trusted relayers: %v", err)
	}

	for _, id := range trustedRelayerIDs {
		if id == relayerID {
			return true, nil
		}
	}

	return false, nil
}
