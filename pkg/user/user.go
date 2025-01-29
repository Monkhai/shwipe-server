package user

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"firebase.google.com/go/auth"
	"github.com/Monkhai/shwipe-server.git/pkg/protocol"
	clientmessages "github.com/Monkhai/shwipe-server.git/pkg/protocol/clientMessages"
	"github.com/gorilla/websocket"
)

type User struct {
	IDToken            string
	Conn               *websocket.Conn
	FirebaseUserRecord *auth.UserRecord
	Ctx                context.Context
	cancelCtx          context.CancelFunc
	ServerMsgChan      chan interface{}
	SessionMsgChan     chan interface{}
	Location           protocol.Location
}

func NewUser(userRecord *auth.UserRecord, idToken string, conn *websocket.Conn, ctx context.Context, location protocol.Location) *User {
	ctx, cancel := context.WithCancel(ctx)
	return &User{
		FirebaseUserRecord: userRecord,
		IDToken:            idToken,
		Conn:               conn,
		Ctx:                ctx,
		cancelCtx:          cancel,
		Location:           location,
		ServerMsgChan:      make(chan interface{}),
		SessionMsgChan:     make(chan interface{}),
	}
}

func (u *User) Listen(wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
		close(u.ServerMsgChan)
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
		case msg := <-messageChan:
			{
				var baseMsg clientmessages.BaseClientMessage
				if err := json.Unmarshal(msg, &baseMsg); err != nil {
					log.Printf("Error unmarshalling message: %v", err)
					continue
				}

				log.Printf("Received message: %v", baseMsg.Type)

				switch baseMsg.Type {
				case clientmessages.UPDATE_INDEX_MESSAGE_TYPE:
					{
						var indexUpdateMessage clientmessages.IndexUpdateMessage
						if err := json.Unmarshal(msg, &indexUpdateMessage); err != nil {
							log.Printf("Error unmarshalling index update message: %v", err)
							continue
						}
						u.SessionMsgChan <- indexUpdateMessage
						log.Println("Index update message sent")
					}
				case clientmessages.CREATE_SESSION_MESSAGE_TYPE:
					{
						var createSessionMessage clientmessages.CreateSessionMessage
						if err := json.Unmarshal(msg, &createSessionMessage); err != nil {
							log.Printf("Error unmarshalling create session message: %v", err)
							continue
						}
						u.ServerMsgChan <- createSessionMessage
						log.Println("Create session message sent")
					}
				case clientmessages.START_SESSION_MESSAGE_TYPE:
					{
						var startSessionMessage clientmessages.StartSessionMessage
						if err := json.Unmarshal(msg, &startSessionMessage); err != nil {
							log.Printf("Error unmarshalling create session message: %v", err)
							continue
						}
						u.ServerMsgChan <- startSessionMessage
						log.Println("Start session message sent")
					}
				case clientmessages.UPDATE_LOCATION_MESSAGE_TYPE:
					{
						var updateLocationMessage clientmessages.UpdateLocationMessage
						if err := json.Unmarshal(msg, &updateLocationMessage); err != nil {
							log.Printf("Error unmarshalling update location message: %v", err)
							continue
						}
						u.SessionMsgChan <- updateLocationMessage
						log.Println("Update location message sent")
					}
				case clientmessages.JOIN_SESSION_MESSAGE_TYPE:
					{
						var joinSessionMessage clientmessages.JoinSessionMessage
						if err := json.Unmarshal(msg, &joinSessionMessage); err != nil {
							log.Printf("Error unmarshalling join session message: %v", err)
							continue
						}
						u.ServerMsgChan <- joinSessionMessage
						log.Println("Join session message sent")
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
						log.Printf("Player %s disconnected gracefully\n", u.FirebaseUserRecord.UID)
					} else {
						log.Printf("Unexpected error reading from player %s: %v\n", u.FirebaseUserRecord.UID, err)
					}
				}
				u.cancelCtx()
				return
			}
		case <-u.Ctx.Done():
			{
				log.Println("User context done (from user)")
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
