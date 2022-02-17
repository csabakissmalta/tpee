package request

import (
	"bytes"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	execconf "github.com/csabakissmalta/tpee/exec"
	postman "github.com/csabakissmalta/tpee/postman"
	task "github.com/csabakissmalta/tpee/task"
	"github.com/csabakissmalta/tpee/timeline"
)

// Mostly the executable request's lifecycle-related operations

// Regex to get the substitution variable name (max length 30 characters)
var r = regexp.MustCompile(`(?P<WHOLE>[\+]{1}(?P<FEED_VAR>.{1,30})[|]{1}.+[\+])`)

func ComposeHttpRequest(t *task.Task, p postman.Request, env []*execconf.ExecEnvironmentElem, fds []*timeline.Feed) (*task.Task, error) {
	var req_url string
	var req_method string = p.Method
	// body_urlencoded

	// check the postman request
	// URL
	// URL.Raw
	out, err := validate_and_substitute_feed_type(&p.URL.Raw, r, fds)
	if err != nil {
		log.Printf("SUBSTITUTE FEED VAR ERROR: %s", err.Error())
	}
	req_url = out
	log.Println(req_url)

	// Body if Urlencoded
	if len(p.Body.Urlencoded) > 0 {
		body_urlencoded := url.Values{}
		for _, b := range p.Body.Urlencoded {
			out, err := validate_and_substitute_feed_type(&b.Value, r, fds)
			if err != nil {
				log.Printf("SUBSTITUTE FEED VAR ERROR: %s", err.Error())
			}
			body_urlencoded.Set(b.Key, out)
		}
		encoded_data := body_urlencoded.Encode()
		r_res, e := http.NewRequest(req_method, req_url, strings.NewReader(encoded_data))
		if e != nil {
			return nil, e
		}
		r_res.Header.Add("Content-Length", strconv.Itoa(len(encoded_data)))
		task.WithRequest(r_res)(t)
	}

	// Body if Raw
	if len(p.Body.Raw) > 0 {
		out, err := validate_and_substitute_feed_type(&p.Body.Raw, r, fds)
		if err != nil {
			log.Printf("SUBSTITUTE FEED VAR ERROR: %s", err.Error())
		}
		r_res, e := http.NewRequest(req_method, req_url, bytes.NewBuffer([]byte(out)))
		if e != nil {
			return nil, e
		}
		task.WithRequest(r_res)(t)
	}

	return t, nil
}
