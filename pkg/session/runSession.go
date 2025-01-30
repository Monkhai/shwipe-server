package session

import (
	"log"
	"sync"

	clientmessages "github.com/Monkhai/shwipe-server.git/pkg/protocol/clientMessages"
	servermessages "github.com/Monkhai/shwipe-server.git/pkg/protocol/serverMessages"
)

const FETCH_THRESHOLD = 2
const BATCH_SIZE = 20

func (s *Session) RunSession(wg *sync.WaitGroup) {
	s.mux.Lock()
	s.IsStarted = true
	s.mux.Unlock()

	defer func() {
		s.mux.Lock()
		s.IsStarted = false
		s.mux.Unlock()
		wg.Done()
	}()

	log.Println("Getting restaurants")
	restaurants, nextPageToken, err := s.restaurantAPI.GetResaturants(s.Location.Lat, s.Location.Lng, nil)
	if err != nil {
		log.Printf("Error getting restaurants: %v", err)
		return
	}
	log.Println("Got restaurants")

	safeUsers := make([]servermessages.SAFE_SessionUser, 0, len(s.UsersMap.UsersMap))
	for _, usr := range s.UsersMap.UsersMap {
		safeUsers = append(safeUsers, servermessages.SAFE_SessionUser{
			PhotoURL: usr.FirebaseUserRecord.PhotoURL,
			UserName: usr.FirebaseUserRecord.DisplayName,
		})
	}
	msg := servermessages.NewSessionStartMessage(s.ID, safeUsers, restaurants)
	for _, usr := range s.UsersMap.UsersMap {
		usr.WriteMessage(msg)
	}

	for {
		select {
		case <-s.RemoveSessionChan:
			{
				log.Println("Session context cancelled (from RemoveSessionChan)")
				return
			}
		case <-s.ctx.Done():
			{
				log.Println("Session context cancelled")
				return
			}
		case msg := <-s.msgChan:
			{
				switch msg := msg.(type) {
				case clientmessages.BaseClientMessage:
				case clientmessages.LeaveSessionMessage:
				case clientmessages.JoinSessionMessage:
				case clientmessages.CreateSessionMessage:
				case clientmessages.StartSessionMessage:
					{
						log.Printf("Unrelated message received: %v", msg)
						continue
					}
				case clientmessages.UpdateLocationMessage:
					{
						s.handleUpdateLocationMessage(msg)
					}
				case clientmessages.IndexUpdateMessage:
					{
						{
							usr, ok := s.GetUser(msg.TokenID)
							if !ok {
								log.Printf("User not found: %v", msg.TokenID)
								return
							}

							if (msg.Index+FETCH_THRESHOLD)%BATCH_SIZE != 0 {
								err := s.UsersMap.SetIndex(usr.IDToken, msg.Index)
								if err != nil {
									log.Printf("Error setting index: %v", err)
								}
								continue
							}

							left := len(restaurants) - msg.Index
							if left <= FETCH_THRESHOLD {
								newRestaurants, newNextPageToken, err := s.restaurantAPI.GetResaturants(s.Location.Lat, s.Location.Lng, nextPageToken)
								if err != nil {
									log.Printf("Error getting restaurants: %v", err)
									return
								}

								updateRestaurantsMsg := servermessages.NewRestaurantUpdateMessage(newRestaurants)
								usr.WriteMessage(updateRestaurantsMsg)

								restaurants = append(restaurants, newRestaurants...)
								nextPageToken = newNextPageToken
							} else {
								nextBatchIndex := msg.Index + FETCH_THRESHOLD
								nextBatch := restaurants[nextBatchIndex : nextBatchIndex+BATCH_SIZE]
								updateRestaurantsMsg := servermessages.NewRestaurantUpdateMessage(nextBatch)
								usr.WriteMessage(updateRestaurantsMsg)
							}
							err := s.UsersMap.SetIndex(usr.IDToken, msg.Index)
							if err != nil {
								log.Printf("Error setting index: %v", err)
							}
						}
					}

				default:
					{
						log.Printf("Unknown message type: %v", msg)
						continue
					}
				}
			}
		}
	}
}
