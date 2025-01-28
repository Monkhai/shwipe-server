package session

import (
	"errors"

	"github.com/Monkhai/shwipe-server.git/pkg/user"
)

func (s *Session) GetUser(userId string) (*user.User, bool) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	usr, ok := s.UsersMap[userId]
	if !ok {
		return &user.User{}, false
	}

	return usr, true
}

func (s *Session) AddUser(usr *user.User) error {
	if s.IsUserInSession(usr.ID) {
		return errors.New("user already in session")
	}
	s.UsersMap[usr.ID] = usr

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		select {
		case <-s.ctx.Done():
			{
				return
			}
		case msg := <-usr.MsgChan:
			{
				s.msgChan <- msg
			}
		}
	}()

	return nil
}

func (s *Session) IsUserInSession(userId string) bool {
	_, isInSession := s.GetUser(userId)
	return isInSession
}
