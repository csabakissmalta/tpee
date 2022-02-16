// The timeline is the passive, blueprint of the test
// It has tasks and duration and other metadata, which will affect the execution
package timeline

import (
	execconf "github.com/csabakissmalta/tpee/exec"
	"github.com/csabakissmalta/tpee/postman"
	"github.com/csabakissmalta/tpee/request"
	"github.com/csabakissmalta/tpee/task"
)

type Timeline struct {
	// Name - identical/matching with the tasks name
	// In any timeline there is only one kind of task
	Name string

	// A channel of tasks, which will be executed at some point in time.
	Tasks chan *task.Task

	// Execution details
	Rules *execconf.ExecRequestsElem

	// Current Task
	CurrectTask *task.Task
}

type Option func(*Timeline)

func WithRules(rls *execconf.ExecRequestsElem) Option {
	return func(tl *Timeline) {
		tl.Rules = rls
	}
}

func New(option ...Option) *Timeline {
	tl := &Timeline{}
	for _, o := range option {
		o(tl)
	}
	return tl
}

func (t *Timeline) Populate(dur int, r *postman.Request) {
	// Create time markers - empty tasks
	t.Tasks = calc_periods(dur, t.Rules, r)

	// pre-process the tasks and check dependencies
	for tsk := range t.Tasks {
		request.ComposeHttpRequest(tsk.Request)
		t.Tasks <- tsk
	}
}
