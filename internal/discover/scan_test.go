package discover

import (
	"net"
	"testing"
)

func TestParse(t *testing.T) {
	from := &net.UDPAddr{IP: net.ParseIP("192.168.1.20"), Port: 51000}
	p, ok := parse([]byte(`{"version":1,"kind":"here","id":"ab12","name":"Dipesh-laptop"}`), from)
	if !ok || p.ID != "ab12" || p.Name != "Dipesh-laptop" || p.Addr != "192.168.1.20:47470" {
		t.Fatalf("got %+v, %v", p, ok)
	}
	bad := []string{
		`{"version":1,"kind":"probe"}`,
		`{"version":2,"kind":"here","id":"ab12"}`,
		`{"version":1,"kind":"here"}`,
		`not json`,
	}
	for _, msg := range bad {
		if _, ok := parse([]byte(msg), from); ok {
			t.Errorf("%s should be ignored", msg)
		}
	}
}
