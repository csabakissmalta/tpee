package request

import (
	"net/http"
	"regexp"
	"strings"
	"time"

	datastore "github.com/csabakissmalta/tpee/datastore"
	sessionstore "github.com/csabakissmalta/tpee/sessionstore"
	timeline "github.com/csabakissmalta/tpee/timeline"
)

// func generate_auth_header_value() string {
// 	return ""
// }

func validate_and_substitute(in *string, r_var *regexp.Regexp, r_ds *regexp.Regexp, r_ss *regexp.Regexp, fds []*timeline.Feed, ds *datastore.DataBroadcaster, ss *sessionstore.Store) (string, error) {
	match_feed := r_var.FindStringSubmatch(*in)
	match_channel := r_ds.FindStringSubmatch(*in)
	match_session := r_ss.FindAllStringSubmatch(*in, -1)

	var ch chan interface{}
	var feed_varname string
	var sessionvar_name string
	var env_var_to_replace string
	var env_var_replace_string string

	// check FEED var match
	if len(match_feed) > 0 {
		for i, name := range r_var.SubexpNames() {
			if i > 0 && i <= len(match_feed) {
				if name == "FEED_VAR" {
					feed_varname = match_feed[i]
				} else if name == "WHOLE" {
					env_var_to_replace = match_feed[i]
				}
			}
		}
		for _, feed := range fds {
			if feed_varname == feed.Name {
				ch = feed.Value
				break
			}
		}
		elem := <-ch
		env_var_replace_string = elem.(map[string]string)[feed_varname]
		// log.Println(env_var_replace_string)
		out := strings.Replace(*in, env_var_to_replace, env_var_replace_string, -1)
		ch <- env_var_replace_string
		return out, nil
	}

	// check DATA var match
	if len(match_channel) > 0 {
		for i, name := range r_ds.SubexpNames() {
			if i > 0 && i <= len(match_channel) {
				if name == "CHAN" {
					feed_varname = match_channel[i]
				} else if name == "WHOLE" {
					env_var_to_replace = match_channel[i]
				}
			}
		}
		var ret bool = true
		for _, chans := range ds.DataOut {
			if feed_varname == chans.Name {
				ret = chans.Retention
				ch = chans.Queue
				break
			}
		}
		elem := <-ch
		env_var_replace_string = elem.(string)
		out := strings.Replace(*in, env_var_to_replace, env_var_replace_string, -1)

		if ret {
			ch <- env_var_replace_string
		}

		return out, nil
	}

	// check SESSION var match
	if len(match_session) > 0 {
		var out string = *in
		var sess *sessionstore.Session

		for _, mtch := range match_session {
			for i, name := range r_ss.SubexpNames() {
				if name == "SESSIONVAR" {
					sessionvar_name = mtch[i]
				} else if name == "WHOLE" {
					env_var_to_replace = mtch[i]
				}
			}

			if sess == nil {
				for {
					sess = <-ss.SessionOut
					if time.Since(sess.Created) < sessionstore.SESSION_VALIDITY {
						break
					}
				}
			}

			for _, c := range sess.ID.([]*http.Cookie) {
				if sessionvar_name == c.Name {
					env_var_replace_string = c.Value
				}
			}

			if env_var_replace_string != "" {
				out = strings.Replace(out, env_var_to_replace, env_var_replace_string, -1)
			}
		}
		ss.SessionIn <- sess
		return out, nil
	}
	return *in, nil
}
