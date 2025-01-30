package session

import (
	"errors"
	"sync"

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

func (sum *SessionUsersMap) AddUser(usr *user.User) error {
	if sum.IsUserInMap(usr.IDToken) {
		return errors.New("user already in map")
	}

	sum.mux.Lock()
	defer sum.mux.Unlock()
	sum.UsersMap[usr.IDToken] = usr
	sum.IndexMap[usr.IDToken] = 0

	return nil
}

func (sum *SessionUsersMap) RemoveUser(userID string) error {
	if !sum.IsUserInMap(userID) {
		return errors.New("user not found")
	}
	sum.mux.Lock()
	defer sum.mux.Unlock()
	delete(sum.UsersMap, userID)
	delete(sum.IndexMap, userID)
	return nil
}

func (sum *SessionUsersMap) GetUser(userID string) (*user.User, bool) {
	sum.mux.RLock()
	defer sum.mux.RUnlock()
	usr, ok := sum.UsersMap[userID]
	return usr, ok
}

func (sum *SessionUsersMap) GetAllUsers() ([]*user.User, error) {
	sum.mux.RLock()
	defer sum.mux.RUnlock()
	users := make([]*user.User, 0, len(sum.UsersMap))
	for _, usr := range sum.UsersMap {
		users = append(users, usr)
	}
	return users, nil
}

func (sum *SessionUsersMap) IsUserInMap(userID string) bool {
	sum.mux.RLock()
	defer sum.mux.RUnlock()
	_, inIndexMap := sum.IndexMap[userID]
	_, inUsersMap := sum.UsersMap[userID]
	if !inIndexMap || !inUsersMap {
		return false
	}
	return true
}

func (sum *SessionUsersMap) GetIndex(userID string) (int, error) {
	if !sum.IsUserInMap(userID) {
		return 0, errors.New("user not found")
	}
	sum.mux.RLock()
	defer sum.mux.RUnlock()
	return sum.IndexMap[userID], nil
}

func (sum *SessionUsersMap) SetIndex(userID string, index int) error {
	if !sum.IsUserInMap(userID) {
		return errors.New("user not found")
	}
	sum.mux.Lock()
	defer sum.mux.Unlock()
	sum.IndexMap[userID] = index
	return nil
}

func (sum *SessionUsersMap) GetUserCount() int {
	sum.mux.RLock()
	defer sum.mux.RUnlock()
	return len(sum.UsersMap)
}
