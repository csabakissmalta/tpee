// Task is the basic element of the test
// The functionality lies in the execution and reporting.
// The execution is bound to a time, when it supposed to happen.
// Single task should never block on itself.
package task

import (
	"net/http"

	"github.com/csabakissmalta/tpee/postman"
)

type Task struct {
	// Planned relative execution time nanoseconds
	PlannedExecTimeNanos int64

	// A request to be executed
	PostmanRequest *postman.Request

	// HTTP request - ready to execute
	Request *http.Request

	// Executed
	Executed bool
}

type Option func(*Task)

func WithRequest(req *postman.Request) Option {
	return func(t *Task) {
		t.PostmanRequest = req
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
	return nil
}

func (ts *Task) Report(taskdata interface{}) interface{} {
	return nil
}
