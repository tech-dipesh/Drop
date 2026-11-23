package main

import (
	"cmp"
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"

	"github.com/tech-dipesh/drop/internal/receive"
	"github.com/tech-dipesh/drop/internal/send"
	"github.com/tech-dipesh/drop/internal/ui"
	"github.com/tech-dipesh/drop/internal/wire"
)

const usage = `usage:
  drop            wait for files from a nearby device
  drop <files>    send files to a nearby device`

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	if err := run(ctx, os.Args[1:]); err != nil && ctx.Err() == nil {
		fmt.Fprintln(os.Stderr, "drop:", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) error {
	if len(args) > 0 && (args[0] == "-h" || args[0] == "--help") {
		fmt.Println(usage)
		return nil
	}
	port := strconv.Itoa(wire.Port)
	if len(args) == 0 {
		return listen(ctx, port)
	}
	host, _ := os.Hostname()
	s := send.Sender{From: cmp.Or(host, "unknown device"), Offered: ui.Offered, Track: ui.Track}
	return s.Send(ctx, net.JoinHostPort("localhost", port), args)
}

func listen(ctx context.Context, port string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	dir := cmp.Or(os.Getenv("DROP_DIR"), filepath.Join(home, "Downloads"))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	ui.Waiting(dir)
	r := receive.Receiver{
		Dir:   dir,
		Ask:   func(o wire.Offer) bool { return ui.Accept(ctx, o) },
		Track: ui.Track,
		Done:  ui.Done,
	}
	return r.Serve(ctx, ln)
}
