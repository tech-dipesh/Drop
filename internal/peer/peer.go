package peer

import (
	"cmp"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
)

type Peer struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Addr string `json:"-"`
}

func Load() (Peer, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return Peer{}, err
	}
	data, err := os.ReadFile(filepath.Join(dir, "drop", "peer.json"))
	var p Peer
	if err == nil && json.Unmarshal(data, &p) == nil && p.ID != "" {
		return p, nil
	}
	id := make([]byte, 8)
	if _, err := rand.Read(id); err != nil {
		return Peer{}, err
	}
	host, _ := os.Hostname()
	p = Peer{ID: hex.EncodeToString(id), Name: cmp.Or(host, "drop device")}
	return p, Save(p)
}

func Save(p Peer) error {
	dir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	dir = filepath.Join(dir, "drop")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, _ := json.Marshal(p)
	return os.WriteFile(filepath.Join(dir, "peer.json"), data, 0o644)
}
