package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/Monkhai/shwipe-server.git/pkg/db"
	"github.com/Monkhai/shwipe-server.git/pkg/protocol"
	servermessages "github.com/Monkhai/shwipe-server.git/pkg/protocol/serverMessages"
	"github.com/Monkhai/shwipe-server.git/pkg/user"
	"github.com/gorilla/websocket"
)

func (s *Server) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgradeConnection(w, r)
	if err != nil {
		log.Printf("error upgrading connection: %s", err)
		return
	}

	idToken, location, err := processParameters(r, conn)
	if err != nil {
		log.Printf("error processing parameters: %s", err)
		conn.Close()
		return
	}

	userID, err := s.app.AuthenticateUser(idToken)
	if err != nil {
		log.Printf("Error authenticating user: %v", err)
		conn.WriteMessage(websocket.CloseMessage, []byte(""))
		conn.Close()
		return
	}

	dbUser, err := s.getDbUser(userID)
	if err != nil {
		log.Printf("Error getting user: %v", err)
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

func upgradeConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	upgrader := websocket.Upgrader{}
	return upgrader.Upgrade(w, r, nil)
}

func processParameters(r *http.Request, conn *websocket.Conn) (string, protocol.Location, error) {
	idToken := r.URL.Query().Get("token_id")
	if idToken == "" {
		return "", protocol.Location{}, errors.New("id_token is required. User not authenticated and not allowed to connect")
	}

	lat, lng := r.URL.Query().Get("lat"), r.URL.Query().Get("lng")
	if lat == "" || lng == "" {
		log.Printf("lat and lng are required. User not authenticated and not allowed to connect")
		conn.WriteMessage(websocket.CloseMessage, []byte(""))
		conn.Close()
		return "", protocol.Location{}, errors.New("lat and lng are required. User not authenticated and not allowed to connect")
	}
	location := protocol.Location{Lat: lat, Lng: lng}
	return idToken, location, nil
}

func (s *Server) getDbUser(userID string) (*db.DBUser, error) {
	var dbUser *db.DBUser
	var err error

	if cachedUser, exists := s.UserCache.Get(userID); exists {
		dbUser = cachedUser
	} else {
		dbUser, err = s.DB.GetUser(userID)
		if err != nil {
			return nil, err
		}
		s.UserCache.Add(userID, dbUser)
	}

	return dbUser, nil
}
