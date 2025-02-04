package session

import (
	"context"
	"sync"

	"github.com/Monkhai/shwipe-server.git/pkg/protocol"
	"github.com/Monkhai/shwipe-server.git/pkg/restaurant"
	"github.com/Monkhai/shwipe-server.git/secrets"
)

type Session struct {
	ID                string
	UsersMap          *SessionUsersMap
	Location          protocol.Location
	Restaurants       []restaurant.Restaurant
	wg                *sync.WaitGroup
	mux               *sync.RWMutex
	ctx               context.Context
	restaurantAPI     *restaurant.RestaurantAPI
	sessionDbOps      *SessionDbOps
	msgChan           chan interface{}
	RemoveSessionChan chan struct{}
	CloseOnce         sync.Once
	IsStarted         bool
}

func NewSession(
	id string,
	location protocol.Location,
	ctx context.Context,
	closeSessionChan chan struct{},
	wg *sync.WaitGroup,
	sessionDbOps *SessionDbOps,
) *Session {
	return &Session{
		ID:                id,
		Location:          location,
		Restaurants:       []restaurant.Restaurant{},
		ctx:               ctx,
		mux:               &sync.RWMutex{},
		UsersMap:          NewSessionUsersMap(id),
		restaurantAPI:     restaurant.NewRestaurantAPI(secrets.BASE_URL, secrets.GOOGLE_PLACES_API_KEY),
		msgChan:           make(chan interface{}),
		wg:                wg,
		RemoveSessionChan: closeSessionChan,
		CloseOnce:         sync.Once{},
		IsStarted:         false,
		sessionDbOps:      sessionDbOps,
	}
}
