// The coil is the orchestrator and executor of the timeline's tasks.
// The timeline is the plan and the tasks are the points in time, which triggers the execution.
// Termination: The coild should stop only on statistical condition, never on a single failure.
package coil

import (
	"context"

	"github.com/csabakissmalta/tpee/timeline"
)

type Coil struct {
	Ctx       context.Context
	Timelines []timeline.Timeline
}

// It is a consumer of tasks relying on the main loop and context.
// It controls only the exact execution of the timeline
// Should start always from the first element and progressively consume the tasks.
func (c *Coil) Start(t timeline.Timeline) error {
	var e error
	for {
		if e != nil {
			return e
		}
	}
}

// Stops the coil loop
// conditions can be different:
// 1. out of tasks
// 2. exeception
// 3. condition based termination
func (c *Coil) Stop(t timeline.Timeline) error {
	return nil
}
