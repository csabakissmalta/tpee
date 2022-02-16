package timeline

import (
	"fmt"
	"regexp"
	"strings"

	execconf "github.com/csabakissmalta/tpee/exec"

	"github.com/csabakissmalta/tpee/postman"
	"github.com/csabakissmalta/tpee/task"
)

var r = regexp.MustCompile(`[\{]{2}(.{1,})[\}]{2}`)

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

func check_env_var_set(vname string, env []*execconf.ExecEnvironmentElem) (bool, string) {
	for _, envElem := range env {
		if envElem.Key == vname && len(envElem.Value) > 0 {
			return true, envElem.Value
		}
	}
	return false, ""
}

func check_postman_request_and_validate_requirements(pr *postman.Request, env []*execconf.ExecEnvironmentElem) error {
	// --- ENVIRONMENT ---
	// check URL
	match := r.FindStringSubmatch(pr.URL.Raw)
	if len(match) > 1 {
		exists, val := check_env_var_set(match[1], env)
		if !exists {
			return fmt.Errorf("%s env variable for URL: %s is not set", match[1], pr.URL.Raw)
		}
		pr.URL.Raw = strings.Replace(pr.URL.Raw, match[0], val, -1)
	}
	// check Headers
	for _, hdr := range pr.Header {
		match = r.FindStringSubmatch(hdr.Value)
		if len(match) > 1 {
			exists, val := check_env_var_set(match[1], env)
			if !exists {
				return fmt.Errorf("%s env variable for Hedaer: %s is not set", match[1], hdr.Key)
			}
			hdr.Value = strings.Replace(hdr.Value, match[0], val, -1)
		}
	}
	// check body
	match = r.FindStringSubmatch(pr.Body.Raw)
	if len(match) > 1 {
		exists, val := check_env_var_set(match[1], env)
		if !exists {
			return fmt.Errorf("%s env variable for Body: %s is not set", match[1], pr.Body.Raw)
		}
		pr.Body.Raw = strings.Replace(pr.Body.Raw, match[0], val, -1)
	}
	return nil
}
