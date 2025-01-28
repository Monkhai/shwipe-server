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
		case <-usr.Ctx.Done():
			{
				log.Printf("User context done")
				return
			}
		case msg := <-usr.MsgChan:
			{
				switch m := msg.(type) {
				/*
					these are the messages that are not related to sessions
					they are handled by the server and not by the session
				*/
				case clientmessages.UpdateLocationMessage:
				case clientmessages.IndexUpdateMessage:
					{
						continue
					}
				case clientmessages.CreateSessionMessage:
					{
						session := s.SessionManager.CreateSession(usr, s.wg, usr.Ctx)
						msg := servermessages.NewSessionCreatedMessage(session.ID)
						usr.WriteMessage(msg)
					}
				case clientmessages.StartSessionMessage:
					{
						session, err := s.SessionManager.GetSession(m.SessionId)
						if err != nil {
							log.Printf("Session not found: %v", m.SessionId)
							continue
						}
						s.wg.Add(1)
						go session.RunSession(s.wg)
					}
				case clientmessages.JoinSessionMessage:
					{
						log.Printf("Join session message received: %v", m)
						session, err := s.SessionManager.GetSession(m.SessionId)
						if err != nil {
							log.Printf("Session not found: %v", m.SessionId)
							continue
						}
						err = session.AddUser(usr)
						if err != nil {
							log.Printf("Error adding user to session: %v", err)
							continue
						}
					}
				default:
					{
						log.Printf("Unknown message type: %v", m)
					}
				}
			}
		}
	}
}
