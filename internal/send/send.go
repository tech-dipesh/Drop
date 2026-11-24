package send

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"

	"github.com/tech-dipesh/drop/internal/transfer"
	"github.com/tech-dipesh/drop/internal/wire"
)

type Sender struct {
	Offered func(wire.Offer)
	Track   func(wire.Offer, *transfer.Meter) func()
}

func (s Sender) Send(ctx context.Context, addr string, offer wire.Offer, paths []string) error {
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	defer context.AfterFunc(ctx, func() { conn.Close() })()
	if err := wire.Send(conn, offer); err != nil {
		return fmt.Errorf("send offer: %w", err)
	}
	s.Offered(offer)
	var answer wire.Answer
	if err := wire.Receive(conn, &answer); err != nil {
		return fmt.Errorf("read answer: %w", err)
	}
	if !answer.Accept {
		return errors.New("the other side declined")
	}
	if err := s.stream(conn, paths, offer); err != nil {
		return err
	}
	var result wire.Result
	if err := wire.Receive(conn, &result); err != nil {
		return fmt.Errorf("read result: %w", err)
	}
	if !result.OK {
		return fmt.Errorf("receiver failed: %s", result.Error)
	}
	return nil
}

func (s Sender) stream(conn net.Conn, paths []string, offer wire.Offer) error {
	m := transfer.NewMeter(offer.Total())
	defer s.Track(offer, m)()
	for i, path := range paths {
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		err = transfer.Copy(conn, f, offer.Files[i].Size, m.Add)
		f.Close()
		if err != nil {
			return fmt.Errorf("send %s: %w", path, err)
		}
	}
	return nil
}

func Build(from string, paths []string) (wire.Offer, error) {
	offer := wire.Offer{Version: wire.Version, From: from}
	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			return offer, err
		}
		if info.IsDir() {
			return offer, fmt.Errorf("%s is a folder, folders aren't supported yet", path)
		}
		offer.Files = append(offer.Files, wire.File{Name: info.Name(), Size: info.Size()})
	}
	return offer, nil
}
