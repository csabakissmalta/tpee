package request

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	datastore "github.com/csabakissmalta/tpee/datastore"
	execconf "github.com/csabakissmalta/tpee/exec"
	postman "github.com/csabakissmalta/tpee/postman"
	task "github.com/csabakissmalta/tpee/task"
	timeline "github.com/csabakissmalta/tpee/timeline"
)

// Mostly the executable request's lifecycle-related operations

// Regex to get the substitution variable name (max length 30 characters)
var r = regexp.MustCompile(`(?P<WHOLE>[\+]{1}(?P<FEED_VAR>[a-z0-9-_]{1,30})[|]{1}.+[\+])`)

// Regex to get the substitution variable for the datastore
var rds = regexp.MustCompile(`(?P<WHOLE>[<-]{2}(?P<CHAN>[a-z0-9-_]{1,30})[<-]{2})`)

func ComposeHttpRequest(t *task.Task, p postman.Request, env []*execconf.ExecEnvironmentElem, fds []*timeline.Feed, ds *datastore.DataBroadcaster) (*task.Task, error) {
	var req_url string
	var req_method string = p.Method
	var r_res *http.Request
	// body_urlencoded

	// check the postman request
	// --- URL.Raw ---
	out, err := validate_and_substitute_feed_type(&p.URL.Raw, r, rds, fds, ds)
	if err != nil {
		log.Printf("SUBSTITUTE FEED VAR ERROR: %s", err.Error())
	}
	req_url = out

	// --- Body if Urlencoded ---
	if len(p.Body.Urlencoded) > 0 {
		body_urlencoded := url.Values{}
		for _, b := range p.Body.Urlencoded {
			out, err := validate_and_substitute_feed_type(&b.Value, r, rds, fds, ds)
			if err != nil {
				log.Printf("SUBSTITUTE FEED VAR ERROR: %s", err.Error())
			}
			body_urlencoded.Set(b.Key, out)
		}
		encoded_data := body_urlencoded.Encode()
		r_res, err = http.NewRequest(req_method, req_url, strings.NewReader(encoded_data))
		if err != nil {
			return nil, err
		}
		r_res.Header.Add("Content-Length", strconv.Itoa(len(encoded_data)))
	}

	// --- Body if Raw ---
	if len(p.Body.Raw) > 0 {
		out, err := validate_and_substitute_feed_type(&p.Body.Raw, r, rds, fds, ds)
		if err != nil {
			log.Printf("SUBSTITUTE FEED VAR ERROR: %s", err.Error())
		}
		r_res, err = http.NewRequest(req_method, req_url, bytes.NewBuffer([]byte(out)))
		if err != nil {
			return nil, err
		}
	}

	// --- Body if Form ---
	if len(p.Body.Formdata) > 0 {
		var body bytes.Buffer
		wtr := multipart.NewWriter(&body)
		for _, fd := range p.Body.Formdata {
			fw, err := wtr.CreateFormField(fd.Key)
			if err != nil {
				log.Printf("ERROR: %s", err.Error())
			}
			out, err := validate_and_substitute_feed_type(&fd.Value, r, rds, fds, ds)
			if err != nil {
				log.Printf("SUBSTITUTE FEED VAR ERROR: %s", err.Error())
			}
			_, err = io.Copy(fw, strings.NewReader(out))
			if err != nil {
				log.Printf("ERROR: %s", err.Error())
			}
		}
		wtr.Close()
		r_res, err = http.NewRequest(req_method, req_url, bytes.NewReader(body.Bytes()))
		if err != nil {
			return nil, err
		}
		r_res.Header.Set("Content-Type", wtr.FormDataContentType())
	}

	// --- Headers ---
	for _, hdr := range p.Header {
		r_res.Header.Add(hdr.Key, hdr.Value)
	}

	if p.Auth.Type != "" {
		switch p.Auth.Type {
		case "basic":
			var uname string
			var pword string
			for _, autAttr := range p.Auth.Basic {
				if autAttr.Key == "username" {
					uname = autAttr.Value.(string)
				} else if autAttr.Key == "password" {
					pword = autAttr.Value.(string)
				}
			}
			r_res.SetBasicAuth(uname, pword)
		case "bearer":
			log.Printf("ERROR: Auth type %s is not implemented yet", p.Auth.Type)
		default:
			log.Printf("ERROR: Auth type %s is not implemented yet", p.Auth.Type)
		}
	}

	task.WithRequest(r_res)(t)

	return t, nil
}
