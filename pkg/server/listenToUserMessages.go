package server

import (
	"log"
	"sync"

	clientmessages "github.com/Monkhai/shwipe-server.git/pkg/protocol/clientMessages"
	servermessages "github.com/Monkhai/shwipe-server.git/pkg/protocol/serverMessages"
	"github.com/Monkhai/shwipe-server.git/pkg/session"
	"github.com/Monkhai/shwipe-server.git/pkg/user"
)

func (s *Server) listenToUserMessages(usr *user.User, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-s.ctx.Done():
			{
				log.Println("Server context done")
				return
			}
		case <-usr.Ctx.Done():
			{
				log.Printf("User context done (from server)")
				s.UserManager.RemoveUser(usr.IDToken)
				s.SessionManager.RemoveUserFromAllSessions(usr)
				return
			}
		case msg := <-usr.ServerMsgChan:
			{
				switch m := msg.(type) {
				/*
					these are the messages that are not related
					to general server operations
				*/
				case clientmessages.UpdateLocationMessage:
				case clientmessages.IndexUpdateMessage:
					continue
				//--------------------------------
				case clientmessages.LeaveSessionMessage:
					s.leaveSession(m, usr)
				//--------------------------------
				case clientmessages.CreateSessionMessage:
					s.createSession(usr)
				//--------------------------------
				case clientmessages.CreateSessionWithFriendsMessage:
					s.createSessionWithUser(usr, m.FriendIds)
				//--------------------------------
				case clientmessages.CreateSessionWithGroupMessage:
					s.createSessionWithGroup(usr, m.GroupId)
				//--------------------------------
				case clientmessages.StartSessionMessage:
					s.startSession(m)
				//--------------------------------
				case clientmessages.JoinSessionMessage:
					s.joinSession(m, usr)
				//--------------------------------
				default:
					log.Printf("Unhandled message type received: %T with content: %+v", m, m)
				}
			}
		}
	}
}

func (s *Server) createSession(usr *user.User) {
	sessionDbOps := session.NewSessionDbOps(s.DB)
	session, err := s.SessionManager.CreateSession(usr, s.wg, sessionDbOps)
	if err != nil {
		log.Printf("Error creating session: %v", err)
		return
	}

	usrs, err := session.SessionUserManager.GetAllUsers()
	if err != nil {
		log.Printf("Error getting users: %v", err)
		return
	}
	var safeUsers []servermessages.SAFE_SessionUser
	for _, usr := range usrs {
		safeUsers = append(safeUsers, servermessages.SAFE_SessionUser{
			ID:          usr.DBUser.PublicID,
			DisplayName: usr.DBUser.DisplayName,
			PhotoURL:    usr.DBUser.PhotoURL,
		})
	}
	newSessionCreatedMessage := servermessages.NewSessionCreatedMessage(session.ID, safeUsers)
	usr.WriteMessage(newSessionCreatedMessage)
}

func (s *Server) leaveSession(m clientmessages.LeaveSessionMessage, usr *user.User) {
	log.Println("Leave session message received in server listener", m.SessionId)
	session, err := s.SessionManager.GetSession(m.SessionId)
	if err != nil {
		log.Printf("Session not found: %v", m.SessionId)
		return
	}

	err = session.RemoveUser(usr)
	if err != nil {
		log.Printf("Error removing user from session: %v", err)
		return
	}
	log.Println("user removed")
}

func (s *Server) startSession(m clientmessages.StartSessionMessage) {
	log.Println("Start session message received")
	session, err := s.SessionManager.GetSession(m.SessionId)
	if err != nil {
		log.Printf("Session not found: %v", m.SessionId)
		return
	}
	log.Println("Session found")
	s.wg.Add(1)
	go session.RunSession(s.wg)
	log.Println("Session started")
}

func (s *Server) joinSession(m clientmessages.JoinSessionMessage, usr *user.User) {
	session, err := s.SessionManager.GetSession(m.SessionId)
	if err != nil {
		log.Printf("Session not found: %v", m.SessionId)
		return
	}

	err = session.AddUser(usr)
	if err != nil {
		log.Printf("Error adding user to session: %v", err)
		return
	}

	usrs, err := session.SessionUserManager.GetAllUsers()
	if err != nil {
		log.Printf("Error getting users: %v", err)
		return
	}
	var safeUsers []servermessages.SAFE_SessionUser
	for _, usr := range usrs {
		safeUsers = append(safeUsers, servermessages.SAFE_SessionUser{
			ID:          usr.IDToken,
			DisplayName: usr.DBUser.DisplayName,
			PhotoURL:    usr.DBUser.PhotoURL,
		})
	}
	usr.WriteMessage(servermessages.NewJointSessionMessage(session.ID, session.Restaurants, safeUsers, session.IsStarted))
	log.Println("joint session message sent")
}
