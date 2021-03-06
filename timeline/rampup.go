package timeline

import (
	"math"
	"sort"
	"time"

	task "github.com/csabakissmalta/tpee/task"
)

type Rampup string

const (
	LINEAR     Rampup = "linear"
	SINUSOIDAL Rampup = "sinusoidal"
)

const (
	L = 100.0 // do not change this, unless multiple rampup types, but then create a corresponding const/value to each type
)

var (
	A float64 = 20.0 // amplitude - the target rps
)

// y\ =\ \frac{A}{2}\left(\cos\left(\frac{\pi x}{b}\ -\ \pi\right)+1\right)

func calc_val(x float64, tp Rampup, dur int, maxrps int) (pt_val float64) {
	// sinusoidal increase
	switch tp {
	case SINUSOIDAL:
		pt_val = float64(maxrps) / 2 * (math.Cos((math.Pi*x)/float64(dur)-math.Pi) + 1)
		// pt_val = A * (math.Cos(math.Pi*(x-1)) + 1) / 2
	case LINEAR:
		pt_val = A * x
	}
	return pt_val
}

func generate_intervals(t Rampup, dur int, maxrps int) (result []float64, count int) {
	rpss := []float64{}
	for x := 0.0; x < float64(dur); x += 1.0 {
		curr := calc_val(x, t, dur, maxrps)
		rpss = append(rpss, curr)
	}
	sort.Float64s(rpss)

	for i := 1; i < len(rpss); i++ {
		for f := 0; f <= int(rpss[i]); f++ {
			microstep := 1.0 / rpss[i]
			result = append(result, float64(i)+float64(f)*microstep)
		}
	}
	count = len(result)
	sort.Float64s(result)
	return result, count
}

// This method should provide a timeline
// with negative timestamps to be executed before the main timeline
// should take:
// time length and the rate of what the plan supposed to reach and the generator function.
func (tl *Timeline) GenerateRampUpTimeline(l int64, targetRPS int64, delay float64, t Rampup, label string) (rampupPts []*task.Task) {
	initPoints, c := PointsPlannedTimestamps(targetRPS, t, int(l))
	second := float64(time.Second)
	tl.RamUpCallsCount = c

	for _, p := range initPoints {
		t := ((p + delay) * second) / float64(time.Nanosecond)
		rampupPts = append(rampupPts, task.New(
			task.WithPlannedExecTimeNanos(int(t)),
			task.WithLabel(label),
		))

	}
	return rampupPts
}

// The function, returning the values
func PointsPlannedTimestamps(maxRps int64, t Rampup, dur int) (pts []float64, count int) {
	pts, count = generate_intervals(t, dur, int(maxRps))
	return pts, count
}
