package session

import (
	"errors"
	"sync"

	servermessages "github.com/Monkhai/shwipe-server.git/pkg/protocol/serverMessages"
	"github.com/Monkhai/shwipe-server.git/pkg/user"
)

type SessionUsersMap struct {
	sessionID string
	IndexMap  map[string]int
	UsersMap  map[string]*user.User
	mux       sync.RWMutex
}

func NewSessionUsersMap(sessionID string) *SessionUsersMap {
	userIndexMap := SessionUsersMap{
		sessionID: sessionID,
		UsersMap:  make(map[string]*user.User),
		IndexMap:  make(map[string]int),
		mux:       sync.RWMutex{},
	}
	for userID := range userIndexMap.IndexMap {
		userIndexMap.IndexMap[userID] = 0
	}
	return &userIndexMap
}

func (u *SessionUsersMap) AddUser(usr *user.User) error {
	if u.IsUserInMap(usr.IDToken) {
		return errors.New("user already in map")
	}

	usrs, err := u.GetAllUsers()
	if err != nil {
		return err
	}
	for _, usr := range usrs {
		usr.WriteMessage(servermessages.NewUserJoinedSessionMessage(u.sessionID, servermessages.SAFE_SessionUser{
			PhotoURL: usr.FirebaseUserRecord.PhotoURL,
			UserName: usr.FirebaseUserRecord.DisplayName,
		}))
	}

	u.mux.Lock()
	defer u.mux.Unlock()
	u.UsersMap[usr.IDToken] = usr
	u.IndexMap[usr.IDToken] = 0
	return nil
}

func (u *SessionUsersMap) GetUser(userID string) (*user.User, bool) {
	u.mux.RLock()
	defer u.mux.RUnlock()
	usr, ok := u.UsersMap[userID]
	return usr, ok
}

func (u *SessionUsersMap) GetAllUsers() ([]*user.User, error) {
	u.mux.RLock()
	defer u.mux.RUnlock()
	users := make([]*user.User, 0, len(u.UsersMap))
	for _, usr := range u.UsersMap {
		users = append(users, usr)
	}
	return users, nil
}

func (u *SessionUsersMap) IsUserInMap(userID string) bool {
	u.mux.RLock()
	defer u.mux.RUnlock()
	_, inIndexMap := u.IndexMap[userID]
	_, inUsersMap := u.UsersMap[userID]
	if !inIndexMap || !inUsersMap {
		return false
	}

	return true
}

func (u *SessionUsersMap) GetIndex(userID string) (int, error) {
	if !u.IsUserInMap(userID) {
		return 0, errors.New("user not found")
	}
	u.mux.RLock()
	defer u.mux.RUnlock()
	return u.IndexMap[userID], nil
}

func (u *SessionUsersMap) SetIndex(userID string, index int) error {
	if !u.IsUserInMap(userID) {
		return errors.New("user not found")
	}
	u.mux.Lock()
	defer u.mux.Unlock()
	u.IndexMap[userID] = index
	return nil
}

func (u *SessionUsersMap) GetUserCount() int {
	u.mux.RLock()
	defer u.mux.RUnlock()
	return len(u.UsersMap)
}

func (u *SessionUsersMap) RemoveUser(userID string) error {
	if !u.IsUserInMap(userID) {
		return errors.New("user not found")
	}
	u.mux.Lock()
	defer u.mux.Unlock()
	delete(u.UsersMap, userID)
	delete(u.IndexMap, userID)
	return nil
}
