// The timeline is the passive, blueprint of the test
// It has tasks and duration and other metadata, which will affect the execution
package timeline

import (
	"log"

	execconf "github.com/csabakissmalta/tpee/exec"
	"github.com/csabakissmalta/tpee/postman"
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

	// feeds
	Feeds []*Feed
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

func (t *Timeline) Populate(dur int, r *postman.Request, env []*execconf.ExecEnvironmentElem) {
	// do checks on the Postman Request instance and log status
	e := check_postman_request_and_validate_requirements(r, env)
	if e != nil {
		log.Fatalf("DATA ERROR: %s", e.Error())
	}

	// check env elements and load feeds if there is any feedValue type
	timeline_dimension := (dur - t.Rules.DelaySeconds) * t.Rules.Frequency
	t.Feeds = load_feeds_if_required(timeline_dimension, env)

	// Create time markers - empty tasks
	t.Tasks = calc_periods(dur, t.Rules, r)
}
