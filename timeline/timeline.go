// The timeline is the passive, blueprint of the test
// It has tasks and duration and other metadata, which will affect the execution
package timeline

import (
	"net/http"

	execconf "github.com/csabakissmalta/tpee/exec"
	postman "github.com/csabakissmalta/tpee/postman"
	task "github.com/csabakissmalta/tpee/task"
)

type Timeline struct {
	// Name - identical/matching with the tasks name
	// In any timeline there is only one kind of task
	Name string

	// A channel of tasks, which will be executed at some point in time.
	Tasks chan *task.Task

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

func New(option ...Option) *Timeline {
	tl := &Timeline{}
	for _, o := range option {
		o(tl)
	}
	return tl
}

func (t *Timeline) Populate(dur int, r *postman.Request, env []*execconf.ExecEnvironmentElem) {
	// check env elements and load feeds if there is any feedValue type
	timeline_dimension := (dur - t.Rules.DelaySeconds) * t.Rules.Frequency
	t.Feeds = load_feeds_if_required(timeline_dimension, env)

	// Create time markers - empty tasks
	t.Tasks = calc_periods(dur, t.Rules, r)

	// set the resulting postman request
	t.RequestBlueprint = r
}

func CheckPostmanRequestAndValidateRequirements(pr *postman.Request, env []*execconf.ExecEnvironmentElem) error {
	// check URL raw
	e := validate_and_substitute(&pr.URL.Raw, r, env)
	if e != nil {
		return e
	}
	for i, h := range pr.URL.Host {
		e = validate_and_substitute(&h, r, env)
		if e != nil {
			return e
		}
		pr.URL.Host[i] = h
	}
	for j, p := range pr.URL.Path {
		e = validate_and_substitute(&p, r, env)
		if e != nil {
			return e
		}
		pr.URL.Path[j] = p
	}

	// check Headers
	for k, hdr := range pr.Header {
		e = validate_and_substitute(&hdr.Value, r, env)
		if e != nil {
			return e
		}
		pr.Header[k].Value = hdr.Value
	}
	// check body
	if len(pr.Body.Raw) > 0 {
		e = validate_and_substitute(&pr.Body.Raw, r, env)
		if e != nil {
			return e
		}
	}
	if len(pr.Body.Formdata) > 0 {
		for l, fd := range pr.Body.Formdata {
			e = validate_and_substitute(&fd.Value, r, env)
			if e != nil {
				return e
			}
			pr.Body.Formdata[l].Value = fd.Value
		}
	}
	if len(pr.Body.Urlencoded) > 0 {
		for l, ue := range pr.Body.Urlencoded {
			e = validate_and_substitute(&ue.Value, r, env)
			if e != nil {
				return e
			}
			pr.Body.Urlencoded[l].Value = ue.Value
		}
	}
	return nil
}
