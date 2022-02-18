package datastore

import (
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

func (db *DataBroadcaster) StartConsumingDataIn() {
	for {
		select {
		case in := <-dataIn:
			log.Println(in.Name)
		default:
			continue
		}
	}
}

func PushDataIn(d *InUnsorted) {
	dataIn <- d
}

func ExtractDataFromResponse(resp *http.Response, extr_rules []*execconfig.ExecRequestsElemDataPersistenceDataOutElem) (interface{}, error) {
	// determine content type, based on response header
	for _, rule := range extr_rules {
		if rule.Target == EXTR_TARGET_BODY {
			raw_ctype := resp.Header.Get("Content-Type")
			ctype := strings.Split(raw_ctype, ";")[0]
			switch {
			case strings.Contains(ctype, "json"):
				log.Println("JSON TYPE")
				PushDataIn(&InUnsorted{
					Name: *rule.Name,
					In:   ctype,
				})
			default:
				log.Println(ctype)
			}
		} else if rule.Target == EXTR_TARGET_HEADER {
			// to make it more precise here
			log.Println(rule.Target)
		}
	}
	return nil, nil
}
