// The coil is the orchestrator and executor of the timeline's tasks.
// The timeline is the plan and the tasks are the points in time, which triggers the execution.
// Termination: The coild should stop only on statistical condition, never on a single failure.
package coil

import (
	"context"
	"time"

	"github.com/csabakissmalta/tpee/timeline"
)

type Coil struct {
	Ctx       context.Context
	Timelines []*timeline.Timeline
}

type Option func(*Coil)

func WithContext(ctx context.Context) Option {
	return func(c *Coil) {
		c.Ctx = ctx
	}
}

func WithTimelines(tls []*timeline.Timeline) Option {
	return func(c *Coil) {
		c.Timelines = tls
	}
}

func New(option ...Option) *Coil {
	c := &Coil{}
	for _, o := range option {
		o(c)
	}
	return c
}

// It is a consumer of tasks relying on the main loop and context.
// It controls only the exact execution of the timeline
// Should start always from the first element and progressively consume the tasks.
func (c *Coil) Start() {
	for _, tLine := range c.Timelines {
		consumeTimeline(tLine)
	}
	<-make(chan bool)
}

// Stops the coil loop
// conditions can be different:
// 1. out of tasks
// 2. exeception
// 3. condition based termination
func (c *Coil) Stop() error {
	return nil
}

// The Coil needs to control timelines in a separate routines
func consumeTimeline(tl *timeline.Timeline) {
	go func() {
		tl.CurrectTask = <-tl.Tasks
		if tl.CurrectTask.PlannedExecTimeNanos > 0 {
			time.Sleep(time.Duration(tl.CurrectTask.PlannedExecTimeNanos * int(time.Nanosecond)))
		}
		// compose/execute task here
		// --->
		tl.CurrectTask.Execute(nil)

		for {
			next := <-tl.Tasks
			dorm_period := (next.PlannedExecTimeNanos - tl.CurrectTask.PlannedExecTimeNanos) * int(time.Nanosecond)
			time.Sleep(time.Duration(dorm_period))
			// compose/execute task here
			// ---> here, in each step a correction needs to be added to the sleep time, due to the overhead of the composition

			tl.CurrectTask = next
		}
	}()
}
