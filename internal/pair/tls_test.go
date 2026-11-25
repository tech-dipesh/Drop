package pair

import (
	"crypto/tls"
	"net"
	"testing"
)

func pipe(t *testing.T) (*tls.Conn, *tls.Conn) {
	t.Helper()
	conf, err := ServerConfig()
	if err != nil {
		t.Fatal(err)
	}
	a, b := net.Pipe()
	client := tls.Client(a, ClientConfig())
	server := tls.Server(b, conf)
	t.Cleanup(func() {
		a.Close()
		b.Close()
	})
	errs := make(chan error, 1)
	go func() { errs <- server.Handshake() }()
	if err := client.Handshake(); err != nil {
		t.Fatal(err)
	}
	if err := <-errs; err != nil {
		t.Fatal(err)
	}
	return client, server
}

func TestHandshakeUsesTLS13(t *testing.T) {
	client, _ := pipe(t)
	if v := client.ConnectionState().Version; v != tls.VersionTLS13 {
		t.Fatalf("got tls version %x, want 1.3", v)
	}
}
