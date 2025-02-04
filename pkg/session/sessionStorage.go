package session

import "github.com/Monkhai/shwipe-server.git/pkg/db"

type SessionStorage interface {
	InsertSessionUser(sessionID string, userID string) error
	DeleteSessionUser(sessionID string, userID string) error
}

type SessionDbOps struct {
	db *db.DB
}

func NewSessionDbOps(db *db.DB) *SessionDbOps {
	return &SessionDbOps{db: db}
}

func (s *SessionDbOps) InsertSessionUser(sessionID string, userID string) error {
	return s.db.InsertSessionUser(sessionID, userID)
}

func (s *SessionDbOps) DeleteSessionUser(sessionID string, userID string) error {
	return s.db.DeleteSessionUser(sessionID, userID)
}
