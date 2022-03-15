package sessionstore

import (
	"errors"
	"net/http"
	"time"
)

// This package is to store and provide session data for requests.
// A session should hold all data necessary data for authorization, or provide a process to gather it.

// capacity
var STORE_CAPACITY int = 10000

// session validity
var SESSION_VALIDITY time.Duration = time.Duration(15 * time.Minute)

type Store struct {
	// Channel to receive the sessions in
	// this channel is filtered - based on criteria, set at store creation
	SessionIn chan *Session

	// Channel to provide sessions
	// The provided session should always be valid
	// TODO: check how to implement filter here
	SessionOut chan *Session
}

// StoreOption type
type StoreOption func(*Store)

// Creates a new store
func NewStore(option ...StoreOption) *Store {
	s := &Store{}
	for _, o := range option {
		o(s)
	}
	return s
}

// set the capacity and create the channels
func WithInOutCapacity(cap int) StoreOption {
	return func(s *Store) {
		s.SessionIn = make(chan *Session, cap)
		s.SessionOut = make(chan *Session, cap)
	}
}

// Start consuming the sessions from the requests
func (s *Store) Start() {
	for {
		select {
		case ns := <-s.SessionIn:
			// check session validity
			if time.Since(ns.Created) < SESSION_VALIDITY {
				s.SessionOut <- ns
				// log.Println("items in out:", len(s.SessionOut))
			}
		default:
			continue
		}

		// ringbuffer-like trait, to prevent channel block
		if len(s.SessionOut) == STORE_CAPACITY {
			<-s.SessionOut
		}
	}
}

// Extracts client session from an http resonse
// the extracted session is set and pushed to the in channel of the sessionstore
func (s *Store) ExtractClientSessionFromResponse(resp *http.Response, req *http.Request, met *Meta) error {
	var cookies []*http.Cookie = resp.Cookies()
	if len(cookies) > 0 {
		s.SessionIn <- NewSession(
			// The session id represented by the cookies
			WithID(cookies),

			// the session is also timestamped for time validation
			WithTimeCreatedNow(),

			// with metadata
			WithMetaData(met),
		)
		return nil
	} else {
		return errors.New("could not extract session from response")
	}
}
