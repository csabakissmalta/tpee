// Task is the basic element of the test
// The functionality lies in the execution and reporting.
// The execution is bound to a time, when it supposed to happen.
// Single task failure should never block the timeline or execution of the test, perhaps a threshold can be established.
package task

import (
	"log"
	"net/http"
	"time"

	datastore "github.com/csabakissmalta/tpee/datastore"
	execconfig "github.com/csabakissmalta/tpee/exec"
	sessionstore "github.com/csabakissmalta/tpee/sessionstore"
)

type Task struct {
	// Planned relative execution time nanoseconds
	PlannedExecTimeNanos int

	// Execution absolute time
	ExecutionTime time.Time

	// HTTP request - ready to execute
	Request *http.Request

	// Executed
	Executed bool

	// Response
	Response *http.Response

	// Response receive time
	ResponseTime int64

	// Task label for reporting
	TaskLabel string
}

type Option func(*Task)

func WithRequest(req *http.Request) Option {
	return func(t *Task) {
		t.Request = req
	}
}

func WithPlannedExecTimeNanos(n int) Option {
	return func(t *Task) {
		t.PlannedExecTimeNanos = n
	}
}

func WithLabel(l string) Option {
	return func(t *Task) {
		t.TaskLabel = l
	}
}

func New(option ...Option) *Task {
	t := &Task{}
	for _, o := range option {
		o(t)
	}
	return t
}

func (ts *Task) Execute(c *http.Client, extract_rules []*execconfig.ExecRequestsElemDataPersistenceDataOutElem, r_ch chan *Task, extract_session bool, ss *sessionstore.Store) *Task {
	go func() {
		ts.ExecutionTime = time.Now()
		res, err := c.Do(ts.Request)
		if err != nil {
			log.Printf("ERROR: error executing request. %s", err.Error())
		}

		ts.ResponseTime = time.Since(ts.ExecutionTime).Milliseconds()
		ts.Response = res
		if res.StatusCode < 400 {
			if extract_session {
				go func() {
					e := ss.ExtractClientSessionFromResponse(res, ts.Request, nil) // <-- this needs to be corrected by the config, instead of nil
					if e != nil {
						log.Printf("ERROR: %s", e.Error())
					}
				}()
			}
			if len(extract_rules) > 0 {
				go datastore.ExtractDataFromResponse(res, extract_rules)
			}
		}
		ts.Executed = true
		if r_ch != nil {
			r_ch <- ts
		}
	}()
	return ts
}
