package discover

import (
	"net"
	"strconv"
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

func TestHosts(t *testing.T) {
	cases := map[string][3]string{
		"192.168.1.37/24": {"254", "192.168.1.1", "192.168.1.254"},
		"172.20.10.2/28":  {"14", "172.20.10.1", "172.20.10.14"},
		"10.0.5.9/16":     {"254", "10.0.5.1", "10.0.5.254"},
	}
	for cidr, want := range cases {
		ip, n, err := net.ParseCIDR(cidr)
		if err != nil {
			t.Fatal(err)
		}
		ips := hosts(&net.IPNet{IP: ip.To4(), Mask: n.Mask})
		got := [3]string{strconv.Itoa(len(ips)), ips[0].String(), ips[len(ips)-1].String()}
		if got != want {
			t.Errorf("%s: got %v, want %v", cidr, got, want)
		}
	}
}
