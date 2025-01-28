package server

import (
	"log"
	"net/http"

	"github.com/Monkhai/shwipe-server.git/pkg/user"
	"github.com/gorilla/websocket"
)

func (s *Server) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error creating the ws connection: %s", err)
	}

	// get user id from the request header
	idToken := r.Header.Get("id_token")
	if idToken == "" {
		log.Printf("id_token is required")
		return
	}
	userID, err := s.app.AuthenticateUser(idToken)
	if err != nil {
		log.Printf("Error authenticating user: %v", err)
		return
	}
	userRecord, err := s.app.GetUserRecord(userID)
	if err != nil {
		log.Printf("Error getting user record: %v", err)
		return
	}
	usr := user.NewUser(userRecord, idToken, conn, s.ctx)
	s.UserManager.AddUser(usr)
	s.wg.Add(2)
	go usr.Listen(s.wg)
	go s.listenToUserMessages(usr, s.wg)
}
