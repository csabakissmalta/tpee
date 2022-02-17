package request

import (
	"regexp"
	"strings"

	timeline "github.com/csabakissmalta/tpee/timeline"
)

func validate_and_substitute_feed_type(in *string, r_var *regexp.Regexp, fds []*timeline.Feed) (string, error) {
	match := r.FindStringSubmatch(*in)
	if len(match) == 0 {
		return *in, nil
	}
	var feed_varname string
	var env_var_to_replace string
	var env_var_replace_string interface{}
	for i, name := range r.SubexpNames() {
		if i > 0 && i <= len(match) {
			// log.Printf("MATCHED: var: %s :: value: %s", name, match[i])
			if name == "FEED_VAR" {
				feed_varname = match[i]
				// log.Println("feed_varname: ", feed_varname)
			} else if name == "WHOLE" {
				env_var_to_replace = match[i]
			}
		}
	}
	for _, feed := range fds {
		if feed_varname == feed.Name {
			env_var_replace_string = <-feed.Value
		}
	}
	out := strings.Replace(*in, env_var_to_replace, env_var_replace_string.(string), -1)
	return out, nil
}
