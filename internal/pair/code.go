package pair

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"encoding/binary"
	"errors"
	"fmt"
	"slices"

	"github.com/tech-dipesh/drop/internal/wire"
)

func AsSender(conn *tls.Conn) (string, error) {
	mine := make([]byte, 32)
	if _, err := rand.Read(mine); err != nil {
		return "", err
	}
	sum := sha256.Sum256(mine)
	if err := wire.Send(conn, wire.Commit{Hash: sum[:]}); err != nil {
		return "", err
	}
	var theirs wire.Nonce
	if err := wire.Receive(conn, &theirs); err != nil {
		return "", err
	}
	if err := wire.Send(conn, wire.Reveal{Nonce: mine}); err != nil {
		return "", err
	}
	return code(conn, mine, theirs.Value)
}

func AsReceiver(conn *tls.Conn) (string, error) {
	var commit wire.Commit
	if err := wire.Receive(conn, &commit); err != nil {
		return "", err
	}
	mine := make([]byte, 32)
	if _, err := rand.Read(mine); err != nil {
		return "", err
	}
	if err := wire.Send(conn, wire.Nonce{Value: mine}); err != nil {
		return "", err
	}
	var reveal wire.Reveal
	if err := wire.Receive(conn, &reveal); err != nil {
		return "", err
	}
	sum := sha256.Sum256(reveal.Nonce)
	if !bytes.Equal(sum[:], commit.Hash) {
		return "", errors.New("sender's nonce doesn't match what it committed to")
	}
	return code(conn, reveal.Nonce, mine)
}

func code(conn *tls.Conn, sender, receiver []byte) (string, error) {
	state := conn.ConnectionState()
	ekm, err := state.ExportKeyingMaterial("drop pairing", nil, 32)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(slices.Concat(ekm, sender, receiver))
	n := binary.BigEndian.Uint64(sum[:8]) % 1_000_000
	return fmt.Sprintf("%03d %03d", n/1000, n%1000), nil
}
