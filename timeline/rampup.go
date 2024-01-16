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

func calc_val(x float64, tp Rampup, dur int, maxrps int, initrps int) (pt_val float64) {
	// sinusoidal increase
	switch tp {
	case SINUSOIDAL:
		rate_diff := float64(maxrps - initrps)
		pt_val = float64(rate_diff)/2*(math.Cos((math.Pi*x)/float64(dur)-math.Pi)+1) + float64(initrps)
	case LINEAR:
		pt_val = A * x
	}
	return pt_val
}

func generate_intervals(t Rampup, dur int, initrps int, maxrps int) (result []float64, count int) {
	rpss := []float64{}
	for x := -1.0; x <= float64(dur+1); x += 1.0 {
		curr := calc_val(x, t, dur, maxrps, initrps)
		rpss = append(rpss, curr)
	}
	init_rpss_step := 0
	if initrps > maxrps {
		sort.Sort(sort.Reverse(sort.Float64Slice(rpss)))
		init_rpss_step = 0
	} else {
		sort.Float64s(rpss)
	}

	for i := init_rpss_step; i < len(rpss); i++ {
		for f := 0; f < int(rpss[i]); f++ {
			microstep := 1.0 / rpss[i]
			result = append(result, float64(i)+float64(f)*microstep)
		}
	}
	count = len(result)
	sort.Float64s(result)
	return result, count
}

// generates values based on rethought formula: y\ =\ A\frac{1\ +\ \cos\left(\pi\left(\frac{x}{D}-1\right)\right)}{2} (paste to Desmos)
// func generate_intervals_2(t Rampup, dur int, maxrps int) (result []float64, count int) {
// 	for i := 1; i < maxrps; i++ {
// 		// val := 2*float64(dur) - (float64(dur)*(math.Acos(float64((2*float64(i))/float64(maxrps)-1)))+math.Pi)/math.Pi

// 		nom := float64(dur) * (math.Acos(2*(float64(i))/float64(maxrps)-1) + math.Pi)
// 		val := 2*float64(dur) - (nom / math.Pi)

// 		result = append(result, val)
// 	}

// 	return result, maxrps
// }

// This method should provide a timeline
// with negative timestamps to be executed before the main timeline
// should take:
// time length and the rate of what the plan supposed to reach and the generator function.
func (tl *Timeline) GenerateRampUpTimeline(l int64, initrps int64, targetRPS int64, delay float64, t Rampup, label string) (rampupPts []*task.Task) {
	initPoints, c := PointsPlannedTimestamps(initrps, targetRPS, t, int(l))
	second := float64(time.Second)
	tl.RamUpCallsCount = c

	for _, p := range initPoints {
		t := int((p+delay)*second) / int(time.Nanosecond)
		// log.Println("------")
		// log.Println("POINT: ", p)
		// log.Println("PLANNED RAMPUP NANO: ", t)
		rampupPts = append(rampupPts, task.New(
			task.WithPlannedExecTimeNanos(int(t)),
			task.WithLabel(label),
		))

	}
	return rampupPts
}

// The function, returning the values
func PointsPlannedTimestamps(initRps int64, maxRps int64, t Rampup, dur int) (pts []float64, count int) {
	pts, count = generate_intervals(t, dur, int(initRps), int(maxRps))
	// pts, count = generate_intervals_2(t, dur, int(maxRps))
	return pts, count
}
