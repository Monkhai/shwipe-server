package server

import (
	"log"
	"sync"

	clientmessages "github.com/Monkhai/shwipe-server.git/pkg/protocol/clientMessages"
	servermessages "github.com/Monkhai/shwipe-server.git/pkg/protocol/serverMessages"
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
		case msg := <-usr.SessionMsgChan:
			{
				log.Printf("Session Message received from user server listener: %T", msg)
			}
		case msg := <-usr.ServerMsgChan:
			{
				log.Printf("Message received from user: %T", msg)
				switch m := msg.(type) {
				/*
					these are the messages that are not related
					to general server operations
				*/
				case clientmessages.UpdateLocationMessage:
				case clientmessages.IndexUpdateMessage:
					{
						continue
					}
				case clientmessages.LeaveSessionMessage:
					{
						log.Println("Leave session message received in server listener", m.SessionId)
						s, err := s.SessionManager.GetSession(m.SessionId)
						if err != nil {
							log.Printf("Session not found: %v", m.SessionId)
							continue
						}

						err = s.RemoveUser(usr)
						if err != nil {
							log.Printf("Error removing user from session: %v", err)
							continue
						}
						log.Println("user removed")
					}
				case clientmessages.CreateSessionMessage:
					{
						session, err := s.SessionManager.CreateSession(usr, s.wg)
						if err != nil {
							log.Printf("Error creating session: %v", err)
							continue
						}
						log.Println("Session created")
						usrs, err := session.UsersMap.GetAllUsers()
						if err != nil {
							log.Printf("Error getting users: %v", err)
							continue
						}
						var safeUsers []servermessages.SAFE_SessionUser
						for _, usr := range usrs {
							safeUsers = append(safeUsers, servermessages.SAFE_SessionUser{
								PhotoURL: usr.FirebaseUserRecord.PhotoURL,
								UserName: usr.FirebaseUserRecord.DisplayName,
							})
						}
						msg := servermessages.NewSessionCreatedMessage(session.ID, safeUsers)
						usr.WriteMessage(msg)
						log.Println("Session created message sent")
					}
				case clientmessages.StartSessionMessage:
					{
						log.Println("Start session message received")
						session, err := s.SessionManager.GetSession(m.SessionId)
						if err != nil {
							log.Printf("Session not found: %v", m.SessionId)
							continue
						}
						log.Println("Session found")
						s.wg.Add(1)
						go session.RunSession(s.wg)
						log.Println("Session started")
					}
				case clientmessages.JoinSessionMessage:
					{
						log.Printf("Join session message received: %v", m)
						session, err := s.SessionManager.GetSession(m.SessionId)
						if err != nil {
							log.Printf("Session not found: %v", m.SessionId)
							continue
						}
						log.Println("found session")
						err = session.AddUser(usr)
						if err != nil {
							log.Printf("Error adding user to session: %v", err)
							continue
						}
						log.Println("user added")

						usrs, err := session.UsersMap.GetAllUsers()
						if err != nil {
							log.Printf("Error getting users: %v", err)
							continue
						}
						var safeUsers []servermessages.SAFE_SessionUser
						for _, usr := range usrs {
							safeUsers = append(safeUsers, servermessages.SAFE_SessionUser{
								PhotoURL: usr.FirebaseUserRecord.PhotoURL,
								UserName: usr.FirebaseUserRecord.DisplayName,
							})
						}
						usr.WriteMessage(servermessages.NewJointSessionMessage(session.ID, session.Restaurants, safeUsers, session.IsStarted))
						log.Println("joint session message sent")
					}

				default:
					{
						log.Printf("Unhandled message type received: %T with content: %+v", m, m)
					}
				}
			}
		}
	}
}
