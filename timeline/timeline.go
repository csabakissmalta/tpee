// The timeline is the passive, blueprint of the test
// It has tasks and duration and other metadata, which will affect the execution
package timeline

import "github.com/csabakissmalta/tpee/task"

// Execution parameters for the timeline
// Loaded from a config, ideally
type Exec struct {
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
	Exec Exec
}

func (t *Timeline) Populate(e *Exec) {

}
