package request

import (
	"net/http"

	execconf "github.com/csabakissmalta/tpee/exec"
	postman "github.com/csabakissmalta/tpee/postman"
	task "github.com/csabakissmalta/tpee/task"
)

// Mostly the executable request's lifecycle-related operations

func ComposeHttpRequest(t *task.Task, p *postman.Request, env []*execconf.ExecEnvironmentElem) (*task.Task, error) {
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

	r_res, e := http.NewRequest(p.Method, p.URL.Raw, nil)
	if e != nil {
		return nil, e
	}

	task.WithRequest(r_res)(t)

	return t, nil
}
