package receive_test

import (
	"context"
	"os"
	"testing"

	"github.com/tech-dipesh/drop/internal/wire"
)

func TestSpeedTestSavesNothing(t *testing.T) {
	dst := t.TempDir()
	offer := wire.Offer{Version: wire.Version, From: "test", Files: []wire.File{{Name: "speed test", Size: 5 << 20}}, Test: true}
	addr, results := start(t, dst, true)
	if err := sender.Send(context.Background(), addr, offer, nil); err != nil {
		t.Fatal(err)
	}
	<-results
	if entries, _ := os.ReadDir(dst); len(entries) != 0 {
		t.Fatalf("speed test left %d files behind", len(entries))
	}
}
