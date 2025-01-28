package session

import (
	"context"
	"sync"

	"github.com/Monkhai/shwipe-server.git/pkg/protocol"
	"github.com/Monkhai/shwipe-server.git/pkg/restaurant"
	"github.com/Monkhai/shwipe-server.git/pkg/user"
)

type Session struct {
	ID            string
	UsersMap      map[string]*user.User
	Location      protocol.Location
	wg            *sync.WaitGroup
	mux           *sync.RWMutex
	ctx           context.Context
	cancel        context.CancelFunc
	restaurantAPI *restaurant.RestaurantAPI
	msgChan       chan interface{}
}

func NewSession(id string, location protocol.Location, ctx context.Context, cancelCtx context.CancelFunc, wg *sync.WaitGroup) *Session {
	return &Session{
		ID:            id,
		Location:      location,
		ctx:           ctx,
		cancel:        cancelCtx,
		mux:           &sync.RWMutex{},
		UsersMap:      make(map[string]*user.User),
		restaurantAPI: restaurant.NewRestaurantAPI(restaurant.BASE_URL, restaurant.GOOGLE_PLACES_API_KEY),
		msgChan:       make(chan interface{}),
		wg:            wg,
	}
}
