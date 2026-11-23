package discover

import (
	"context"
	"encoding/json"
	"net"

	"github.com/tech-dipesh/drop/internal/peer"
	"github.com/tech-dipesh/drop/internal/wire"
)

type beacon struct {
	Version int    `json:"version"`
	Kind    string `json:"kind"`
	ID      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
}

func Answer(ctx context.Context, me peer.Peer) error {
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{Port: wire.Port})
	if err != nil {
		return err
	}
	defer conn.Close()
	defer context.AfterFunc(ctx, func() { conn.Close() })()
	reply, _ := json.Marshal(beacon{Version: wire.Version, Kind: "here", ID: me.ID, Name: me.Name})
	buf := make([]byte, 1024)
	for {
		n, from, err := conn.ReadFromUDP(buf)
		if err != nil {
			return err
		}
		var b beacon
		if json.Unmarshal(buf[:n], &b) == nil && b.Kind == "probe" && b.Version == wire.Version {
			conn.WriteToUDP(reply, from)
		}
	}
}
