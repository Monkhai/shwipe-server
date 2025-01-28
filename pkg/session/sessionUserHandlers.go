package session

import (
	"log"

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
		for {
			select {
			case <-s.ctx.Done():
				return
			case msg, ok := <-usr.SessionMsgChan:
				{
					if !ok {
						log.Println("User context done (from session)")
						return
					}
					s.msgChan <- msg
				}
			}
		}
	}()

	return nil
}

func (s *Session) RemoveUser(userId string) error {
	//TODO: update all users with the new user list
	return s.UsersMap.RemoveUser(userId)
}

func (s *Session) IsUserInSession(userId string) bool {
	_, isInSession := s.GetUser(userId)
	return isInSession
}
