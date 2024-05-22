package relayers

import (
	"encoding/json"
	"github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
	"os"
	"path/filepath"
)

type Relayer struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

func GetTrustedRelayers(filename ...string) ([]peer.AddrInfo, error) {
	var file string
	if len(filename) > 0 {
		file = filename[0]
	} else {
		// Use the default JSON file in your repository
		file = filepath.Join(".", "relayers.json")
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
