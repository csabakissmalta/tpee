package request

import (
	"log"
	"regexp"
	"strings"

	datastore "github.com/csabakissmalta/tpee/datastore"
	timeline "github.com/csabakissmalta/tpee/timeline"
)

// func generate_auth_header_value() string {
// 	return ""
// }

func validate_and_substitute_feed_type(in *string, r_var *regexp.Regexp, r_ds *regexp.Regexp, fds []*timeline.Feed, ds *datastore.DataBroadcaster) (string, error) {
	match_feed := r_var.FindStringSubmatch(*in)
	match_channel := r_ds.FindStringSubmatch(*in)

	// log.Println("len(match_feed)", len(match_feed))
	// log.Println("len(match_channel)", len(match_channel))

	var ch chan interface{}
	var feed_varname string
	var env_var_to_replace string
	var env_var_replace_string interface{}

	if len(match_feed) > 0 {
		for i, name := range r.SubexpNames() {
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
		env_var_replace_string = <-ch
		out := strings.Replace(*in, env_var_to_replace, env_var_replace_string.(string), -1)
		ch <- env_var_replace_string
		return out, nil
	}

	if len(match_channel) > 0 {
		for i, name := range r.SubexpNames() {
			if i > 0 && i <= len(match_channel) {
				if name == "CHAN" {
					log.Println(name)
					feed_varname = match_channel[i]
				} else if name == "WHOLE" {
					log.Println(name)
					env_var_to_replace = match_channel[i]
				}
			}
		}
		for _, chans := range ds.DataOut {
			if feed_varname == chans.Name {
				ch = chans.Queue
				break
			}
		}
		env_var_replace_string = <-ch
		out := strings.Replace(*in, env_var_to_replace, env_var_replace_string.(string), -1)
		// ch <- env_var_replace_string
		log.Println(out)
		return out, nil
	}
	return *in, nil
}
