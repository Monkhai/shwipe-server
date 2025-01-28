package user

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/Monkhai/shwipe-server.git/pkg/protocol"
	"github.com/gorilla/websocket"
)

type User struct {
	ID        string
	IDToken   string
	Conn      *websocket.Conn
	ctx       context.Context
	cancelCtx context.CancelFunc
	MsgChan   chan interface{}
}

func NewUser(id string, idToken string, conn *websocket.Conn, ctx context.Context) *User {
	ctx, cancel := context.WithCancel(ctx)
	return &User{ID: id, IDToken: idToken, Conn: conn, ctx: ctx, cancelCtx: cancel}
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
		case <-u.ctx.Done():
			{
				u.cancelCtx()
				return
			}
		case msg := <-messageChan:
			{
				var baseMsg protocol.BaseClientMessage
				if err := json.Unmarshal(msg, &baseMsg); err != nil {
					log.Printf("Error unmarshalling message: %v", err)
					continue
				}

				switch baseMsg.Type {
				case protocol.INDEX_UPDATE_MESSAGE_TYPE:
					{
						var indexUpdateMessage protocol.IndexUpdateMessage
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
						log.Printf("Player %s disconnected gracefully\n", u.ID)
					} else {
						log.Printf("Unexpected error reading from player %s: %v\n", u.ID, err)
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
