package timeline

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"

	execconf "github.com/csabakissmalta/tpee/exec"
	postman "github.com/csabakissmalta/tpee/postman"
	task "github.com/csabakissmalta/tpee/task"
)

var r = regexp.MustCompile(`[\{]{2}(.{1,32})[\}]{2}`)

type Feed struct {
	Name  string
	Value chan interface{}
}

func load_feed(dim int, e *execconf.ExecEnvironmentElem) *Feed {
	f := &Feed{
		Name:  e.Key,
		Value: make(chan interface{}, dim),
	}
	// determine feed file type and load the file accordingly
	f_name := strings.Split(e.Value, "|")[1]
	f_extension := strings.Split(f_name, ".")[1]
	f_extension = f_extension[:len(f_extension)-1]
	switch f_extension {
	case "csv":
		var rec_index int
		file, err := os.Open(f_name)
		if err != nil {
			log.Fatalf("DATA READ ERROR: cannot open feed file %s. %s", f_name, err.Error())
		}
		csvReader := csv.NewReader(file)
		rec, err := csvReader.ReadAll()
		if err != nil {
			log.Fatalf("DATA READ ERROR: cannot read file %s. %s", f_name, err.Error())
		}
		csv_header := rec[0]
		for idx, hkey := range csv_header {
			if hkey == e.Key {
				rec_index = idx
				break
			}
		}
		// shuffle the feed records - they tend to be ordered in some way
		rec = rec[1:]
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(rec), func(i, j int) {
			rec[i], rec[j] = rec[j], rec[i]
		})

		for i := 0; i < dim; i++ {
			f.Value <- rec[i][rec_index]
		}
	case "json":
	case "txt":
	default:
		return f
	}
	return f
}

func calc_periods(dur int, er *execconf.ExecRequestsElem, rq *postman.Request) chan *task.Task {
	// the count of markers is (duration - delay) * frequency
	m_count := (dur - er.DelaySeconds) * er.Frequency

	// the step between the markers
	second := time.Second
	convers := int(second / time.Nanosecond)
	step := int(convers / er.Frequency)

	ch := make(chan *task.Task, m_count)
	for i := er.DelaySeconds * er.Frequency; i < int(m_count); i++ {
		curr_step := i * step
		ch <- task.New(
			task.WithPlannedExecTimeNanos(curr_step),
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

func load_feeds_if_required(dim int, env []*execconf.ExecEnvironmentElem) []*Feed {
	fds := []*Feed{}
	for _, envElem := range env {
		if *envElem.Type == execconf.FEED_VALUE {
			fd := load_feed(dim, envElem)
			fds = append(fds, fd)
		}
	}
	return fds
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
