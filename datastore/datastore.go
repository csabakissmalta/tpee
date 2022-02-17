package datastore

import (
	"log"
	"net/http"
	"strings"
	// execconfig "github.com/csabakissmalta/tpee/exec"
)

func ExtractDataFromResponse(resp *http.Response) (interface{}, error) {
	// determine content type, based on response header
	raw_ctype := resp.Header.Get("Content-Type")
	ctype := strings.Split(raw_ctype, ";")[0]
	switch {
	case strings.Contains(ctype, "json"):
		log.Println("JSON TYPE")
	default:
		log.Println(ctype)
	}
	return nil, nil
}
