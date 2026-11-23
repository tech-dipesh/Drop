package transfer

import (
	"testing"
	"time"
)

const mb = 1 << 20

func TestMeterSpeedOverWindow(t *testing.T) {
	m := NewMeter(100 * mb)
	start := time.Now()
	stat := m.Tick(start)
	if stat.Speed != 0 || stat.Left != 0 {
		t.Fatalf("first tick got %+v, want no speed yet", stat)
	}
	for sec := 1; sec <= 3; sec++ {
		m.Add(10 * mb)
		stat = m.Tick(start.Add(time.Duration(sec) * time.Second))
	}
	if stat.Speed != 10*mb {
		t.Fatalf("speed %.0f, want %d", stat.Speed, 10*mb)
	}
	if stat.Left != 7*time.Second {
		t.Fatalf("left %v, want 7s", stat.Left)
	}
}

func TestPercent(t *testing.T) {
	if p := (Stat{Done: 25, Total: 100}).Percent(); p != 25 {
		t.Fatalf("got %v, want 25", p)
	}
	if p := (Stat{}).Percent(); p != 100 {
		t.Fatalf("empty transfer got %v, want 100", p)
	}
}
