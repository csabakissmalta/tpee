package request

import (
	"log"
	"regexp"

	timeline "github.com/csabakissmalta/tpee/timeline"
)

func validate_and_substitute_feed_type(in *string, r_var *regexp.Regexp, fds []*timeline.Feed) (*string, error) {
	// match := r_var.FindStringSubmatch(*in)
	// if len(match) > 1 {
	// 	varname := match[1]
	// 	var ch chan interface{}
	// 	for _, feed := range fds {
	// 		if varname == feed.Name {
	// 			ch = feed.Value
	// 			break
	// 		}
	// 	}
	// 	if ch != nil {
	// 		repl := <-ch
	// 		*in = strings.Replace(*in, match[1], repl.(string), -1)
	// 	} else {
	// 		e := fmt.Errorf("the variable match has failed, no feed value is substituted")
	// 		return in, e
	// 	}
	// }
	match := r.FindStringSubmatch(*in)
	for i, name := range r.SubexpNames() {
		if i > 0 && i <= len(match) {
			log.Printf("MATHED: var: %s :: value: %s", name, match[i])
		}
	}
	return in, nil
}
