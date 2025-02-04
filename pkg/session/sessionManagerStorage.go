package session

import "github.com/Monkhai/shwipe-server.git/pkg/db"

type SessionManagerStorage interface {
	InsertSession(sessionID string) error
	DeleteSession(sessionID string) error
}

type SessionManagerDbOps struct {
	db *db.DB
}

func NewSessionMangerDbOps(db *db.DB) *SessionManagerDbOps {
	return &SessionManagerDbOps{db: db}
}

func (s *SessionManagerDbOps) InsertSession(sessionID string) error {
	return s.db.InsertSession(sessionID)
}

func (s *SessionManagerDbOps) DeleteSession(sessionID string) error {
	return s.db.DeleteSession(sessionID)
}
