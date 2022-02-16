package timeline

import (
	"fmt"
	"regexp"
	"strings"

	execconf "github.com/csabakissmalta/tpee/exec"

	"github.com/csabakissmalta/tpee/postman"
	"github.com/csabakissmalta/tpee/task"
)

var r = regexp.MustCompile(`[\{]{2}(.{1,32})[\}]{2}`)

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
		if *envElem.Type == execconf.GENERATED_VALUE {
			continue
		}
		if envElem.Key == vname && len(envElem.Value) > 0 {
			return true, envElem.Value
		}
	}
	return false, ""
}

func validate_and_substitute(src *string, rgx *regexp.Regexp, env []*execconf.ExecEnvironmentElem) error {
	match := rgx.FindStringSubmatch(*src)
	if len(match) > 1 {
		exists, val := check_env_var_set(match[1], env)
		if !exists {
			return fmt.Errorf("%s env variable for URL: %s is not set", match[1], *src)
		}
		*src = strings.Replace(*src, match[0], val, -1)
	}
	return nil
}

func check_postman_request_and_validate_requirements(pr *postman.Request, env []*execconf.ExecEnvironmentElem) error {
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
