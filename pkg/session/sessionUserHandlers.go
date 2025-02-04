package session

import (
	"log"

	servermessages "github.com/Monkhai/shwipe-server.git/pkg/protocol/serverMessages"
	"github.com/Monkhai/shwipe-server.git/pkg/user"
)

func (s *Session) GetUser(userId string) (*user.User, bool) {
	return s.UsersMap.GetUser(userId)
}

func (s *Session) AddUser(usr *user.User) error {
	log.Println("Adding user to session")
	err := s.UsersMap.AddUser(usr)
	if err != nil {
		return err
	}

	doneChan, ok := s.UsersMap.GetDoneChan(usr.IDToken)
	if !ok {
		log.Panicln("Done chan not found for user")
		return nil
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for {
			select {
			case <-usr.Ctx.Done():
				log.Println("User context done (from session Add User)")
				s.UsersMap.CloseDoneChan(usr.IDToken)
				return
			case <-doneChan:
				log.Println("User done chan closed (from session Add User)")
				return
			case msg, ok := <-usr.SessionMsgChan:
				{
					log.Println("User message received from AddUser")
					if !ok {
						log.Println("User context done (from session)")
						return
					}

					if !s.IsUserInSession(usr.IDToken) {
						log.Println("Received message from user not in session, skipping. session.AddUser")
						s.RemoveUserSilent(usr)
						return
					}
					s.msgChan <- msg
				}
			}
		}
	}()
	s.UpdateUserList(&usr.IDToken)
	go func() {
		err := s.sessionDbOps.InsertSessionUser(s.ID, usr.DBUser.ID)
		if err != nil {
			log.Println("Error inserting session user", err)
		}
		log.Println("Session user inserted into db")
	}()

	return nil
}

func (s *Session) RemoveUserSilent(usr *user.User) error {
	err := s.UsersMap.RemoveUser(usr.IDToken)
	if err != nil {
		return err
	}

	usrCount := s.UsersMap.GetUserCount()
	if usrCount == 0 {
		log.Println("Session", s.ID, "is empty, closing")
		s.RemoveSessionChan <- struct{}{}
	}

	go func() {
		err := s.sessionDbOps.DeleteSessionUser(s.ID, usr.DBUser.ID)
		if err != nil {
			log.Println("Error deleting session user", err)
		}
		log.Println("Session user deleted from db")
	}()
	return nil
}

func (s *Session) RemoveUser(usr *user.User) error {
	log.Println("Removing user from session")
	err := s.UsersMap.RemoveUser(usr.IDToken)
	if err != nil {
		return err
	}
	s.UsersMap.CloseDoneChan(usr.IDToken)
	log.Println("User removed from session")

	msg := servermessages.NewRemovedFromSessionMessage(s.ID)
	usr.WriteMessage(msg)

	usrCount := s.UsersMap.GetUserCount()
	if usrCount == 0 {
		log.Println("Session", s.ID, "is empty, closing")
		s.RemoveSessionChan <- struct{}{}
		return nil
	}

	s.UpdateUserList(nil)
	go func() {
		err := s.sessionDbOps.DeleteSessionUser(s.ID, usr.DBUser.ID)
		if err != nil {
			log.Println("Error deleting session user", err)
		}
		log.Println("Session user deleted from db")
	}()
	return nil
}

func (s *Session) IsUserInSession(userId string) bool {
	_, isInSession := s.GetUser(userId)
	return isInSession
}

func (s *Session) UpdateUserList(usrIDToAvoid *string) error {
	usrs, err := s.UsersMap.GetAllUsers()
	if err != nil {
		return err
	}
	safeUsrs := make([]servermessages.SAFE_SessionUser, len(usrs))
	for i, usr := range usrs {
		safeUsrs[i] = servermessages.SAFE_SessionUser{
			PhotoURL:    usr.DBUser.PhotoURL,
			ID:          usr.DBUser.PublicID,
			DisplayName: usr.DBUser.DisplayName,
		}
	}
	msg := servermessages.NewUpdateUserListMessage(safeUsrs, s.ID)
	for _, usr := range usrs {
		if usrIDToAvoid != nil && usr.IDToken == *usrIDToAvoid {
			continue
		}
		usr.WriteMessage(msg)
	}
	return nil
}
