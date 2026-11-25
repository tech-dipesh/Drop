package ui

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/tech-dipesh/drop/internal/wire"
)

var stdin = bufio.NewReader(os.Stdin)

func Waiting(name, dir string) {
	fmt.Printf("visible as %s, saving to %s (ctrl+c to stop)\n", clean(name), dir)
}

func Paired(code string) {
	fmt.Printf("\npairing code %s, check it matches the other screen\n", code)
}

func Accept(ctx context.Context, o wire.Offer) bool {
	what := fmt.Sprintf("send %s (%s)", label(o), size(o.Total()))
	if o.Test {
		what = fmt.Sprintf("run a speed test (%s, nothing is saved)", size(o.Total()))
	}
	fmt.Printf("%s wants to %s. accept? [y/n] ", clean(o.From), what)
	select {
	case a := <-readLine():
		return a == "y" || a == "yes"
	case <-ctx.Done():
		return false
	}
}

func Offered(o wire.Offer, to string) {
	fmt.Printf("offered %s (%s) to %s, waiting for them to accept\n", label(o), size(o.Total()), clean(to))
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
		return clean(o.Files[0].Name)
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
