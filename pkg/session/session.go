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
	cancel            context.CancelFunc
	restaurantAPI     *restaurant.RestaurantAPI
	msgChan           chan interface{}
	RemoveSessionChan chan struct{}
	IsStarted         bool
}

func NewSession(
	id string,
	location protocol.Location,
	ctx context.Context,
	cancelCtx context.CancelFunc,
	closeSessionChan chan struct{},
	wg *sync.WaitGroup,
) *Session {
	return &Session{
		ID:                id,
		Location:          location,
		Restaurants:       []restaurant.Restaurant{},
		ctx:               ctx,
		cancel:            cancelCtx,
		mux:               &sync.RWMutex{},
		UsersMap:          NewSessionUsersMap(id),
		restaurantAPI:     restaurant.NewRestaurantAPI(secrets.BASE_URL, secrets.GOOGLE_PLACES_API_KEY),
		msgChan:           make(chan interface{}),
		wg:                wg,
		RemoveSessionChan: closeSessionChan,
		IsStarted:         false,
	}
}
