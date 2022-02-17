package request

import (
	"regexp"
	"strings"

	timeline "github.com/csabakissmalta/tpee/timeline"
)

// func generate_auth_header_value() string {
// 	return ""
// }

func validate_and_substitute_feed_type(in *string, r_var *regexp.Regexp, fds []*timeline.Feed) (string, error) {
	match := r.FindStringSubmatch(*in)
	if len(match) == 0 {
		return *in, nil
	}
	var ch chan interface{}
	var feed_varname string
	var env_var_to_replace string
	var env_var_replace_string interface{}
	for i, name := range r.SubexpNames() {
		if i > 0 && i <= len(match) {
			if name == "FEED_VAR" {
				feed_varname = match[i]
			} else if name == "WHOLE" {
				env_var_to_replace = match[i]
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
