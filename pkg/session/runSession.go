package session

import (
	"log"

	"github.com/Monkhai/shwipe-server.git/pkg/protocol"
)

func (s *Session) RunSession() {
	usersIndexMap := make(map[string]int)
	restaurants, nextPageToken, err := s.restaurantAPI.GetResaturants(s.Location.Lat, s.Location.Lng, nil)
	if err != nil {
		log.Printf("Error getting restaurants: %v", err)
		return
	}

	for _, usr := range s.UsersMap {
		msg := protocol.NewRestaurantListMessage(restaurants)
		usr.WriteMessage(protocol.IndexUpdateMessage{
			Index: 0,
		})
	}
}
