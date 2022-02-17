package request

import (
	"regexp"

	timeline "github.com/csabakissmalta/tpee/timeline"
)

func validate_and_substitute_feed_type(in string, r *regexp.Regexp, fds []*timeline.Feed) error {
	// match := r.FindStringSubmatch(in)
	// if len(match) > 1 {
	// 	repl := <-f
	// 	in = strings.Replace(in, match[1], repl.(string), -1)
	// }
	return nil
}
