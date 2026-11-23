package discover

import (
	"net"
	"testing"
)

func TestBroadcast(t *testing.T) {
	cases := map[string]string{
		"192.168.1.37/24": "192.168.1.255",
		"10.0.5.9/16":     "10.0.255.255",
		"172.20.10.2/28":  "172.20.10.15",
	}
	for cidr, want := range cases {
		ip, n, err := net.ParseCIDR(cidr)
		if err != nil {
			t.Fatal(err)
		}
		got := broadcast(&net.IPNet{IP: ip.To4(), Mask: n.Mask})
		if got.String() != want {
			t.Errorf("%s: got %s, want %s", cidr, got, want)
		}
	}
}

func TestNetworksAreIPv4(t *testing.T) {
	nets, err := networks()
	if err != nil {
		t.Fatal(err)
	}
	for _, n := range nets {
		if len(n.IP) != 4 || len(n.Mask) != 4 {
			t.Errorf("%v isn't a 4 byte ipv4 network", n)
		}
	}
}
