package datastore

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"

	execconfig "github.com/csabakissmalta/tpee/exec"
)

const (
	EXTR_TARGET_BODY         = "body"
	EXTR_TARGET_HEADER       = "header"
	OUT_CHANNELS_BUFFER_SIZE = 100
	IN_CHANNELS_BUFFER_SIZE  = 10
)

var cases []reflect.SelectCase

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
	cases = []reflect.SelectCase{}

	for _, nm := range dss {
		_dq := &Data{
			Name:  nm,
			Queue: make(chan interface{}, OUT_CHANNELS_BUFFER_SIZE),
		}
		d = append(d, _dq)
		_cs := reflect.SelectCase{
			Dir:  reflect.SelectSend,
			Chan: reflect.ValueOf(_dq.Queue),
			// Send: ,
		}
		cases = append(cases, _cs)
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
				log.Println(in.In)
				ch_obj.Queue <- in.In
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
			log.Println(rule.Target)
		}
	}
}

func extractFromJSONBody(b []byte, key string) string {
	intf := make(map[string]interface{})
	e := json.Unmarshal(b, &intf)
	if e != nil {
		log.Println("DATA EXTRACTION ERROR:", e.Error())
	}
	result := intf[key].(string)
	return result
}
