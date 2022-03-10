package sessionstore

import "time"

type Session struct {
	// id - unique
	ID interface{}

	// created time - should be invalidated when expires
	Created time.Time
}

// option type
type SessionOption func(*Session)

// option with created time
func WithTimeCreatedNow() SessionOption {
	return func(s *Session) {
		s.Created = time.Now()
	}
}

// with session id
func WithID(id interface{}) SessionOption {
	return func(s *Session) {
		s.ID = id
	}
}

// create a session
func NewSession(option ...SessionOption) *Session {
	s := &Session{}
	for _, o := range option {
		o(s)
	}
	return s
}
