package user

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/Monkhai/shwipe-server.git/pkg/db"
	"github.com/Monkhai/shwipe-server.git/pkg/protocol"
	clientmessages "github.com/Monkhai/shwipe-server.git/pkg/protocol/clientMessages"
	"github.com/gorilla/websocket"
)

type User struct {
	IDToken             string
	DBUser              *db.DBUser
	Conn                *websocket.Conn
	Ctx                 context.Context
	cancelCtx           context.CancelFunc
	ServerMsgChan       chan interface{}
	SessionMsgChan      chan interface{}
	Location            protocol.Location
	AuthenticateMessage func(token string) (string, error)
}

func NewUser(dbUser *db.DBUser, idToken string, conn *websocket.Conn, ctx context.Context, location protocol.Location, authenticateMessage func(token string) (string, error)) *User {
	ctx, cancel := context.WithCancel(ctx)
	return &User{
		IDToken:             idToken,
		DBUser:              dbUser,
		Conn:                conn,
		Ctx:                 ctx,
		cancelCtx:           cancel,
		Location:            location,
		ServerMsgChan:       make(chan interface{}),
		SessionMsgChan:      make(chan interface{}),
		AuthenticateMessage: authenticateMessage,
	}
}

func (u *User) Listen(wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
		close(u.ServerMsgChan)
		close(u.SessionMsgChan)
		log.Println("user.Listen function finished")
	}()

	messageChan := make(chan []byte)
	errorChan := make(chan error)

	go func() {
		for {
			_, msg, err := u.Conn.ReadMessage()
			if err != nil {
				errorChan <- err
				return
			}
			messageChan <- msg
		}
	}()

	for {
		select {
		case <-u.Ctx.Done():
			{
				log.Println("User context done (from user)")
				return
			}
		case msg := <-messageChan:
			{
				var baseMsg clientmessages.BaseClientMessage
				if err := json.Unmarshal(msg, &baseMsg); err != nil {
					log.Printf("Error unmarshalling message: %v", err)
					continue
				}

				log.Printf("user received message: %v", baseMsg.Type)
				userID, err := u.AuthenticateMessage(u.IDToken)
				if err != nil {
					log.Printf("Error authenticating user: %v", err)
					continue
				}
				if userID != u.DBUser.ID {
					log.Printf("User ID mismatch: %s != %s", userID, u.DBUser.ID)
					continue
				}

				switch baseMsg.Type {
				case clientmessages.UPDATE_INDEX_MESSAGE_TYPE:
					{
						var indexUpdateMessage clientmessages.IndexUpdateMessage
						if err := json.Unmarshal(msg, &indexUpdateMessage); err != nil {
							log.Printf("Error unmarshalling index update message: %v", err)
							continue
						}
						u.SessionMsgChan <- indexUpdateMessage
					}
				case clientmessages.CREATE_SESSION_MESSAGE_TYPE:
					{
						var createSessionMessage clientmessages.CreateSessionMessage
						if err := json.Unmarshal(msg, &createSessionMessage); err != nil {
							log.Printf("Error unmarshalling create session message: %v", err)
							continue
						}
						u.ServerMsgChan <- createSessionMessage
					}
				case clientmessages.CREATE_SESSION_WITH_FRIENDS_MESSAGE_TYPE:
					{
						var createSessionWithFriendsMessage clientmessages.CreateSessionWithFriendsMessage
						if err := json.Unmarshal(msg, &createSessionWithFriendsMessage); err != nil {
							log.Printf("Error unmarshalling create session with friends message: %v", err)
							continue
						}
						u.ServerMsgChan <- createSessionWithFriendsMessage
					}
				case clientmessages.START_SESSION_MESSAGE_TYPE:
					{
						var startSessionMessage clientmessages.StartSessionMessage
						if err := json.Unmarshal(msg, &startSessionMessage); err != nil {
							log.Printf("Error unmarshalling create session message: %v", err)
							continue
						}
						u.ServerMsgChan <- startSessionMessage
					}
				case clientmessages.UPDATE_LOCATION_MESSAGE_TYPE:
					{
						var updateLocationMessage clientmessages.UpdateLocationMessage
						if err := json.Unmarshal(msg, &updateLocationMessage); err != nil {
							log.Printf("Error unmarshalling update location message: %v", err)
							continue
						}
						u.SessionMsgChan <- updateLocationMessage
					}
				case clientmessages.JOIN_SESSION_MESSAGE_TYPE:
					{
						var joinSessionMessage clientmessages.JoinSessionMessage
						if err := json.Unmarshal(msg, &joinSessionMessage); err != nil {
							log.Printf("Error unmarshalling join session message: %v", err)
							continue
						}
						u.ServerMsgChan <- joinSessionMessage
					}
				case clientmessages.LEAVE_SESSION_MESSAGE_TYPE:
					{
						var leaveSessionMessage clientmessages.LeaveSessionMessage
						if err := json.Unmarshal(msg, &leaveSessionMessage); err != nil {
							log.Printf("Error unmarshalling leave session message: %v", err)
							continue
						}
						u.ServerMsgChan <- leaveSessionMessage
					}
				}
			}
		case err := <-errorChan:
			{
				if err != nil {
					if websocket.IsCloseError(
						err,
						websocket.CloseNormalClosure,
						websocket.CloseGoingAway,
						websocket.CloseAbnormalClosure,
						websocket.CloseNoStatusReceived) {
						log.Printf("Player %s disconnected gracefully\n", u.DBUser.ID)
					} else {
						log.Printf("Unexpected error reading from player %s: %v\n", u.DBUser.ID, err)
					}
				}
				u.cancelCtx()
				return
			}

		}
	}
}

func (u *User) WriteMessage(msg interface{}) {
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshalling message: %v", err)
		return
	}
	u.Conn.WriteMessage(websocket.TextMessage, jsonMsg)
}
