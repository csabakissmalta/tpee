// Task is the basic element of the test
// The functionality lies in the execution and reporting.
// The execution is bound to a time, when it supposed to happen.
// Single task failure should never block the timeline or execution of the test, perhaps a threshold can be established.
package task

import (
	"log"
	"net/http"
	"time"

	data "github.com/csabakissmalta/tpee/data"
	"github.com/csabakissmalta/tpee/datastore"
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

	// rampup flag
	IsRampup bool
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

// util method does the request require session?
func isSessionRequired(dp []*execconfig.ExecRequestsElemDataPersistenceDataInElem) bool {
	for _, d := range dp {
		if d.Storage.(string) == "session-meta" {
			return true
		}
	}
	return false
}

func (ts *Task) Execute(c *http.Client, extract_rules []*execconfig.ExecRequestsElemDataPersistenceDataOutElem, data_in_rules []*execconfig.ExecRequestsElemDataPersistenceDataInElem, r_ch chan *Task, extract_session bool, ss *sessionstore.Store, ds *datastore.DataBroadcaster) *Task {

	go func() {
		ts.ExecutionTime = time.Now()
		res, err := c.Do(ts.Request)
		if err != nil {
			log.Printf("ERROR: error executing request. %s", err.Error())
		}
		var session *sessionstore.Session
		if len(extract_rules) > 0 {
			session_required := isSessionRequired(data_in_rules)
			if session_required {
				for {
					session = <-ss.SessionOut
					if time.Since(session.Created) < sessionstore.SESSION_VALIDITY {
						break
					}
				}
			}
		}

		ts.ResponseTime = time.Since(ts.ExecutionTime).Milliseconds()
		ts.Response = res
		if res.StatusCode < 400 {

			// go func() {

			if extract_session {
				// meta := &sessionstore.Meta{}
				session, err = ss.ExtractClientSessionFromResponse(res, ts.Request, nil) // <-- this needs to be corrected by the config, instead of nil
				if err != nil {
					log.Printf("ERROR: %s", err.Error())
				}

			}

			for _, erule := range extract_rules {
				log.Println(erule.Storage.(string))
				switch erule.Storage.(string) {
				case "data-store":
					data.ExtractDataFromResponse(res, erule, ds)
				case "session-meta":
					data.ExtractDataFromResponse(res, erule, session)
				default:
					log.Printf("tpee: %s", "default")
					// nothing happens
				}
			}

			ss.SessionIn <- session
			// }()

		}
		ts.Executed = true
		if r_ch != nil {
			r_ch <- ts
		}
	}()
	return ts
}
