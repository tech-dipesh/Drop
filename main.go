package main

import (
	"cmp"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"

	"github.com/tech-dipesh/drop/internal/discover"
	"github.com/tech-dipesh/drop/internal/peer"
	"github.com/tech-dipesh/drop/internal/receive"
	"github.com/tech-dipesh/drop/internal/send"
	"github.com/tech-dipesh/drop/internal/ui"
	"github.com/tech-dipesh/drop/internal/wire"
)

const usage = `usage:
  drop                 wait for files, visible to nearby devices
  drop <files>         pick a nearby device and send files to it
  drop -name <name>    rename this device, then wait for files`

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	if err := run(ctx, os.Args[1:]); err != nil && ctx.Err() == nil {
		fmt.Fprintln(os.Stderr, "drop:", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) error {
	flags := flag.NewFlagSet("drop", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.Usage = func() { fmt.Println(usage) }
	rename := flags.String("name", "", "")
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	me, err := peer.Load()
	if err != nil {
		return err
	}
	if *rename != "" {
		me.Name = *rename
		if err := peer.Save(me); err != nil {
			return err
		}
	}
	if flags.NArg() == 0 {
		return listen(ctx, me)
	}
	return share(ctx, me, flags.Args())
}

func listen(ctx context.Context, me peer.Peer) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	dir := cmp.Or(os.Getenv("DROP_DIR"), filepath.Join(home, "Downloads"))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(wire.Port))
	if err != nil {
		return err
	}
	ui.Waiting(me.Name, dir)
	r := receive.Receiver{
		Dir:   dir,
		Ask:   func(o wire.Offer) bool { return ui.Accept(ctx, o) },
		Track: ui.Track,
		Done:  ui.Done,
	}
	errs := make(chan error, 2)
	go func() { errs <- discover.Answer(ctx, me) }()
	go func() { errs <- r.Serve(ctx, ln) }()
	return <-errs
}

func share(ctx context.Context, me peer.Peer, paths []string) error {
	offer, err := send.Build(me.Name, paths)
	if err != nil {
		return err
	}
	scan, stop := context.WithCancel(ctx)
	defer stop()
	found, err := discover.Scan(scan)
	if err != nil {
		return err
	}
	to, err := ui.Pick(ctx, found)
	if err != nil {
		return err
	}
	stop()
	s := send.Sender{
		Offered: func(o wire.Offer) { ui.Offered(o, to.Name) },
		Track:   ui.Track,
	}
	return s.Send(ctx, to.Addr, offer, paths)
}
