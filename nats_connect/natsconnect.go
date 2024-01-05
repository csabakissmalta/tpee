package natsconnect

import (
	"log"

	"github.com/nats-io/nats.go"
)

// The connect structure to hold config and client
type NATSClient struct {
	// NATS server URL
	ConnectionUrl string

	// Subjects
	Subjects []string

	// Creds path
	CredsPath string

	// connection
	conn *nats.Conn
}

// option type
type Option func(*NATSClient)

// create a client instance and return the pointer
func New(opts ...Option) *NATSClient {
	nc := &NATSClient{}
	for _, opt := range opts {
		opt(nc)
	}

	return nc
}

// Set connection Url (NATS_URL) - if not set, the default NATS url is set
// "nats://127.0.0.1:4222
func WithConnectionUrl(url string) Option {
	return func(ncl *NATSClient) {
		ncl.ConnectionUrl = url
	}
}

// The subject, what the client will listen on
func WithSubjects(subj []string) Option {
	return func(ncl *NATSClient) {
		ncl.Subjects = subj
	}
}

// If Authentication and Authorization is required for the NATS server,
// the creds file path needs to be provided
func WithCredsFilePath(cfp string) Option {
	return func(ncl *NATSClient) {
		ncl.CredsPath = cfp
	}
}

// -------------------- SUBSCRIBE ----------------------

// Connect to the NATS server and subscribe to the subjects
func (ncl *NATSClient) ConnectAndSubscribe() (map[string]chan *nats.Msg, error) {
	nc, err := nats.Connect(ncl.ConnectionUrl, nats.UserCredentials(ncl.CredsPath))
	if err != nil {
		return nil, err
	}
	channels := map[string]chan *nats.Msg{}
	// now subscribe
	for _, sbj := range ncl.Subjects {
		subs_chan := make(chan *nats.Msg)
		// sub, err := nc.ChanSubscribe(sbj, subs_chan)
		_, err := nc.QueueSubscribeSyncWithChan(sbj, "plab", subs_chan)
		if err != nil {
			log.Fatal("ERR: NATS subscription issue: ", err.Error())
		}
		// if err = sub.SetPendingLimits(1024*500, 1024*5000); err != nil {
		// 	log.Fatalf("Unable to set pending limits: %v", err)
		// }
		channels[sbj] = subs_chan
	}
	return channels, nil
}

// -------------------- PUBLISH ---------------------

// Connect to NATS server
func (ncl *NATSClient) connectAndReadyToPublish() error {
	nc, err := nats.Connect(ncl.ConnectionUrl, nats.UserCredentials(ncl.CredsPath))
	if err != nil {
		return err
	}
	ncl.conn = nc
	return nil
}

// Publish message
func (ncl *NATSClient) Publish(subj string, raw_msg []byte) error {
	if ncl.conn != nil {
		return ncl.conn.Publish(subj, raw_msg)
	} else {
		err := ncl.connectAndReadyToPublish()
		if err != nil {
			return err
		} else {
			return ncl.conn.Publish(subj, raw_msg)
		}
	}
}
