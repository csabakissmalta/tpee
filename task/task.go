// Task is the basic element of the test
// The functionality lies in the execution and reporting.
// The execution is bound to a time, when it supposed to happen.
// Single task failure should never block the timeline or execution of the test, perhaps a threshold can be established.
package task

import (
	"log"
	"net/http"
	"time"
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

func New(option ...Option) *Task {
	t := &Task{}
	for _, o := range option {
		o(t)
	}
	return t
}

func (ts *Task) Execute(c *http.Client) error {
	go func() {
		res, err := c.Do(ts.Request)
		if err != nil {
			log.Printf("ERROR: error executing request. %s", err.Error())
		}
		log.Println("STATUS: ", res.StatusCode)
	}()
	return nil
}

func (ts *Task) Report(taskdata interface{}) interface{} {
	return nil
}
