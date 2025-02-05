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

	restaurants, nextPageToken, err := s.restaurantAPI.GetResaturants(s.Location.Lat, s.Location.Lng, nil)
	if err != nil {
		log.Printf("Error getting restaurants: %v", err)
		return
	}

	safeUsers := make([]servermessages.SAFE_SessionUser, 0, len(s.SessionUserManager.UsersMap))
	for _, usr := range s.SessionUserManager.UsersMap {
		safeUsers = append(safeUsers, servermessages.SAFE_SessionUser{
			ID:          usr.DBUser.PublicID,
			DisplayName: usr.DBUser.DisplayName,
			PhotoURL:    usr.DBUser.PhotoURL,
		})
	}
	msg := servermessages.NewSessionStartMessage(s.ID, safeUsers, restaurants)
	for _, usr := range s.SessionUserManager.UsersMap {
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
						usr, ok := s.GetUser(msg.TokenID)
						if !ok {
							log.Printf("User not found: %v", msg.TokenID)
							return
						}

						isAllLiked := s.VoteManager.SetVote(msg.Index, msg.TokenID, msg.Liked)
						// if all users have liked the restaurant, send a match found message
						if isAllLiked {
							msg := servermessages.NewMatchFoundMessage(msg.Index)
							s.SessionUserManager.Broadcast(msg)
						}

						// session is over.
						nextIndex := msg.Index + 1

						/*
							if the index is not 2 less than a multiple of BATCH_SIZE
							then the user does not need more restaurants
							and we don't need to fetch moe restaurants
						*/
						if (nextIndex+FETCH_THRESHOLD)%BATCH_SIZE != 0 {
							err := s.SessionUserManager.SetIndex(usr.IDToken, nextIndex)
							if err != nil {
								log.Printf("Error setting index: %v", err)
							}
							continue
						}

						/*
							if the number of of restaurants left is less than or equal to the fetch threshold
							then we need to fetch more restaurants from the API
						*/
						if len(restaurants)-nextIndex <= FETCH_THRESHOLD {
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
							nextBatchIndex := nextIndex + FETCH_THRESHOLD
							nextBatch := restaurants[nextBatchIndex : nextBatchIndex+BATCH_SIZE]
							updateRestaurantsMsg := servermessages.NewRestaurantUpdateMessage(nextBatch)
							usr.WriteMessage(updateRestaurantsMsg)
						}
						err := s.SessionUserManager.SetIndex(usr.IDToken, nextIndex)
						if err != nil {
							log.Printf("Error setting index: %v", err)
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
