package server

import (
	uuid "github.com/satori/go.uuid"
	"net"
)

// Session struct
type Session struct {
	sID      string
	uID      string
	conn     *net.Conn
	settings map[string]interface{}
}

// NewSession create a new session
func NewSession(conn *net.Conn) *Session {
	id, _ := uuid.NewV4()
	session := &Session{
		sID:      id.String(),
		uID:      "",
		conn:     conn,
		settings: make(map[string]interface{}),
	}

	return session
}

// GetSessionID get session ID
func (s *Session) GetSessionID() string {
	return s.sID
}

// BindUserID bind a user ID to session
func (s *Session) BindUserID(uid string) {
	s.uID = uid
}

// GetUserID get user ID
func (s *Session) GetUserID() string {
	return s.uID
}

// GetConn get zero.Conn pointer
func (s *Session) GetConn() *net.Conn {
	return s.conn
}

// SetConn set a zero.Conn to session
func (s *Session) SetConn(conn *net.Conn) {
	s.conn = conn
}

// GetSetting get setting
func (s *Session) GetSetting(key string) interface{} {

	if v, ok := s.settings[key]; ok {
		return v
	}

	return nil
}

// SetSetting set setting
func (s *Session) SetSetting(key string, value interface{}) {
	s.settings[key] = value
}
