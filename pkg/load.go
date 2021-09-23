package pkg

import (
	"fmt"
	"time"
)

const (
	FSHIFT  = 11
	FIXED_1 = (1 << FSHIFT)
	EXP_1   = 1884 /* 1/exp(5sec/1min) as fixed-point */
	EXP_5   = 2014 /* 1/exp(5sec/5min) */
	EXP_15  = 2037 /* 1/exp(5sec/15min) */
)

type LoadMonitor struct {
	AvenRun1  uint64
	AvenRun5  uint64
	AvenRun15 uint64
}

func (l *LoadMonitor) Run() {
	ticker := time.NewTicker(100 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			l.refreshLoad()
		}
	}
}

// load1 = load0 * exp + active * (1 - exp)
func calcLoad(load0 uint64, exp uint64, active uint64) uint64 {
	if active > 0 {
		active = active * FIXED_1
	}
	load1 := load0*exp + active*(FIXED_1-exp)
	if active >= load0 {
		load1 += FIXED_1 - 1
	}

	return load1 / FIXED_1
}

func (l *LoadMonitor) refreshLoad() {
	tc := GetTaskCount()
	runPid := tc.Running + tc.Uninterruptible
	l.AvenRun1 = calcLoad(l.AvenRun1, EXP_1, uint64(runPid))
	l.AvenRun5 = calcLoad(l.AvenRun5, EXP_5, uint64(runPid))
	l.AvenRun15 = calcLoad(l.AvenRun15, EXP_15, uint64(runPid))
}

func loadInt(x uint64) uint64 {
	return x >> FSHIFT
}

func loadFrac(x uint64) uint64 {
	return loadInt(((x) & (FIXED_1 - 1)) * 100)
}

func (l *LoadMonitor) GetLoad() string {
	a := l.AvenRun1 + (FIXED_1 / 200)
	b := l.AvenRun5 + (FIXED_1 / 200)
	c := l.AvenRun15 + (FIXED_1 / 200)
	return fmt.Sprintf("%d.%02d, %d.%02d, %d.%02d", loadInt(a), loadFrac(a), loadInt(b), loadFrac(b), loadInt(c), loadFrac(c))
}

func NewLoadMonitor() *LoadMonitor {
	l := LoadMonitor{}
	return &l
}
