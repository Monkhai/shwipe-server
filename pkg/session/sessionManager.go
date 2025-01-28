package session

import (
	"errors"
	"sync"

	"github.com/Monkhai/shwipe-server.git/pkg/user"
)

type SessionManager struct {
	Sessions map[string]*Session
	mux      *sync.RWMutex
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		mux:      &sync.RWMutex{},
		Sessions: make(map[string]*Session),
	}
}

func (sm *SessionManager) GetSession(id string) (*Session, error) {
	sm.mux.RLock()
	defer sm.mux.RUnlock()
	session, ok := sm.Sessions[id]
	if !ok {
		return nil, errors.New("session not found")
	}
	return session, nil
}

func (sm *SessionManager) DeleteSession(id string) {
	sm.mux.Lock()
	defer sm.mux.Unlock()
	delete(sm.Sessions, id)
}

func (sm *SessionManager) GetAllSessions() []*Session {
	sm.mux.RLock()
	defer sm.mux.RUnlock()
	sessions := make([]*Session, 0, len(sm.Sessions))
	for _, session := range sm.Sessions {
		sessions = append(sessions, session)
	}
	return sessions
}

func (sm *SessionManager) GetSessionCount() int {
	sm.mux.RLock()
	defer sm.mux.RUnlock()
	return len(sm.Sessions)
}

func (sm *SessionManager) AddUserToSession(id string, usr *user.User) error {
	s, err := sm.GetSession(id)
	if err != nil {
		return err
	}

	err = s.AddUser(usr)
	if err != nil {
		return err
	}
	return nil
}

func (sm *SessionManager) IsSessionIn(session *Session) bool {
	sm.mux.RLock()
	defer sm.mux.RUnlock()
	_, ok := sm.Sessions[session.ID]
	return ok
}

func (sm *SessionManager) AddSession(session *Session) error {
	if sm.IsSessionIn(session) {
		return errors.New("session already in")
	}
	sm.mux.Lock()
	defer sm.mux.Unlock()
	sm.Sessions[session.ID] = session
	return nil
}
