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
	L = 50.0 // do not change this, unless multiple rampup types, but then create a corresponding const/value to each type
)

var (
	multiplier float64 = 4.0  // multiplies L and gets as much seconds duration for the rampup
	A          float64 = 20.0 // amplitude - the target rps
)

func AdjustMultiplier(total float64) {
	multiplier = total / L
}

func GetRampupDuration() (dur float64) {
	dur = L * multiplier
	return dur
}

func calc_val(x float64, tp Rampup) (pt_val float64) {
	// sinusoidal increase
	switch tp {
	case SINUSOIDAL:
		pt_val = A * (math.Cos(math.Pi*(x-1)) + 1) / 2
	case LINEAR:
		pt_val = A * x
	}

	return pt_val
}

func generate_intervals(t Rampup, dur int) (result []float64, count int) {
	stepper := float64(0.1 / multiplier)
	rpss := []float64{}
	for x := 0.0; x < float64(dur)*multiplier; x += 1.0 {
		curr := calc_val(x*stepper, t)
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
	A = float64(maxRps)
	pts, count = generate_intervals(t, dur)
	return pts, count
}
