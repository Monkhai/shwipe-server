package session

import (
	"errors"
	"log"
	"sync"

	"github.com/Monkhai/shwipe-server.git/pkg/user"
)

type SessionUserManager struct {
	sessionID string
	IndexMap  map[string]int
	UsersMap  map[string]*user.User
	doneChans map[string]chan struct{}
	mux       sync.RWMutex
}

func NewSessionUsersMap(sessionID string) *SessionUserManager {
	userIndexMap := SessionUserManager{
		sessionID: sessionID,
		UsersMap:  make(map[string]*user.User),
		IndexMap:  make(map[string]int),
		doneChans: make(map[string]chan struct{}),
		mux:       sync.RWMutex{},
	}
	for userID := range userIndexMap.IndexMap {
		userIndexMap.IndexMap[userID] = 0
	}
	return &userIndexMap
}

func (sum *SessionUserManager) AddUser(usr *user.User) error {
	if sum.IsUserInMap(usr.IDToken) {
		return errors.New("user already in map")
	}

	sum.mux.Lock()
	sum.UsersMap[usr.IDToken] = usr
	sum.IndexMap[usr.IDToken] = 0
	sum.mux.Unlock()
	sum.SetDoneChan(usr.IDToken, make(chan struct{}))

	return nil
}

func (sum *SessionUserManager) RemoveUser(userID string) error {
	if !sum.IsUserInMap(userID) {
		return errors.New("user not found")
	}
	sum.mux.Lock()
	defer sum.mux.Unlock()
	delete(sum.UsersMap, userID)
	delete(sum.IndexMap, userID)
	return nil
}

func (sum *SessionUserManager) GetUser(userID string) (*user.User, bool) {
	sum.mux.RLock()
	defer sum.mux.RUnlock()
	usr, ok := sum.UsersMap[userID]
	return usr, ok
}

func (sum *SessionUserManager) GetAllUsers() ([]*user.User, error) {
	sum.mux.RLock()
	defer sum.mux.RUnlock()
	users := make([]*user.User, 0, len(sum.UsersMap))
	for _, usr := range sum.UsersMap {
		users = append(users, usr)
	}
	return users, nil
}

func (sum *SessionUserManager) IsUserInMap(userID string) bool {
	sum.mux.RLock()
	defer sum.mux.RUnlock()
	_, inIndexMap := sum.IndexMap[userID]
	_, inUsersMap := sum.UsersMap[userID]
	if !inIndexMap || !inUsersMap {
		return false
	}
	return true
}

func (sum *SessionUserManager) GetIndex(userID string) (int, error) {
	if !sum.IsUserInMap(userID) {
		return 0, errors.New("user not found")
	}
	sum.mux.RLock()
	defer sum.mux.RUnlock()
	return sum.IndexMap[userID], nil
}

func (sum *SessionUserManager) SetIndex(userID string, index int) error {
	if !sum.IsUserInMap(userID) {
		return errors.New("user not found")
	}
	sum.mux.Lock()
	defer sum.mux.Unlock()
	sum.IndexMap[userID] = index
	return nil
}

func (sum *SessionUserManager) GetUserCount() int {
	sum.mux.RLock()
	defer sum.mux.RUnlock()
	return len(sum.UsersMap)
}

func (sum *SessionUserManager) SetDoneChan(userID string, doneChan chan struct{}) {
	log.Println("Setting done chan for user")
	sum.mux.Lock()
	defer sum.mux.Unlock()
	sum.doneChans[userID] = doneChan
	log.Println("Set done chan for user")
}

func (sum *SessionUserManager) GetDoneChan(userID string) (chan struct{}, bool) {
	sum.mux.RLock()
	defer sum.mux.RUnlock()
	doneChan, ok := sum.doneChans[userID]
	return doneChan, ok
}

func (sum *SessionUserManager) CloseDoneChan(userID string) {
	sum.mux.Lock()
	defer sum.mux.Unlock()
	doneChan, ok := sum.doneChans[userID]
	if ok {
		close(doneChan)
		delete(sum.doneChans, userID)
	}
}

func (sum *SessionUserManager) CloseAllDoneChans() {
	sum.mux.Lock()
	defer sum.mux.Unlock()
	for _, doneChan := range sum.doneChans {
		close(doneChan)
	}
	sum.doneChans = make(map[string]chan struct{})
}

func (sum *SessionUserManager) Broadcast(msg interface{}) {
	sum.mux.RLock()
	defer sum.mux.RUnlock()
	for _, user := range sum.UsersMap {
		user.WriteMessage(msg)
	}
}
