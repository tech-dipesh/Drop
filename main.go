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

	"github.com/tech-dipesh/drop/internal/config"
	"github.com/tech-dipesh/drop/internal/discover"
	"github.com/tech-dipesh/drop/internal/peer"
	"github.com/tech-dipesh/drop/internal/receive"
	"github.com/tech-dipesh/drop/internal/send"
	"github.com/tech-dipesh/drop/internal/ui"
	"github.com/tech-dipesh/drop/internal/wire"
)

var version = "dev"

const usage = `usage:
  drop                  wait for files, visible to nearby devices
  drop <files>          pick a nearby device and send files to it
  drop -speed           measure the link to a nearby device, nothing is saved
  drop -name <name>     set the name other devices see
  drop -dir <folder>    set where received files are saved
  drop -version         print the version`

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
	name := flags.String("name", "", "")
	dir := flags.String("dir", "", "")
	speed := flags.Bool("speed", false, "")
	showVersion := flags.Bool("version", false, "")
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	if *showVersion {
		fmt.Println("drop", version)
		return nil
	}
	conf, err := config.Load()
	if err != nil {
		return err
	}
	if *name != "" || *dir != "" {
		folder, err := filepath.Abs(cmp.Or(*dir, conf.Dir))
		if err != nil {
			return err
		}
		conf.Name = cmp.Or(*name, conf.Name)
		conf.Dir = folder
		if err := config.Save(conf); err != nil {
			return err
		}
	}
	if *speed {
		return share(ctx, send.SpeedTest(conf.Name), nil)
	}
	if flags.NArg() == 0 {
		return listen(ctx, peer.Peer{ID: conf.ID, Name: conf.Name}, cmp.Or(os.Getenv("DROP_DIR"), conf.Dir))
	}
	offer, err := send.Build(conf.Name, flags.Args())
	if err != nil {
		return err
	}
	return share(ctx, offer, flags.Args())
}

func listen(ctx context.Context, me peer.Peer, dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(wire.Port))
	if err != nil {
		return err
	}
	ui.Waiting(me.Name, dir)
	r := receive.Receiver{
		Dir:    dir,
		Paired: ui.Paired,
		Ask:    func(o wire.Offer) bool { return ui.Accept(ctx, o) },
		Track:  ui.Track,
		Done:   ui.Done,
	}
	errs := make(chan error, 2)
	go func() { errs <- discover.Answer(ctx, me) }()
	go func() { errs <- r.Serve(ctx, ln) }()
	return <-errs
}

func share(ctx context.Context, offer wire.Offer, paths []string) error {
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
		Paired:  ui.Paired,
		Offered: func(o wire.Offer) { ui.Offered(o, to.Name) },
		Track:   ui.Track,
	}
	return s.Send(ctx, to.Addr, offer, paths)
}
