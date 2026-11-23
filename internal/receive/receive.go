package receive

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/tech-dipesh/drop/internal/transfer"
	"github.com/tech-dipesh/drop/internal/wire"
)

type Receiver struct {
	Dir   string
	Ask   func(wire.Offer) bool
	Track func(wire.Offer, *transfer.Meter) func()
	Done  func(saved []string, err error)
}

func (r Receiver) Serve(ctx context.Context, ln net.Listener) error {
	defer context.AfterFunc(ctx, func() { ln.Close() })()
	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		saved, err := r.handle(ctx, conn)
		r.Done(saved, err)
	}
}

func (r Receiver) handle(ctx context.Context, conn net.Conn) ([]string, error) {
	defer conn.Close()
	defer context.AfterFunc(ctx, func() { conn.Close() })()
	var offer wire.Offer
	if err := wire.Receive(conn, &offer); err != nil {
		return nil, fmt.Errorf("read offer: %w", err)
	}
	if err := check(offer); err != nil {
		return nil, err
	}
	accept := r.Ask(offer)
	if err := wire.Send(conn, wire.Answer{Accept: accept}); err != nil || !accept {
		return nil, err
	}
	m := transfer.NewMeter(offer.Total())
	defer r.Track(offer, m)()
	var saved []string
	for _, f := range offer.Files {
		path, err := save(r.Dir, f, conn, m.Add)
		if err != nil {
			wire.Send(conn, wire.Result{Error: err.Error()})
			return saved, fmt.Errorf("save %s: %w", f.Name, err)
		}
		saved = append(saved, path)
	}
	return saved, wire.Send(conn, wire.Result{OK: true})
}

func check(o wire.Offer) error {
	if o.Version != wire.Version {
		return fmt.Errorf("sender uses protocol v%d, this drop uses v%d", o.Version, wire.Version)
	}
	if len(o.Files) == 0 {
		return errors.New("offer has no files")
	}
	for _, f := range o.Files {
		if !validName(f.Name) || f.Size < 0 {
			return fmt.Errorf("offer has an unsafe file entry %q", f.Name)
		}
	}
	return nil
}
