package server

import (
	uuid "github.com/satori/go.uuid"
	"log"
	"net"
	"time"
)

type SessionInfo struct {
	ReqTime       time.Time
	RespTime      time.Time
	RemoteAddress string
	Data          string
	Duration      float64
}

// Session struct
type Session struct {
	sID      string
	uID      string
	conn     *net.Conn
	settings map[string][]SessionInfo
}

// NewSession create a new session
func NewSession(conn *net.Conn) *Session {
	var (
		id  uuid.UUID
		err error
	)
	if id, err = uuid.NewV4(); err != nil {
		log.Println("failed to create uuid: ", err.Error())
		return nil
	}

	session := &Session{
		sID:      id.String(),
		uID:      "",
		conn:     conn,
		settings: make(map[string][]SessionInfo, 0),
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
	var (
		v  []SessionInfo
		ok bool
	)

	if v, ok = s.settings[key]; ok {
		return v
	}

	log.Println("failed to get value of ", s.settings[key])
	return nil
}

// SetSetting set setting
func (s *Session) SetSetting(key string, value SessionInfo) {
	s.settings[key] = append(s.settings[key], value)
}
