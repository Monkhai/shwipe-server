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

type DBFunc func(sessionID string) error

type SessionManager struct {
	Sessions       map[string]*Session
	mux            *sync.RWMutex
	ctx            context.Context
	sessionStorage SessionStorage
}

func NewSessionManager(ctx context.Context, sessionStorage SessionStorage) *SessionManager {
	return &SessionManager{
		mux:            &sync.RWMutex{},
		Sessions:       make(map[string]*Session),
		ctx:            ctx,
		sessionStorage: sessionStorage,
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
	log.Println("Creating session")
	sessionID := createSessionID()
	removeSessionChan := make(chan struct{})
	session := NewSession(sessionID, usr.Location, sm.ctx, removeSessionChan, wg)
	err := sm.addSession(session)
	if err != nil {
		log.Printf("Error adding session: %v", err)
		return nil, err
	}
	log.Println("Session created")

	err = session.AddUser(usr)
	if err != nil {
		log.Printf("Error adding user to session: %v", err)
		return nil, err
	}
	log.Println("User added to session")
	return session, nil
}

func (sm *SessionManager) addSession(session *Session) error {
	if sm.IsSessionIn(session) {
		return errors.New("session already in")
	}

	go func() {
		select {
		case <-session.ctx.Done():
			log.Println("Session context done (from addSession)")
			sm.RemoveSession(session.ID)
			return
		case <-session.RemoveSessionChan:
			log.Println("Removing session (from RemoveSessionChan)")
			sm.RemoveSession(session.ID)
			return
		}
	}()

	go func() {
		err := sm.sessionStorage.InsertSession(session.ID)
		if err != nil {
			log.Printf("Error inserting session: %v", err)
			session.RemoveSessionChan <- struct{}{}
		}
		log.Println("Session inserted into db")
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
	session, err := sm.GetSession(sessionID)
	if err != nil {
		return err
	}

	close(session.RemoveSessionChan)
	close(session.msgChan)
	session.UsersMap.CloseAllDoneChans()

	msg := servermessages.NewSessionClosedMessage()
	usrs, err := session.UsersMap.GetAllUsers()
	if err != nil {
		return err
	}
	for _, usr := range usrs {
		usr.WriteMessage(msg)
	}

	err = sm.sessionStorage.DeleteSession(sessionID)
	if err != nil {
		return err
	}
	log.Println("Session deleted from db (from RemoveSession)")

	delete(sm.Sessions, sessionID)
	log.Println("Session removed (from RemoveSession)")

	return nil
}

func (sm *SessionManager) RemoveUserFromAllSessions(usr *user.User) error {
	for _, session := range sm.Sessions {
		if session.IsUserInSession(usr.IDToken) {
			session.RemoveUser(usr)
			log.Println("User removed from session (from RemoveUserFromAllSessions)")
			if session.UsersMap.GetUserCount() == 0 {
				sm.RemoveSession(session.ID)
				log.Println("Session removed (from RemoveUserFromAllSessions)")
			}

		}
	}
	return nil
}

func (sm *SessionManager) RemoveAllSessions() error {
	for _, session := range sm.Sessions {
		go sm.RemoveSession(session.ID)
	}
	return nil
}
