package request

import (
	"log"
	"net/http"

	execconf "github.com/csabakissmalta/tpee/exec"
	"github.com/csabakissmalta/tpee/postman"
)

// Mostly the executable request's lifecycle-related operations

func ComposeHttpRequest(r *postman.Request, env []*execconf.ExecEnvironmentElem) *http.Request {
	log.Println(r.Method)
	log.Println(r.Auth)
	if len(r.Body.Urlencoded) > 0 {
		for _, b := range r.Body.Urlencoded {
			log.Println(b.Value)
		}
	}
	if len(r.Body.Formdata) > 0 {
		for _, f := range r.Body.Formdata {
			log.Println(f.Value)
		}
	}
	log.Println(r.Header)
	log.Println(r.URL)

	// r_res := &http.Request{
	// 	Method: r.Method,
	// 	URL: ,
	// }
	return nil
}
