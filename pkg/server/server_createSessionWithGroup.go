package server

import (
	"log"

	"github.com/Monkhai/shwipe-server.git/pkg/notifications"
	servermessages "github.com/Monkhai/shwipe-server.git/pkg/protocol/serverMessages"
	"github.com/Monkhai/shwipe-server.git/pkg/session"
	"github.com/Monkhai/shwipe-server.git/pkg/user"
)

func (s *Server) createSessionWithGroup(usr *user.User, groupId string) {
	sessionDbOps := session.NewSessionDbOps(s.DB)
	session, err := s.SessionManager.CreateSession(usr, s.wg, sessionDbOps)
	if err != nil {
		log.Printf("Error creating session: %v", err)
		return
	}
	log.Println("Session created")

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
	msg := servermessages.NewSessionCreatedMessage(session.ID, safeUsers)
	usr.WriteMessage(msg)

	pushTokens, err := s.DB.GetGroupMembersPushTokens(groupId)
	if err != nil {
		log.Printf("Error getting users push tokens: %v", err)
		return
	}

	for _, pushToken := range pushTokens {
		s.sendNotification(pushToken, session.ID, notifications.NotificationTypeSessionInvitation)
	}
}
