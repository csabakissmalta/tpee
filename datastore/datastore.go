package datastore

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	execconfig "github.com/csabakissmalta/tpee/exec"
)

const (
	EXTR_TARGET_BODY         = "body"
	EXTR_TARGET_HEADER       = "header"
	EXTR_TARGET_SESSION      = "cookies"
	OUT_CHANNELS_BUFFER_SIZE = 1000
	IN_CHANNELS_BUFFER_SIZE  = 1000
)

var dataIn chan *InUnsorted

// Data type, stores a channel and name
type Data struct {
	// Identifier of the data stored in the channel
	Name string

	// The data channel
	Queue chan interface{}
}

// Data coming in needs to be sorted and controlled based on key
type InUnsorted struct {
	// Name - the key/channel name wehere it belongs
	Name string

	// Data
	In interface{}
}

// Datastore controller
type DataBroadcaster struct {
	// Channels containing sorted data
	DataOut []*Data
}

type Option func(*DataBroadcaster)

func New(option ...Option) *DataBroadcaster {
	db := &DataBroadcaster{}
	for _, o := range option {
		o(db)
	}
	dataIn = make(chan *InUnsorted, IN_CHANNELS_BUFFER_SIZE)
	return db
}

func WithDataOutSocketNames(dss []string) Option {
	d := []*Data{}

	for _, nm := range dss {
		_dq := &Data{
			Name:  nm,
			Queue: make(chan interface{}, OUT_CHANNELS_BUFFER_SIZE),
		}
		d = append(d, _dq)
	}

	return func(db *DataBroadcaster) {
		db.DataOut = d
	}
}

func (db *DataBroadcaster) getDataChannelByName(name string) *Data {
	for _, d := range db.DataOut {
		if name == d.Name {
			return d
		}
	}
	return nil
}

func (db *DataBroadcaster) StartConsumingDataIn() {
	for {
		select {
		case in := <-dataIn:
			ch_obj := db.getDataChannelByName(in.Name)
			if ch_obj != nil {
				ch_obj.Queue <- in.In

				// this is the ringbuffer-like trait
				// if the buffer is full, it will remove an element to unblock the queue
				if len(ch_obj.Queue) == OUT_CHANNELS_BUFFER_SIZE {
					<-ch_obj.Queue
				}
			} else {
				log.Println("DATA EXTRACTION ERROR: out channel doesn't exist")
			}
		default:
			continue
		}
	}
}

func PushDataIn(d *InUnsorted) {
	dataIn <- d
}

func ExtractDataFromResponse(resp *http.Response, extr_rules []*execconfig.ExecRequestsElemDataPersistenceDataOutElem) {
	// determine content type, based on response header
	for _, rule := range extr_rules {
		if rule.Target == EXTR_TARGET_BODY {
			raw_ctype := resp.Header.Get("Content-Type")
			ctype := strings.Split(raw_ctype, ";")[0]
			switch {
			case strings.Contains(ctype, "json"):
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Println("DATA EXTRACTION ERROR:", err.Error())
				}
				defer resp.Body.Close()

				to_push := extractFromJSONBody(body, *rule.Name)

				PushDataIn(&InUnsorted{
					Name: *rule.Name,
					In:   to_push,
				})
			default:
				log.Println(ctype)
			}
		} else if rule.Target == EXTR_TARGET_HEADER {
			// to make it more precise here
			// log.Println("ERROR: Extraction from", rule.Target, "is not implmented yet")
			mems := strings.Split(rule.ContentType, "§")
			ctype := resp.Header.Get(mems[0])
			if len(mems) > 1 {
				regex_ptr := mems[1]
				matchmap := RegexpIt(regex_ptr, ctype)
				for key, val := range matchmap {
					if key == mems[0] {
						to_push := val
						log.Printf("tpee: %s", val)
						PushDataIn(&InUnsorted{
							Name: *rule.Name,
							In:   to_push,
						})
					}
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

func extractFromJSONBody(b []byte, key string) string {
	intf := make(map[string]interface{})
	e := json.Unmarshal(b, &intf)
	if e != nil {
		log.Println("DATA EXTRACTION ERROR:", e.Error())
		return ""
	}
	result := intf[key].(string)
	return result
}
