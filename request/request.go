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
	log.Println(r.Body)
	log.Println(r.Header)
	log.Println(r.URL)

	// r_res := &http.Request{
	// 	Method: r.Method,
	// 	URL: ,
	// }
	return nil
}
