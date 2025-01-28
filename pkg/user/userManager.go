package user

import (
	"errors"
	"sync"
)

type UserManager struct {
	usersMap map[string]*User
	mux      *sync.RWMutex
}

func NewUserManager() *UserManager {
	return &UserManager{
		usersMap: make(map[string]*User),
		mux:      &sync.RWMutex{},
	}
}

func (um *UserManager) GetUser(id string) (*User, error) {
	um.mux.RLock()
	defer um.mux.RUnlock()
	user, ok := um.usersMap[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (um *UserManager) IsUserIn(user *User) bool {
	um.mux.RLock()
	defer um.mux.RUnlock()
	_, ok := um.usersMap[user.ID]
	return ok
}

func (um *UserManager) AddUser(user *User) error {
	if um.IsUserIn(user) {
		return errors.New("user already in")
	}
	um.mux.Lock()
	defer um.mux.Unlock()
	um.usersMap[user.ID] = user
	return nil
}
