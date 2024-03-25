package service

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
	log "github.com/sirupsen/logrus"
	"proto-snapshot-server/config"
)

var rpctorelay, _ = libp2p.New(libp2p.EnableRelay())
var RelayerId peer.ID

func ConfigureRelayer() {
	ctx := context.Background()

	var relayAddr ma.Multiaddr
	var err error
	relayAddr, err = ma.NewMultiaddr(fmt.Sprintf("%s/p2p/%s", config.SettingsObj.RelayerUrl, config.SettingsObj.RelayerId))
	if err != nil {
		log.Debugln(err.Error())
		return
	}

	relayerinfo, err := peer.AddrInfoFromP2pAddr(relayAddr)
	RelayerId = relayerinfo.ID
	log.Debugln(err)

	//Establish connections
	if err = rpctorelay.Connect(ctx, *relayerinfo); err != nil {
		log.Debugln("Failed to connect grpc server to relayer")
	}
}
