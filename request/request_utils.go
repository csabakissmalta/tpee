package request

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	datastore "github.com/csabakissmalta/tpee/datastore"
	execconf "github.com/csabakissmalta/tpee/exec"
	sessionstore "github.com/csabakissmalta/tpee/sessionstore"
	timeline "github.com/csabakissmalta/tpee/timeline"
)

// func generate_auth_header_value() string {
// 	return ""
// }

// whichDataStore is a function to determine from where to retrieve the data
func whichDataStore(name string, dp []*execconf.ExecRequestsElemDataPersistenceDataInElem) (string, bool) {
	for _, d := range dp {
		if name == d.Name {
			return d.Storage.(string), d.Retention
		}
	}
	return "", false
}

func validate_and_substitute(in *string, r_var *regexp.Regexp, r_ds *regexp.Regexp, r_ss *regexp.Regexp, fds []*timeline.Feed, ds *datastore.DataBroadcaster, ss *sessionstore.Session, dp []*execconf.ExecRequestsElemDataPersistenceDataInElem) (string, error) {
	match_feed := r_var.FindStringSubmatch(*in)
	match_data_in_storage := r_ds.FindStringSubmatch(*in)
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
		fmt.Printf(":: Conversion: %v", elem)
		elem_map, ok := elem.(map[string]string)
		if !ok {
			return "", fmt.Errorf("conversion error: %v, %v", elem_map, elem)
		}

		env_var_replace_string = elem_map[feed_varname]
		// log.Println(env_var_replace_string)
		out := strings.Replace(*in, env_var_to_replace, env_var_replace_string, -1)
		ch <- elem
		return out, nil
	}

	// check DATA var match
	if len(match_data_in_storage) > 0 {
		for i, name := range r_ds.SubexpNames() {
			if i > 0 && i <= len(match_data_in_storage) {
				if name == "CHAN" {
					feed_varname = match_data_in_storage[i]
				} else if name == "WHOLE" {
					env_var_to_replace = match_data_in_storage[i]
				}
			}
		}
		// var ret bool = true

		datasource_in, retention := whichDataStore(feed_varname, dp)
		var elem interface{}
		switch datasource_in {
		case "data-store":
			elem = ds.RetrieveData(feed_varname, retention)
			if retention {
				ch <- env_var_replace_string
			}
		case "session-meta":
			elem = ss.RetrieveData(feed_varname, retention)
		default:
			// do nothing
		}

		env_var_replace_string = elem.(string)
		out := strings.Replace(*in, env_var_to_replace, env_var_replace_string, -1)

		// log.Println("************************")
		// log.Println("*** DATASTORE_IN:", datasource_in)
		// log.Println("*** ELEM:", elem)
		// log.Println("*** IN:", *in)
		// log.Println("*** OUT:", out)
		// log.Println("************************")

		return out, nil
	}

	// check SESSION var match
	if len(match_session) > 0 && ss != nil {
		var out string = *in

		for _, mtch := range match_session {
			for i, name := range r_ss.SubexpNames() {
				if name == "SESSIONVAR" {
					sessionvar_name = mtch[i]
				} else if name == "WHOLE" {
					env_var_to_replace = mtch[i]
				}
			}

			for _, c := range ss.ID.([]*http.Cookie) {
				if sessionvar_name == c.Name {
					env_var_replace_string = c.Value
				}
			}

			if env_var_replace_string != "" {
				out = strings.Replace(out, env_var_to_replace, env_var_replace_string, -1)
			}
		}
		return out, nil
	}
	return *in, nil
}
