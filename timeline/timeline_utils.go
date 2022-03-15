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
	f_name = f_name[:len(f_name)-1]
	f_extension := strings.Split(f_name, ".")[1]
	switch f_extension {
	case "csv":
		// var rec_index int
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

		// shuffle the feed records - they tend to be ordered in some way
		rec = rec[1:]
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(rec), func(i, j int) {
			rec[i], rec[j] = rec[j], rec[i]
		})

		for i := 0; i < dim; i++ {
			d := make(map[string]string)
			for idx, hkey := range csv_header {
				d[hkey] = rec[i][idx]
			}
			f.Value <- d
		}
	case "json":
	case "txt":
	default:
		return f
	}
	return f
}

func calc_periods(dur int, step int, er *execconf.ExecRequestsElem, rq *postman.Request) chan *task.Task {
	// the count of markers is (duration - delay) * frequency
	m_count := (dur - er.DelaySeconds) * er.Frequency

	ch := make(chan *task.Task, m_count)
	for i := er.DelaySeconds * er.Frequency; i < int(m_count); i++ {
		curr_step := i * step
		ch <- task.New(
			task.WithPlannedExecTimeNanos(curr_step),
			task.WithLabel(er.Name),
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

func validate_and_substitute(src string, rgx *regexp.Regexp, env []*execconf.ExecEnvironmentElem) (string, error) {
	match := rgx.FindAllStringSubmatch(src, -1)
	var rpl string = src
	if len(match) > 0 {
		for _, mtch := range match {
			exists, val := check_env_var_set(mtch[1], env)
			if !exists {
				return "", fmt.Errorf("%s env variable for URL: %s is not set", mtch[1], src)
			}
			rpl = strings.Replace(rpl, mtch[0], val, -1)
			// log.Println(rpl)
		}
		return rpl, nil
	}
	return src, nil
}
