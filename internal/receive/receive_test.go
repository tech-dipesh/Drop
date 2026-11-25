package receive_test

import (
	"bytes"
	"context"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tech-dipesh/drop/internal/receive"
	"github.com/tech-dipesh/drop/internal/send"
	"github.com/tech-dipesh/drop/internal/transfer"
	"github.com/tech-dipesh/drop/internal/wire"
)

func noTrack(wire.Offer, *transfer.Meter) func() { return func() {} }

var sender = send.Sender{Paired: func(string) {}, Offered: func(wire.Offer) {}, Track: noTrack}

func start(t *testing.T, dir string, accept bool) (string, chan []string) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	results := make(chan []string, 1)
	r := receive.Receiver{
		Dir:    dir,
		Paired: func(string) {},
		Ask:    func(wire.Offer) bool { return accept },
		Track:  noTrack,
		Done: func(saved []string, err error) {
			if err != nil {
				t.Error(err)
			}
			results <- saved
		},
	}
	go r.Serve(ctx, ln)
	return ln.Addr().String(), results
}

func TestTransferAccepted(t *testing.T) {
	src, dst := t.TempDir(), t.TempDir()
	want := map[string][]byte{"a.txt": []byte("hello"), "b.bin": bytes.Repeat([]byte{7}, 5<<20)}
	var paths []string
	for name, data := range want {
		path := filepath.Join(src, name)
		if err := os.WriteFile(path, data, 0o644); err != nil {
			t.Fatal(err)
		}
		paths = append(paths, path)
	}
	offer, err := send.Build("test", paths)
	if err != nil {
		t.Fatal(err)
	}
	addr, results := start(t, dst, true)
	if err := sender.Send(context.Background(), addr, offer, paths); err != nil {
		t.Fatal(err)
	}
	saved := <-results
	if len(saved) != len(want) {
		t.Fatalf("saved %d files, want %d", len(saved), len(want))
	}
	for _, path := range saved {
		got, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(got, want[filepath.Base(path)]) {
			t.Fatalf("%s doesn't match what was sent", path)
		}
	}
}

func TestTransferDeclined(t *testing.T) {
	src, dst := t.TempDir(), t.TempDir()
	path := filepath.Join(src, "a.txt")
	if err := os.WriteFile(path, []byte("hi"), 0o644); err != nil {
		t.Fatal(err)
	}
	offer, err := send.Build("test", []string{path})
	if err != nil {
		t.Fatal(err)
	}
	addr, results := start(t, dst, false)
	err = sender.Send(context.Background(), addr, offer, []string{path})
	if err == nil || !strings.Contains(err.Error(), "declined") {
		t.Fatalf("got %v, want declined", err)
	}
	if saved := <-results; len(saved) != 0 {
		t.Fatalf("saved %v after declining", saved)
	}
	if entries, _ := os.ReadDir(dst); len(entries) != 0 {
		t.Fatalf("found %d files after declining", len(entries))
	}
}
