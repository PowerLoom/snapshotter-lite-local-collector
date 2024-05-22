package relayers

import (
	"context"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	ma "github.com/multiformats/go-multiaddr"
	"testing"
	"time"
)

func TestReadRelayersFromJSON(t *testing.T) {
	t.Log("Starting TestReadRelayersFromJSON")

	relayers, err := GetTrustedRelayers("fixtures/relayers_test.json")
	if err != nil {
		t.Fatalf("Failed to read relayers: %v", err)
	}

	t.Logf("Read %d relayers", len(relayers))

	if len(relayers) != 2 {
		t.Fatalf("Expected 2 relayers, got %d", len(relayers))
	}

	expectedID1, _ := peer.Decode("QmY7Yh4UquoXHLPFo2XbhXkhBvFoPwmQUSa92pxnxjQuPU")
	expectedAddr1, _ := ma.NewMultiaddr("/ip4/127.0.0.1/tcp/4001/p2p/QmY7Yh4UquoXHLPFo2XbhXkhBvFoPwmQUSa92pxnxjQuPU")

	t.Logf("Checking relayer 1: ID = %s, Address = %s", relayers[0].ID, relayers[0].Addrs[0])

	if relayers[0].ID != expectedID1 || !relayers[0].Addrs[0].Equal(expectedAddr1) {
		t.Fatalf("Unexpected relayer data: %v", relayers[0])
	}

	expectedID2, _ := peer.Decode("QmY7Yh4UquoXHLPFo2XbhXkhBvFoPwmQUSa92pxnxjQuPV")
	expectedAddr2, _ := ma.NewMultiaddr("/ip4/127.0.0.1/tcp/4002/p2p/QmY7Yh4UquoXHLPFo2XbhXkhBvFoPwmQUSa92pxnxjQuPV")

	t.Logf("Checking relayer 2: ID = %s, Address = %s", relayers[1].ID, relayers[1].Addrs[0])

	if relayers[1].ID != expectedID2 || !relayers[1].Addrs[0].Equal(expectedAddr2) {
		t.Fatalf("Unexpected relayer data: %v", relayers[1])
	}

	t.Log("TestReadRelayersFromJSON completed successfully")
}

func BenchmarkRelayConnection(b *testing.B) {
	relayers, err := GetTrustedRelayers("./relayers.json")
	if err != nil {
		b.Fatalf("Failed to read relayers: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	p2pClient, err := libp2p.New()
	if err != nil {
		b.Fatalf("Failed to create host: %v", err)
	}

	b.ResetTimer()

	for _, relay := range relayers {
		b.Run(relay.ID.String(), func(b *testing.B) {
			err := attemptConnect(ctx, p2pClient, relay)
			if err != nil {
				b.Errorf("Failed to connect to relay: %v", err)
			}
		})
	}
}

func attemptConnect(ctx context.Context, h host.Host, relayInfo peer.AddrInfo) error {
	h.Peerstore().AddAddrs(relayInfo.ID, relayInfo.Addrs, peerstore.AddressTTL)
	err := h.Connect(ctx, relayInfo)
	if err != nil {
		return err
	}
	return nil
}
