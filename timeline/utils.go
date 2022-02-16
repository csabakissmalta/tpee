package timeline

import (
	execconf "github.com/csabakissmalta/tpee/exec"

	"github.com/csabakissmalta/tpee/postman"
	"github.com/csabakissmalta/tpee/task"
)

func calc_periods(dur int, er *execconf.ExecRequestsElem, rq *postman.Request) chan *task.Task {
	// the count of markers is (duration - delay) * frequency
	m_count := (dur - er.DelaySeconds) * er.Frequency
	ch := make(chan *task.Task, m_count)
	for i := 0; i < int(m_count); i++ {
		ch <- task.New(
			task.WithRequest(rq),
		)
	}
	return ch
}
