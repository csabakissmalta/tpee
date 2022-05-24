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
	"time"

	datastore "github.com/csabakissmalta/tpee/datastore"
	execconf "github.com/csabakissmalta/tpee/exec"
	postman "github.com/csabakissmalta/tpee/postman"
	"github.com/csabakissmalta/tpee/sessionstore"
	task "github.com/csabakissmalta/tpee/task"
	timeline "github.com/csabakissmalta/tpee/timeline"
)

// Mostly the executable request's lifecycle-related operations

// --------------- REGEXP DEFINITIONS ---------------
// Regex to get the substitution variable name (max length 30 characters)
var r = regexp.MustCompile(`(?P<WHOLE>[\+]{1}(?P<FEED_VAR>[a-z0-9-_]{1,30})[|]{1}.+[\+])`)

// Regex to get the substitution variable for the datastore
var rds = regexp.MustCompile(`(?P<WHOLE>[\<]{1}(?P<CHAN>[a-z0-9\-_]{1,30})[\>]{1})`)

// Regex to get the substitution variable for SESSION variables
var rss = regexp.MustCompile(`(?P<WHOLE>[\<]{1}(?P<SESSIONVAR>[A-Za-z0-9\-_]{1,30})[\>]{1})`)

// ---------------------------------------------------

// does the request require session?
func isSessionRequired(dp []*execconf.ExecRequestsElemDataPersistenceDataInElem) bool {
	for _, d := range dp {
		if d.Storage == "session-meta" {
			return true
		}
	}
	return false
}

func ComposeHttpRequest(t *task.Task, p postman.Request, dp []*execconf.ExecRequestsElemDataPersistenceDataInElem, env []*execconf.ExecEnvironmentElem, fds []*timeline.Feed, ds *datastore.DataBroadcaster, ss *sessionstore.Store) (*task.Task, error) {
	var req_url string
	var req_method string = p.Method
	var r_res *http.Request
	var sess *sessionstore.Session
	// body_urlencoded

	if isSessionRequired(dp) {
		for {
			sess = <-ss.SessionOut
			if time.Since(sess.Created) < sessionstore.SESSION_VALIDITY {
				break
			}
		}
	}
	// get a session, if required, based on the data in props

	// check the postman request
	// --- URL.Raw ---
	out, err := validate_and_substitute(&p.URL.Raw, r, rds, rss, fds, ds, sess, dp)
	if err != nil {
		log.Printf("SUBSTITUTE VAR ERROR: %s", err.Error())
	}
	req_url = out

	// --- Body if Urlencoded ---
	if len(p.Body.Urlencoded) > 0 {
		body_urlencoded := url.Values{}
		for _, b := range p.Body.Urlencoded {
			out, err := validate_and_substitute(&b.Value, r, rds, rss, fds, ds, sess, dp)
			if err != nil {
				log.Printf("SUBSTITUTE VAR ERROR: %s", err.Error())
			}
			body_urlencoded.Set(b.Key, out)
		}
		encoded_data := body_urlencoded.Encode()
		r_res, err = http.NewRequest(req_method, req_url, strings.NewReader(encoded_data))
		if err != nil {
			return nil, err
		}
		r_res.Header.Add("Content-Length", strconv.Itoa(len(encoded_data)))
	} else if len(p.Body.Raw) > 0 {
		out, err := validate_and_substitute(&p.Body.Raw, r, rds, rss, fds, ds, sess, dp)
		if err != nil {
			log.Printf("SUBSTITUTE VAR ERROR: %s", err.Error())
		}
		r_res, err = http.NewRequest(req_method, req_url, bytes.NewBuffer([]byte(out)))
		if err != nil {
			return nil, err
		}
	} else if len(p.Body.Formdata) > 0 {
		var body bytes.Buffer
		wtr := multipart.NewWriter(&body)
		for _, fd := range p.Body.Formdata {
			fw, err := wtr.CreateFormField(fd.Key)
			if err != nil {
				log.Printf("ERROR: %s", err.Error())
			}
			out, err := validate_and_substitute(&fd.Value, r, rds, rss, fds, ds, sess, dp)
			if err != nil {
				log.Printf("SUBSTITUTE VAR ERROR: %s", err.Error())
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
	} else {
		r_res, err = http.NewRequest(req_method, req_url, nil)
		if err != nil {
			return nil, err
		}
	}

	// --- Headers ---
	for _, hdr := range p.Header {
		out, err := validate_and_substitute(&hdr.Value, r, rds, rss, fds, ds, sess, dp)
		if err != nil {
			log.Printf("SUBSTITUTE VAR ERROR: %s", err.Error())
		}
		// log.Println(out)
		r_res.Header.Add(hdr.Key, out)
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
			var token string
			for _, autAttr := range p.Auth.Bearer {
				if autAttr.Key == "username" {
					token = autAttr.Value.(string)
				}
			}
			out, err := validate_and_substitute(&token, r, rds, rss, fds, ds, sess, dp)
			if err != nil {
				log.Printf("SUBSTITUTE VAR ERROR: %s", err.Error())
			}
			r_res.Header.Add("Authorization", "Bearer"+out)
		case "noauth":
			// do nothing
		default:
			log.Printf("ERROR: Auth type %s is not implemented yet", p.Auth.Type)
		}
	}

	task.WithRequest(r_res)(t)

	return t, nil
}
