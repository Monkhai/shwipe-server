package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Monkhai/shwipe-server.git/pkg/protocol"
	servermessages "github.com/Monkhai/shwipe-server.git/pkg/protocol/serverMessages"
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
	idToken := r.URL.Query().Get("token_id")
	if idToken == "" {
		log.Printf("id_token is required. User not authenticated and not allowed to connect")
		conn.WriteMessage(websocket.CloseMessage, []byte(""))
		conn.Close()
		return
	}

	lat, lng := r.URL.Query().Get("lat"), r.URL.Query().Get("lng")
	if lat == "" || lng == "" {
		log.Printf("lat and lng are required. User not authenticated and not allowed to connect")
		conn.WriteMessage(websocket.CloseMessage, []byte(""))
		conn.Close()
		return
	}
	location := protocol.Location{Lat: lat, Lng: lng}

	userID, err := s.app.AuthenticateUser(idToken)
	if err != nil {
		log.Printf("Error authenticating user: %v", err)
		conn.WriteMessage(websocket.CloseMessage, []byte(""))
		conn.Close()
		return
	}

	loadingConnectionMsg := servermessages.NewLoadingConnectionMessage()
	msgBytes, err := json.Marshal(loadingConnectionMsg)
	if err != nil {
		log.Printf("Error marshalling message: %v", err)
		conn.WriteMessage(websocket.CloseMessage, []byte(""))
		conn.Close()
		return
	}
	conn.WriteMessage(websocket.TextMessage, msgBytes)

	dbUser, err := s.DB.GetUser(userID)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		conn.WriteMessage(websocket.CloseMessage, []byte(""))
		conn.Close()
		return
	}
	log.Println("Got user from db")

	usr := user.NewUser(dbUser, idToken, conn, s.ctx, location, s.app.AuthenticateUser)
	s.UserManager.AddUser(usr)

	connectionEstablised := servermessages.NewConnectionEstablishedMessage()
	msgBytes, err = json.Marshal(connectionEstablised)
	if err != nil {
		log.Printf("Error marshalling message: %v", err)
		conn.WriteMessage(websocket.CloseMessage, []byte(""))
		conn.Close()
		return
	}
	conn.WriteMessage(websocket.TextMessage, msgBytes)

	s.wg.Add(2)
	go usr.Listen(s.wg)
	go s.listenToUserMessages(usr, s.wg)
}
