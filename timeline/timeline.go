// The timeline is the passive, blueprint of the test
// It has tasks and duration and other metadata, which will affect the execution
package timeline

import (
	"net/http"
	"time"

	execconf "github.com/csabakissmalta/tpee/exec"
	postman "github.com/csabakissmalta/tpee/postman"
	task "github.com/csabakissmalta/tpee/task"
	"github.com/nats-io/nats.go"
)

type Timeline struct {
	// Name - identical/matching with the tasks name
	// In any timeline there is only one kind of task
	Name string

	// A channel of tasks, which will be executed at some point in time.
	Tasks chan *task.Task

	// rampup calls count
	RamUpCallsCount int

	// rampup tasks
	RampupTasks chan *task.Task

	// Request blueprint
	RequestBlueprint *postman.Request

	// Execution details
	Rules *execconf.ExecRequestsElem

	// Current Task
	CurrectTask *task.Task

	// Feeds
	Feeds []*Feed

	// HTTP client
	HTTPClient *http.Client

	// Step duration
	StepDuration int
}

// rate change delta
type Transition struct {
	// target rate
	TargetRate int

	// rampup to achieve the changed trate
	// this is identical with the rampup seconds
	TransitionRampupTimeSeconds int

	// rampup type - so far sinisoidal implemented only
	RampupType string
}

type Option func(*Timeline)

func WithRules(rls *execconf.ExecRequestsElem) Option {
	return func(tl *Timeline) {
		tl.Rules = rls
	}
}

func WithHTTPClient(c *http.Client) Option {
	return func(tl *Timeline) {
		tl.HTTPClient = c
	}
}

func WithName(n string) Option {
	return func(t *Timeline) {
		t.Name = n
	}
}

func WithRampup(rmp chan *task.Task) Option {
	return func(t *Timeline) {
		t.RampupTasks = rmp
	}
}

func New(option ...Option) *Timeline {
	tl := &Timeline{}
	for _, o := range option {
		o(tl)
	}
	return tl
}

func (t *Timeline) Populate(dur int, r *postman.Request, env []*execconf.ExecEnvironmentElem, rmp *execconf.ExecRampup, subs map[string]chan *nats.Msg) {
	// populate rampup period if set
	if rmp != nil {
		if rmp != nil {
			rmp_points := t.GenerateRampUpTimeline(int64(*rmp.DurationSeconds), 0, int64(t.Rules.Frequency), float64(t.Rules.DelaySeconds), Rampup(*rmp.RampupType), t.Rules.Name)
			t.RampupTasks = make(chan *task.Task, len(rmp_points))
			for _, p := range rmp_points {
				t.RampupTasks <- p
			}
		}
	}

	// check env elements and load feeds if there is any feedValue type
	timeline_dimension := (dur-t.Rules.DelaySeconds)*t.Rules.Frequency + len(t.RampupTasks) + 1000
	t.Feeds = load_feeds_if_required(timeline_dimension, env, subs)

	// The step between the markers
	second := time.Second
	convers := int(second / time.Nanosecond)
	step := int(convers / t.Rules.Frequency)

	// Set the step duration for the timeline as well
	t.StepDuration = step

	// Create time markers - empty tasks
	t.Tasks = calc_periods(dur, step, t.Rules, r)

	// set the resulting postman request
	t.RequestBlueprint = r
}

// repopulate
func (t *Timeline) Repopulate(tr *Transition, test_duration int, start_time time.Time, ticker *time.Ticker) {
	rmp_points := t.GenerateRampUpTimeline(int64(tr.TransitionRampupTimeSeconds), int64(t.Rules.Frequency), int64(tr.TargetRate), float64(time.Since(start_time).Seconds()), Rampup(tr.RampupType), t.Rules.Name)
	t.RampupTasks = make(chan *task.Task, len(rmp_points))
	for _, p := range rmp_points {
		t.RampupTasks <- p
	}

	// check env elements and load feeds if there is any feedValue type
	// elapsed := time.Since(start_time)
	// dur := test_duration - int(elapsed.Seconds())
	// timeline_dimension := (dur-t.Rules.DelaySeconds)*t.Rules.Frequency + len(t.RampupTasks) + 1000

	// The step between the markers
	second := time.Second
	convers := int(second / time.Nanosecond)
	step := int(convers / tr.TargetRate)

	// Set the step duration for the timeline as well
	t.StepDuration = step

	// Create time markers - empty tasks
	t.Rules.Frequency = tr.TargetRate

	required_task_count := tr.TargetRate * (test_duration - int(time.Since(start_time).Seconds()))
	if len(t.Tasks) < required_task_count {
		generateAdditionalTasks(required_task_count, step, t.Tasks, t.Rules, t.RequestBlueprint)
	}

	// t.Tasks = new_tasks
	ticker.Reset(time.Duration(step))

	// L:
	//
	//	for {
	//		select {
	//		case <-ch_to_empty:
	//		default:
	//			// close(ch_to_empty)
	//			break L
	//		}
	//	}
}

func CheckPostmanRequestAndValidateRequirements(pr *postman.Request, env []*execconf.ExecEnvironmentElem) error {
	// check URL raw
	res, e := validate_and_substitute(pr.URL.Raw, r, env)
	if e != nil {
		return e
	}
	pr.URL.Raw = res
	for i, h := range pr.URL.Host {
		res, e = validate_and_substitute(h, r, env)
		if e != nil {
			return e
		}
		pr.URL.Host[i] = res
	}
	for j, p := range pr.URL.Path {
		res, e = validate_and_substitute(p, r, env)
		if e != nil {
			return e
		}
		pr.URL.Path[j] = res
	}

	// check Headers
	for k, hdr := range pr.Header {
		res, e = validate_and_substitute(hdr.Value, r, env)
		if e != nil {
			return e
		}
		pr.Header[k].Value = res
	}
	// check body
	if len(pr.Body.Raw) > 0 {
		res, e = validate_and_substitute(pr.Body.Raw, r, env)
		if e != nil {
			return e
		}
		pr.Body.Raw = res
	}
	if len(pr.Body.Formdata) > 0 {
		for l, fd := range pr.Body.Formdata {
			res, e = validate_and_substitute(fd.Value, r, env)
			if e != nil {
				return e
			}
			pr.Body.Formdata[l].Value = res
		}
	}
	if len(pr.Body.Urlencoded) > 0 {
		for l, ue := range pr.Body.Urlencoded {
			res, e = validate_and_substitute(ue.Value, r, env)
			if e != nil {
				return e
			}
			pr.Body.Urlencoded[l].Value = res
		}
	}
	return nil
}
