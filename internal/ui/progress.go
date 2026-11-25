package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/tech-dipesh/drop/internal/transfer"
	"github.com/tech-dipesh/drop/internal/wire"
)

const width = 78

func Track(o wire.Offer, m *transfer.Meter) func() {
	name := label(o)
	start := time.Now()
	quit := make(chan struct{})
	finished := make(chan struct{})
	go func() {
		defer close(finished)
		tick := time.NewTicker(100 * time.Millisecond)
		defer tick.Stop()
		for {
			select {
			case now := <-tick.C:
				draw(line(name, m.Tick(now)))
			case <-quit:
				draw(summary(name, m.Tick(time.Now()).Done, time.Since(start)))
				fmt.Println()
				return
			}
		}
	}()
	return func() {
		close(quit)
		<-finished
	}
}

func line(name string, s transfer.Stat) string {
	pct := s.Percent()
	full := int(pct / 5)
	bar := strings.Repeat("=", full) + strings.Repeat(" ", 20-full)
	left := "--"
	if s.Left > 0 {
		left = s.Left.Round(time.Second).String()
	}
	return fmt.Sprintf("%-20.20s [%s] %5.1f%% %10s/s %8s left", name, bar, pct, size(int64(s.Speed)), left)
}

func summary(name string, done int64, elapsed time.Duration) string {
	avg := float64(done) / elapsed.Seconds()
	return fmt.Sprintf("%-20.20s %s in %.2fs (%s/s)", name, size(done), elapsed.Seconds(), size(int64(avg)))
}

func draw(text string) {
	fmt.Printf("\r%-*s", width, text)
}
