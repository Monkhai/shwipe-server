package session

import "github.com/Monkhai/shwipe-server.git/pkg/db"

type SessionStorage interface {
	InsertSession(sessionID string) error
	DeleteSession(sessionID string) error
}

type SessionStorageManager struct {
	db *db.DB
}

func NewSessionStorageManager(db *db.DB) *SessionStorageManager {
	return &SessionStorageManager{db: db}
}

func (s *SessionStorageManager) InsertSession(sessionID string) error {
	return s.db.InsertSession(sessionID)
}

func (s *SessionStorageManager) DeleteSession(sessionID string) error {
	return s.db.DeleteSession(sessionID)
}
