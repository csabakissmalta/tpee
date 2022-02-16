// The timeline is the passive, blueprint of the test
// It has tasks and duration and other metadata, which will affect the execution
package timeline

import (
	"github.com/csabakissmalta/tpee/postman"
	"github.com/csabakissmalta/tpee/task"
)

// Execution parameters for the timeline
// Loaded from a config, ideally
type ExecRules struct {
	// Duration
	DurationSec int64

	// frequency aka RPS
	Frequency int64

	// delay
	Delay int64
}

type Timeline struct {
	// Name - identical/matching with the tasks name
	// In any timeline there is only one kind of task
	Name string

	// A channel of tasks, which will be executed at some point in time.
	Tasks chan *task.Task

	// Execution details
	Rules ExecRules

	// Current Task
	CurrectTask *task.Task
}

func (t *Timeline) Populate(r *postman.Request) {
	// Create time markers - empty tasks
	calc_periods(&t.Rules, r)
}
