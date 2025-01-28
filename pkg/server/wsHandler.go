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
	userID := r.Header.Get("user_id")
	if userID == "" {
		log.Printf("user_id is required")
		return
	}
	idToken := r.Header.Get("id_token")
	if idToken == "" {
		log.Printf("id_token is required")
		return
	}

	usr := user.NewUser(userID, idToken, conn, s.ctx)
	s.UserManager.AddUser(usr)
}
