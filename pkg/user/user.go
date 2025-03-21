package user

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/Monkhai/shwipe-server.git/pkg/app"
	"github.com/Monkhai/shwipe-server.git/pkg/db"
	"github.com/Monkhai/shwipe-server.git/pkg/protocol"
	clientmessages "github.com/Monkhai/shwipe-server.git/pkg/protocol/clientMessages"
	"github.com/gorilla/websocket"
)

type User struct {
	IDToken        string
	DBUser         *db.DBUser
	Conn           *websocket.Conn
	Ctx            context.Context
	cancelCtx      context.CancelFunc
	ServerMsgChan  chan any
	SessionMsgChan chan any
	Location       protocol.Location
	authenticator  app.Authenticator
}

func NewUser(dbUser *db.DBUser, idToken string, conn *websocket.Conn, ctx context.Context, location protocol.Location, authenticator app.Authenticator) *User {
	ctx, cancel := context.WithCancel(ctx)
	return &User{
		IDToken:        idToken,
		DBUser:         dbUser,
		Conn:           conn,
		Ctx:            ctx,
		cancelCtx:      cancel,
		Location:       location,
		ServerMsgChan:  make(chan any),
		SessionMsgChan: make(chan any),
		authenticator:  authenticator,
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
				userID, err := u.authenticator.VerifyIDToken(u.IDToken)
				if err != nil {
					log.Printf("Error authenticating user: %v", err)
					continue
				}
				if userID != u.DBUser.ID {
					log.Printf("User ID mismatch: %s != %s", userID, u.DBUser.ID)
					continue
				}

				switch baseMsg.Type {
				// SESSION MESSAGES
				case clientmessages.UPDATE_INDEX_MESSAGE_TYPE:
					processedMessage, err := clientmessages.ProcessMessage[clientmessages.IndexUpdateMessage](msg)
					if err != nil {
						log.Printf("Error processing message of type %s. %q", processedMessage.Type, err)
					}
					u.SessionMsgChan <- processedMessage
				//--------------------------------
				//--------------------------------
				case clientmessages.UPDATE_LOCATION_MESSAGE_TYPE:
					proecssedMessage, err := clientmessages.ProcessMessage[clientmessages.UpdateLocationMessage](msg)
					if err != nil {
						log.Printf("Error processing message of type %s. %q", proecssedMessage.Type, err)
					}
					u.SessionMsgChan <- proecssedMessage
				//--------------------------------
				//--------------------------------

				// SERVER MESSAGES
				case clientmessages.CREATE_SESSION_MESSAGE_TYPE:
					proecssedMessage, err := clientmessages.ProcessMessage[clientmessages.CreateSessionMessage](msg)
					if err != nil {
						log.Printf("Error processing message of type %s. %q", proecssedMessage.Type, err)
					}
					u.ServerMsgChan <- proecssedMessage
				//--------------------------------
				//--------------------------------
				case clientmessages.CREATE_SESSION_WITH_FRIENDS_MESSAGE_TYPE:
					processedMessage, err := clientmessages.ProcessMessage[clientmessages.CreateSessionWithFriendsMessage](msg)
					if err != nil {
						log.Printf("Error processing message of type %s. %q", processedMessage.Type, err)
					}
					u.ServerMsgChan <- processedMessage
				//--------------------------------
				//--------------------------------
				case clientmessages.CREATE_SESSION_WITH_GROUP_MESSAGE_TYPE:
					processedmessage, err := clientmessages.ProcessMessage[clientmessages.CreateSessionWithGroupMessage](msg)
					if err != nil {
						log.Printf("Error processing message of type %s. %q", processedmessage.Type, err)
					}
					u.ServerMsgChan <- processedmessage
				//--------------------------------
				//--------------------------------
				case clientmessages.START_SESSION_MESSAGE_TYPE:
					processedMessage, err := clientmessages.ProcessMessage[clientmessages.StartSessionMessage](msg)
					if err != nil {
						log.Printf("Error processing message of type %s. %q", processedMessage.Type, err)
					}
					u.ServerMsgChan <- processedMessage
				//--------------------------------
				//--------------------------------
				case clientmessages.JOIN_SESSION_MESSAGE_TYPE:
					processedMessage, err := clientmessages.ProcessMessage[clientmessages.JoinSessionMessage](msg)
					if err != nil {
						log.Printf("Error processing message of type %s. %q", processedMessage.Type, err)
					}
					u.ServerMsgChan <- processedMessage
				//--------------------------------
				//--------------------------------
				case clientmessages.LEAVE_SESSION_MESSAGE_TYPE:
					processedMessage, err := clientmessages.ProcessMessage[clientmessages.LeaveSessionMessage](msg)
					if err != nil {
						log.Printf("Error processing message of type %s. %q", processedMessage.Type, err)
					}
					u.ServerMsgChan <- processedMessage
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

func (u *User) WriteMessage(msg any) {
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshalling message: %v", err)
		return
	}
	u.Conn.WriteMessage(websocket.TextMessage, jsonMsg)
}
