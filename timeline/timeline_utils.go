package timeline

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strings"

	execconf "github.com/csabakissmalta/tpee/exec"
	postman "github.com/csabakissmalta/tpee/postman"
	task "github.com/csabakissmalta/tpee/task"
	"github.com/nats-io/nats.go"
)

var r = regexp.MustCompile(`[\{]{2}(.{1,32})[\}]{2}`)

type Feed struct {
	Name      string
	Value     chan interface{}
	Type      string
	NATSValue chan *nats.Msg
}

func load_feed(dim int, e *execconf.ExecEnvironmentElem) *Feed {
	f := &Feed{
		Name:  e.Key,
		Value: make(chan interface{}, dim),
	}
	// determine feed file type and load the file accordingly
	f_name := strings.Split(e.Value, "|")[1]
	f_name = f_name[:len(f_name)-1]
	f_extension_raw := strings.Split(f_name, ".")
	f_extension := f_extension_raw[len(f_extension_raw)-1]
	log.Println(":::EXTENSION :::", f_extension)
	switch f_extension {
	case "csv":
		// var rec_index int
		f.Type = "csv"
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
		// rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(rec), func(i, j int) {
			rec[i], rec[j] = rec[j], rec[i]
		})

		var i int = 0
		var cntr int = 0
		for {
			d := make(map[string]string)
			for idx, hkey := range csv_header {
				d[hkey] = rec[i][idx]
			}
			f.Value <- d
			i++
			cntr++
			if cntr == dim {
				break
			} else if i == len(rec) {
				i = 0
			}
		}
	case "nats":
		// setting up nats client channel consumes needs to be added here
		// -----
		log.Println(":::EXTENSION NATS :::", f_extension)
		return nil
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

	ch := make(chan *task.Task, 10000)
	var i int = er.DelaySeconds * er.Frequency
	// ; i < int(m_count); i++
	go func() {
		for {
			switch {
			case len(ch) < cap(ch):
				curr_step := i * step
				ch <- task.New(
					task.WithPlannedExecTimeNanos(curr_step),
					task.WithLabel(er.Name),
				)
				i++
			case i == int(m_count):
				return
			default:
				continue
				// do nothing
			}
		}
	}()
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

func load_feeds_if_required(dim int, env []*execconf.ExecEnvironmentElem, subs map[string]chan *nats.Msg) []*Feed {
	fds := []*Feed{}
	for _, envElem := range env {
		if *envElem.Type == execconf.FEED_VALUE {
			fd := load_feed(dim, envElem)

			if fd == nil {
				// that means, it is a NATS subscription
				ch_item := envElem.Key
				fd = &Feed{
					Name:      ch_item,
					NATSValue: subs[ch_item],
					Type:      "nats_msg",
				}
				fds = append(fds, fd)
			} else {
				fds = append(fds, fd)
			}
		}
	}
	return fds
}

func validate_and_substitute(src string, rgx *regexp.Regexp, env []*execconf.ExecEnvironmentElem) (string, error) {
	log.Println("timeline_utils.validate_and_substitute", src)
	match := rgx.FindAllStringSubmatch(src, -1)
	var rpl string = src
	if len(match) > 0 {
		for _, mtch := range match {
			exists, val := check_env_var_set(mtch[1], env)
			if !exists {
				return "", fmt.Errorf("%s env variable for URL: %s is not set", mtch[1], src)
			}
			rpl = strings.Replace(rpl, mtch[0], val, -1)
			log.Println("REPLACED: ", rpl)
		}
		return rpl, nil
	}
	return src, nil
}
