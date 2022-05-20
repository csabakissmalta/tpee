package datastore

import (
	"log"

	"github.com/csabakissmalta/tpee/exec"
)

const (
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

	// Retetntion policy
	Retention bool
}

// Data coming in needs to be sorted and controlled based on key
type InUnsorted struct {
	// Name - the key/channel name wehere it belongs
	Name string

	// Data
	In interface{}

	// Retention policy
	Retetntion bool
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

// Store interface impl
func (db *DataBroadcaster) SaveData(extracted interface{}, rule *exec.ExecRequestsElemDataPersistenceDataOutElem) {
	dataIn <- &InUnsorted{
		Name:       *rule.Name,
		In:         extracted,
		Retetntion: rule.Retention,
	}
}

// Store interafce impl
func (db *DataBroadcaster) RetrieveData(name string) interface{} {
	return nil
}
