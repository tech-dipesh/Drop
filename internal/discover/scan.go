package discover

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"strconv"
	"time"

	"github.com/tech-dipesh/drop/internal/peer"
	"github.com/tech-dipesh/drop/internal/wire"
)

func Scan(ctx context.Context) (<-chan peer.Peer, error) {
	nets, err := networks()
	if err != nil {
		return nil, err
	}
	if len(nets) == 0 {
		return nil, errors.New("no usable network found, connect to wifi or a hotspot first")
	}
	conn, err := net.ListenUDP("udp4", nil)
	if err != nil {
		return nil, err
	}
	context.AfterFunc(ctx, func() { conn.Close() })
	found := make(chan peer.Peer)
	go probe(ctx, conn, nets)
	go collect(ctx, conn, found)
	return found, nil
}

func probe(ctx context.Context, conn *net.UDPConn, nets []*net.IPNet) {
	msg, _ := json.Marshal(beacon{Version: wire.Version, Kind: "probe"})
	tick := time.NewTicker(time.Second)
	defer tick.Stop()
	for {
		for _, n := range nets {
			conn.WriteToUDP(msg, &net.UDPAddr{IP: broadcast(n), Port: wire.Port})
		}
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
		}
	}
}

func collect(ctx context.Context, conn *net.UDPConn, found chan<- peer.Peer) {
	seen := map[string]bool{}
	buf := make([]byte, 1024)
	for {
		n, from, err := conn.ReadFromUDP(buf)
		if err != nil {
			return
		}
		p, ok := parse(buf[:n], from)
		if !ok || seen[p.ID] {
			continue
		}
		seen[p.ID] = true
		select {
		case found <- p:
		case <-ctx.Done():
			return
		}
	}
}

func parse(b []byte, from *net.UDPAddr) (peer.Peer, bool) {
	var msg beacon
	if json.Unmarshal(b, &msg) != nil || msg.Kind != "here" || msg.Version != wire.Version || msg.ID == "" {
		return peer.Peer{}, false
	}
	addr := net.JoinHostPort(from.IP.String(), strconv.Itoa(wire.Port))
	return peer.Peer{ID: msg.ID, Name: msg.Name, Addr: addr}, true
}
