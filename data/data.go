package data

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	execconfig "github.com/csabakissmalta/tpee/exec"
	"github.com/csabakissmalta/tpee/store"
)

const (
	EXTR_TARGET_BODY    = "body"
	EXTR_TARGET_HEADER  = "header"
	EXTR_TARGET_SESSION = "cookies"
)

func ExtractDataFromResponse(resp *http.Response, rule *execconfig.ExecRequestsElemDataPersistenceDataOutElem, storage store.Store) {
	// determine content type, based on response header
	if rule.Target == EXTR_TARGET_BODY {
		raw_ctype := resp.Header.Get("Content-Type")
		ctype := strings.Split(raw_ctype, ";")[0]
		switch {
		case strings.Contains(ctype, "json"):
			body, _ := io.ReadAll(resp.Body)

			// log.Println(string(body))

			var intf_ interface{}
			e := json.Unmarshal(body, &intf_)
			if e != nil {
				log.Println("DATA EXTRACTION ERROR (Unmarshal bytes into interface):", e.Error())
				// return ""
			}

			// err := json.NewDecoder(resp.Body).Decode(&intf_)
			// if err != nil {
			// 	log.Println("DATA EXTRACTION ERROR (Read body):", err.Error())
			// }

			intf, ok := intf_.(map[string]interface{})
			if !ok {
				log.Fatal("Error converting unmarshalled data to map")
			}
			// to_push := extractFromJSONBody(body, rule.Name)
			to_push := intf[rule.Name].(string)
			// call to the store to save the data
			storage.SaveData(to_push, rule)
		default:
			log.Println(ctype)
		}
	} else if rule.Target == EXTR_TARGET_HEADER {
		// to make it more precise here
		mems := strings.Split(rule.ContentType, "§")
		ctype := resp.Header.Get(mems[0])
		if len(mems) > 1 {
			regex_ptr := mems[1]
			matchmap := RegexpIt(regex_ptr, ctype)
			for key, val := range matchmap {
				if key == rule.Name {
					to_push := val
					log.Printf("tpee: %s", val)
					// call to the store to save the data
					storage.SaveData(to_push, rule)
				}
			}
		}
	}

}

func RegexpIt(regEx, src string) (rgxmap map[string]string) {
	var compRegEx = regexp.MustCompile(regEx)
	match := compRegEx.FindAllStringSubmatch(src, -1)
	rgxmap = make(map[string]string)
	for i := range compRegEx.SubexpNames() {
		if i < len(match) {
			rgxmap[match[i][2]] = match[i][3]
		}
	}
	return rgxmap
}

// func extractFromJSONBody(b []byte, key string) string {
// 	intf := make(map[string]interface{})
// 	e := json.Unmarshal(b, &intf)
// 	if e != nil {
// 		log.Println("DATA EXTRACTION ERROR (Unmarshal bytes into map):", e.Error())
// 		return ""
// 	}
// 	result := intf[key].(string)
// 	return result
// }
