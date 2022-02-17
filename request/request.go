package request

import (
	"log"
	"net/http"
	"regexp"

	execconf "github.com/csabakissmalta/tpee/exec"
	postman "github.com/csabakissmalta/tpee/postman"
	task "github.com/csabakissmalta/tpee/task"
	"github.com/csabakissmalta/tpee/timeline"
)

// Mostly the executable request's lifecycle-related operations

// Regex to get the substitution variable name (max length 30 characters)
var r = regexp.MustCompile(`(?P<WHOLE>[\+]{1}(?P<FEED_VAR>.{1,30})[|]{1}.+[\+])`)

func ComposeHttpRequest(t *task.Task, p postman.Request, env []*execconf.ExecEnvironmentElem, fds []*timeline.Feed) (*task.Task, error) {
	// check the postman request
	// URL
	// URL.Raw
	_, err := validate_and_substitute_feed_type(&p.URL.Raw, r, fds)
	if err != nil {
		log.Printf("SUBSTITUTE FEED VAR ERROR: %s", err.Error())
	}
	// log.Println(out)
	// p.URL.Raw = out

	// Body if Urlencoded
	if len(p.Body.Urlencoded) > 0 {
		for _, b := range p.Body.Urlencoded {
			_, err := validate_and_substitute_feed_type(&b.Value, r, fds)
			if err != nil {
				log.Printf("SUBSTITUTE FEED VAR ERROR: %s", err.Error())
			}
			// log.Println(out)
			// b.Value = out
		}
	}
	if len(p.Body.Urlencoded) > 0 {
		for _, b := range p.Body.Urlencoded {
			log.Println(b.Value)
		}
	}

	r_res, e := http.NewRequest(p.Method, p.URL.Raw, nil)
	if e != nil {
		return nil, e
	}

	task.WithRequest(r_res)(t)

	return t, nil
}

// log.Println(r.Method)
// log.Println(r.Auth)
// if len(r.Body.Urlencoded) > 0 {
// 	for _, b := range r.Body.Urlencoded {
// 		log.Println(b.Value)
// 	}
// }
// if len(r.Body.Formdata) > 0 {
// 	for _, f := range r.Body.Formdata {
// 		log.Println(f.Value)
// 	}
// }
// log.Println(r.Header)
// log.Println(r.URL)
