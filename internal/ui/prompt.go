package ui

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/tech-dipesh/drop/internal/wire"
)

var stdin = bufio.NewReader(os.Stdin)

func Waiting(dir string) {
	fmt.Printf("waiting for files, saving to %s (ctrl+c to stop)\n", dir)
}

func Accept(ctx context.Context, o wire.Offer) bool {
	fmt.Printf("\n%s wants to send %s (%s). accept? [y/n] ", o.From, label(o), size(o.Total()))
	answer := make(chan string, 1)
	go func() {
		line, _ := stdin.ReadString('\n')
		answer <- strings.ToLower(strings.TrimSpace(line))
	}()
	select {
	case a := <-answer:
		return a == "y" || a == "yes"
	case <-ctx.Done():
		return false
	}
}

func Offered(o wire.Offer) {
	fmt.Printf("offered %s (%s), waiting for the other side to accept\n", label(o), size(o.Total()))
}

func Done(saved []string, err error) {
	for _, path := range saved {
		fmt.Println("saved", path)
	}
	if err != nil {
		fmt.Println("transfer failed:", err)
	}
}

func label(o wire.Offer) string {
	if len(o.Files) == 1 {
		return o.Files[0].Name
	}
	return fmt.Sprintf("%d files", len(o.Files))
}

func size(n int64) string {
	units := []string{"B", "KB", "MB", "GB", "TB"}
	v := float64(n)
	i := 0
	for v >= 1024 && i < len(units)-1 {
		v /= 1024
		i++
	}
	return fmt.Sprintf("%.1f %s", v, units[i])
}
