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
	MsgChan            chan interface{}
	Location           protocol.Location
}

func NewUser(userRecord *auth.UserRecord, idToken string, conn *websocket.Conn, ctx context.Context, location protocol.Location) *User {
	ctx, cancel := context.WithCancel(ctx)
	return &User{FirebaseUserRecord: userRecord, IDToken: idToken, Conn: conn, Ctx: ctx, cancelCtx: cancel, Location: location}
}

func (u *User) Listen(wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
		close(u.MsgChan)
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
				u.cancelCtx()
				return
			}
		case msg := <-messageChan:
			{
				var baseMsg clientmessages.BaseClientMessage
				if err := json.Unmarshal(msg, &baseMsg); err != nil {
					log.Printf("Error unmarshalling message: %v", err)
					continue
				}

				switch baseMsg.Type {
				case clientmessages.INDEX_UPDATE_MESSAGE_TYPE:
					{
						var indexUpdateMessage clientmessages.IndexUpdateMessage
						if err := json.Unmarshal(msg, &indexUpdateMessage); err != nil {
							log.Printf("Error unmarshalling index update message: %v", err)
							continue
						}
						u.MsgChan <- indexUpdateMessage
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
