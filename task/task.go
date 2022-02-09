// Task is the basic element of the test
// The functionality lies in the execution and reporting.
// The execution is bound to a time, when it supposed to happen.
// Single task should never block on itself.
package task

type ITask interface {
	Execute() error
	Report(interface{}) interface{}
}

type Task struct {
}

func (ts *Task) Execute() error {
	return nil
}

func (ts *Task) Report(taskdata interface{}) interface{} {
	return nil
}
