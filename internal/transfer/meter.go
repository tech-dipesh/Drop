package transfer

import (
	"slices"
	"sync/atomic"
	"time"
)

const window = 2 * time.Second

type sample struct {
	at   time.Time
	done int64
}

type Stat struct {
	Done  int64
	Total int64
	Speed float64
	Left  time.Duration
}

type Meter struct {
	total   int64
	done    atomic.Int64
	samples []sample
}

func NewMeter(total int64) *Meter {
	return &Meter{total: total}
}

func (m *Meter) Add(n int64) {
	m.done.Add(n)
}

func (m *Meter) Tick(now time.Time) Stat {
	done := m.done.Load()
	m.samples = append(m.samples, sample{now, done})
	cut := slices.IndexFunc(m.samples, func(s sample) bool { return now.Sub(s.at) <= window })
	m.samples = m.samples[cut:]
	first := m.samples[0]
	stat := Stat{Done: done, Total: m.total}
	if elapsed := now.Sub(first.at).Seconds(); elapsed > 0 {
		stat.Speed = float64(done-first.done) / elapsed
	}
	if stat.Speed > 0 {
		stat.Left = time.Duration(float64(m.total-done) / stat.Speed * float64(time.Second))
	}
	return stat
}

func (s Stat) Percent() float64 {
	if s.Total == 0 {
		return 100
	}
	return float64(s.Done) / float64(s.Total) * 100
}
