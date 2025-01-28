package session

import (
	"github.com/Monkhai/shwipe-server.git/pkg/user"
)

func (s *Session) GetUser(userId string) (*user.User, bool) {
	return s.UsersMap.GetUser(userId)
}

func (s *Session) AddUser(usr *user.User) error {
	err := s.UsersMap.AddUser(usr)
	if err != nil {
		return err
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		select {
		case <-s.ctx.Done():
			return
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
