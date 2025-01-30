package session

import (
	"context"
	"errors"
	"log"
	"sync"

	servermessages "github.com/Monkhai/shwipe-server.git/pkg/protocol/serverMessages"
	"github.com/Monkhai/shwipe-server.git/pkg/user"
	"github.com/google/uuid"
)

type SessionManager struct {
	Sessions map[string]*Session
	mux      *sync.RWMutex
	ctx      context.Context
}

func NewSessionManager(ctx context.Context) *SessionManager {
	return &SessionManager{
		mux:      &sync.RWMutex{},
		Sessions: make(map[string]*Session),
		ctx:      ctx,
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

func (sm *SessionManager) CreateSession(usr *user.User, wg *sync.WaitGroup) (*Session, error) {
	sessionID := createSessionID()
	_, cancel := context.WithCancel(sm.ctx)
	removeSessionChan := make(chan struct{})
	session := NewSession(sessionID, usr.Location, sm.ctx, cancel, removeSessionChan, wg)
	err := sm.addSession(session)
	if err != nil {
		log.Printf("Error adding session: %v", err)
		return nil, err
	}
	err = session.AddUser(usr)
	if err != nil {
		log.Printf("Error adding user to session: %v", err)
		return nil, err
	}
	return session, nil
}

func (sm *SessionManager) addSession(session *Session) error {
	if sm.IsSessionIn(session) {
		return errors.New("session already in")
	}

	go func() {
		select {
		case <-session.ctx.Done():
			sm.RemoveSession(session.ID)
			log.Println("Session removed")
			return
		case <-session.RemoveSessionChan:
			sm.RemoveSession(session.ID)
			log.Println("Session removed")
			return
		}
	}()

	sm.mux.Lock()
	defer sm.mux.Unlock()
	sm.Sessions[session.ID] = session
	return nil
}

func createSessionID() string {
	return uuid.New().String()
}

func (sm *SessionManager) RemoveSession(sessionID string) error {
	sm.mux.Lock()
	defer sm.mux.Unlock()
	session, ok := sm.Sessions[sessionID]
	if !ok {
		return errors.New("session not found")
	}
	close(session.RemoveSessionChan)
	close(session.msgChan)
	session.cancel()

	msg := servermessages.NewSessionClosedMessage()
	usrs, err := session.UsersMap.GetAllUsers()
	if err != nil {
		return err
	}
	for _, usr := range usrs {
		usr.WriteMessage(msg)
	}

	delete(sm.Sessions, sessionID)
	log.Println("Session removed")
	return nil
}

func (sm *SessionManager) RemoveUserFromAllSessions(usr *user.User) error {
	for _, session := range sm.Sessions {
		if session.IsUserInSession(usr.IDToken) {
			session.RemoveUser(usr)
			log.Println("User removed from session")
			if session.UsersMap.GetUserCount() == 0 {
				session.cancel()
				sm.RemoveSession(session.ID)
				log.Println("Session removed")
			}

		}
	}
	return nil
}
