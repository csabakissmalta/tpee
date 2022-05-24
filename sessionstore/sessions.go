package sessionstore

import (
	"time"

	"github.com/csabakissmalta/tpee/exec"
)

type Session struct {
	// id - unique
	ID interface{}

	// created time - should be invalidated when expires
	Created time.Time

	// metadata - any kind of data, which can be attached to the session and passed along with it
	Meta *Meta
}

// meta type - which should instruct and hold any additional session data
type Meta struct {
	// data
	Data map[string]interface{}

	// instruction
	Instruction string
}

// option type
type SessionOption func(*Session)

// option with created time
func WithTimeCreatedNow() SessionOption {
	return func(s *Session) {
		s.Created = time.Now()
	}
}

// with metadata
func WithMetaData(m *Meta) SessionOption {
	return func(s *Session) {
		s.Meta = m
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

// store.Store interface impl for the session
func (sess *Session) SaveData(extracted interface{}, rule *exec.ExecRequestsElemDataPersistenceDataOutElem) {
	meta := sess.GetSessionMeta()
	if meta.Data == nil {
		meta.Data = make(map[string]interface{})
	}
	meta.Data[rule.Name] = extracted
}

// store.Store interface impl for the session
func (sess *Session) RetrieveData(name string) interface{} {
	return sess.Meta.Data[name]
}
