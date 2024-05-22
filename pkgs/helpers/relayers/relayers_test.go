package relayers

import (
	"context"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"testing"
	"time"
)

func TestReadRelayersFromJSON(t *testing.T) {
	t.Log("Starting TestReadRelayersFromJSON")

	relayers, err := GetTrustedRelayersWithAddresses("fixtures/relayers_test.json")
	if err != nil {
		t.Fatalf("Failed to read relayers: %v", err)
	}

	t.Logf("Read %d relayers", len(relayers))

	if len(relayers) != 2 {
		t.Fatalf("Expected 2 relayers, got %d", len(relayers))
	}

	expectedID1, _ := peer.Decode("QmY7Yh4UquoXHLPFo2XbhXkhBvFoPwmQUSa92pxnxjQuPU")
	//expectedAddr1, _ := ma.NewMultiaddr("/ip4/127.0.0.1/tcp/4001/p2p/QmY7Yh4UquoXHLPFo2XbhXkhBvFoPwmQUSa92pxnxjQuPU")

	t.Logf("Checking relayer 1: ID = %s", relayers[0].ID)

	if relayers[0].ID != expectedID1 {
		t.Fatalf("Unexpected relayer data: %v", relayers[0])
	}

	//if relayers[0].ID != expectedID1 || !relayers[0].Addrs[0].Equal(expectedAddr1) {
	//	t.Fatalf("Unexpected relayer data: %v", relayers[0])
	//}

	expectedID2, _ := peer.Decode("QmY7Yh4UquoXHLPFo2XbhXkhBvFoPwmQUSa92pxnxjQuPV")
	//expectedAddr2, _ := ma.NewMultiaddr("/ip4/127.0.0.1/tcp/4002/p2p/QmY7Yh4UquoXHLPFo2XbhXkhBvFoPwmQUSa92pxnxjQuPV")

	t.Logf("Checking relayer 2: ID = %s", relayers[1].ID)

	if relayers[1].ID != expectedID2 {
		t.Fatalf("Unexpected relayer data: %v", relayers[1])
	}

	//if relayers[1].ID != expectedID2 || !relayers[1].Addrs[0].Equal(expectedAddr2) {
	//	t.Fatalf("Unexpected relayer data: %v", relayers[1])
	//}

	t.Log("TestReadRelayersFromJSON completed successfully")
}

func TestGetTrustedRelayerIDs(t *testing.T) {
	t.Log("Starting TestGetTrustedRelayerIDs")

	rendezvousPoint := "Relayer_POP_test_simulation_phase_1"
	relayerIDs, err := GetTrustedRelayerIDs(rendezvousPoint, "fixtures/relayers_no_address.json")
	t.Logf("Reading relayers with rendezvous point %s", rendezvousPoint)
	if err != nil {
		t.Fatalf("Failed to read relayer IDs: %v", err)
	}

	t.Logf("Read %d relayer IDs", len(relayerIDs))

	if len(relayerIDs) != 2 {
		t.Fatalf("Expected 2 relayer IDs, got %d", len(relayerIDs))
	}

	expectedID1 := "QmaNQc1MbGzrwbmUytPfikkrJeJ5sZF24ZZa47RGRnYoo8"
	expectedID2 := "QmWG6hesiCV5uFJwPaThtSrfnwEnwQovuXHJB6j6j7p4TL"
	notExpectedID := "QmY7Yh4UquoXHLPFo2XbhXkhBvFoPwmQUSa92pxnxjQuPV"

	if relayerIDs[0] != expectedID1 || relayerIDs[1] != expectedID2 {
		t.Fatalf("Unexpected relayer IDs: %v", relayerIDs)
	}

	for _, id := range relayerIDs {
		if id == notExpectedID {
			t.Fatalf("Relayer with ID %s should not be included", notExpectedID)
		}
	}

	t.Log("TestGetTrustedRelayerIDs completed successfully")
}

func TestIsRelayerTrusted(t *testing.T) {
	// Define the rendezvous point and the relayer ID
	rendezvousPoint := "Relayer_POP_test_simulation_phase_1"
	trustedRelayerID := "QmaNQc1MbGzrwbmUytPfikkrJeJ5sZF24ZZa47RGRnYoo8"
	untrustedRelayerID := "QmY7Yh4UquoXHLPFo2XbhXkhBvFoPwmQUSa92pxnxjQuPV"

	t.Logf("Call the function with the trusted relayer ID")
	isTrusted, err := IsRelayerTrusted(trustedRelayerID, rendezvousPoint)
	if err != nil {
		t.Fatalf("Failed to check if relayer is trusted: %v", err)
	}

	t.Logf("Check if the relayer is correctly identified as trusted")
	if !isTrusted {
		t.Fatalf("Expected relayer with ID %s to be trusted", trustedRelayerID)
	}
	t.Logf("Relayer with ID %s is trusted", trustedRelayerID)

	t.Logf("Call the function with the untrusted relayer ID")
	isTrusted, err = IsRelayerTrusted(untrustedRelayerID, rendezvousPoint)
	if err != nil {
		t.Fatalf("Failed to check if relayer is trusted: %v", err)
	}

	t.Logf("Check if the relayer is correctly identified as untrusted")
	if isTrusted {
		t.Fatalf("Expected relayer with ID %s to be untrusted", untrustedRelayerID)
	}
	t.Logf("Relayer with ID %s is untrusted", untrustedRelayerID)
}

func BenchmarkRelayConnection(b *testing.B) {
	relayers, err := GetTrustedRelayersWithAddresses("./relayers.json")
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
