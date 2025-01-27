package service

import (
	"context"
	"fmt"
	"proto-snapshot-server/config"
	"sync"
	"time"

	circuitv2 "github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/client"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	rcmgr "github.com/libp2p/go-libp2p/p2p/host/resource-manager"
	"github.com/libp2p/go-libp2p/p2p/muxer/yamux"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
	"github.com/libp2p/go-libp2p/p2p/security/noise"
	libp2ptls "github.com/libp2p/go-libp2p/p2p/security/tls"
	ma "github.com/multiformats/go-multiaddr"
	log "github.com/sirupsen/logrus"
)

var (
	SequencerHostConn host.Host
	SequencerID       peer.ID
	sequencerMu       sync.RWMutex
	ConnManager       *connmgr.BasicConnMgr
	TcpAddr           ma.Multiaddr
	rm                network.ResourceManager
)

// Thread-safe getter for connection state
func GetSequencerConnection() (host.Host, peer.ID, error) {
	sequencerMu.RLock()
	defer sequencerMu.RUnlock()

	if SequencerHostConn == nil || SequencerID.String() == "" {
		return nil, "", fmt.Errorf("sequencer connection not established")
	}

	return SequencerHostConn, SequencerID, nil
}

func ConnectToSequencerP2P(relayers []Relayer, p2pHost host.Host) bool {
	for _, relayer := range relayers {
		relayerMA, _ := ma.NewMultiaddr(relayer.Maddr)
		relayerInfo, _ := peer.AddrInfoFromP2pAddr(relayerMA)

		if reservation, err := circuitv2.Reserve(context.Background(), p2pHost, *relayerInfo); err != nil {
			log.Fatalf("Failed to request reservation with relay: %v", err)
		} else {
			fmt.Println("Reservation with relay successful", reservation.Expiration, reservation.LimitDuration)
		}

		sequencerAddr, err := ma.NewMultiaddr(fmt.Sprintf("%s/p2p-circuit/p2p/%s", relayer.Maddr, config.SettingsObj.SequencerID))
		if err != nil {
			log.Debugln(err.Error())
		}
		log.Debugln("Connecting to Sequencer: ", sequencerAddr.String())

		isConnected := AddPeerConnection(context.Background(), p2pHost, sequencerAddr.String())
		if isConnected {
			return true
		}
	}

	return false
}

func CreateLibP2pHost() error {
	var err error
	TcpAddr, _ = ma.NewMultiaddr("/ip4/0.0.0.0/tcp/9000")

	ConnManager, _ = connmgr.NewConnManager(
		40960,
		81920,
		connmgr.WithGracePeriod(1*time.Minute))

	scalingLimits := rcmgr.DefaultLimits

	libp2p.SetDefaultServiceLimits(&scalingLimits)

	scaledDefaultLimits := scalingLimits.AutoScale()

	cfg := rcmgr.PartialLimitConfig{
		System: rcmgr.ResourceLimits{
			StreamsOutbound: rcmgr.Unlimited,
			StreamsInbound:  rcmgr.Unlimited,
			Streams:         rcmgr.Unlimited,
			Conns:           rcmgr.Unlimited,
			ConnsOutbound:   rcmgr.Unlimited,
			ConnsInbound:    rcmgr.Unlimited,
			FD:              rcmgr.Unlimited,
			Memory:          rcmgr.LimitVal64(rcmgr.Unlimited),
		},
	}

	limits := cfg.Build(scaledDefaultLimits)

	limiter := rcmgr.NewFixedLimiter(limits)

	rm, err = rcmgr.NewResourceManager(limiter, rcmgr.WithMetricsDisabled())

	if err != nil {
		log.Debugln("Error instantiating resource manager: ", err.Error())
		return err
	}

	SequencerHostConn, err = libp2p.New(
		libp2p.EnableRelay(),
		libp2p.ConnectionManager(ConnManager),
		libp2p.ListenAddrs(TcpAddr),
		libp2p.ResourceManager(rm),
		libp2p.Security(libp2ptls.ID, libp2ptls.New),
		libp2p.Security(noise.ID, noise.New),
		libp2p.DefaultTransports,
		libp2p.NATPortMap(),
		libp2p.EnableRelayService(),
		libp2p.EnableNATService(),
		libp2p.EnableHolePunching(),
		libp2p.Muxer(yamux.ID, yamux.DefaultTransport))

	if err != nil {
		log.Debugln("Error instantiating libp2p host: ", err.Error())
		return err
	}
	return nil
}

// EstablishSequencerConnection should only be called during initialization
// or explicit reconnection logic, not during stream operations
func EstablishSequencerConnection() error {
	sequencerMu.Lock()
	defer sequencerMu.Unlock()

	// Clear existing connection if any
	if SequencerHostConn != nil {
		if err := SequencerHostConn.Close(); err != nil {
			log.Warnf("Error closing existing connection: %v", err)
		}
		// Important: Signal that connection is being reset
		// This should trigger cleanup of existing stream pool
		SequencerHostConn = nil
		SequencerID = ""
	}

	// 1. Create properly configured host
	if err := CreateLibP2pHost(); err != nil {
		return fmt.Errorf("failed to create libp2p host: %w", err)
	}

	// No need to reassign SequencerHostConn as it's already set in CreateLibP2pHost()

	// 2. Get sequencer info
	sequencer, err := fetchSequencer(
		"https://raw.githubusercontent.com/PowerLoom/snapshotter-lite-local-collector/feat/trusted-relayers/sequencers.json",
		config.SettingsObj.DataMarketAddress,
	)
	if err != nil {
		return fmt.Errorf("failed to fetch sequencer info: %w", err)
	}

	// 3. Parse multiaddr and create peer info
	maddr, err := ma.NewMultiaddr(sequencer.Maddr)
	if err != nil {
		return fmt.Errorf("failed to parse multiaddr: %w", err)
	}

	sequencerInfo, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		return fmt.Errorf("failed to get addr info: %w", err)
	}

	// 4. Set sequencer ID
	SequencerID = sequencerInfo.ID
	if SequencerID.String() == "" {
		return fmt.Errorf("empty sequencer ID")
	}

	// 5. Establish connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := SequencerHostConn.Connect(ctx, *sequencerInfo); err != nil {
		return fmt.Errorf("failed to connect to sequencer: %w", err)
	}

	log.Infof("Successfully connected to Sequencer: %s with ID: %s", sequencer.Maddr, SequencerID.String())
	return nil
}
