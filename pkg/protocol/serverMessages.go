package protocol

import "github.com/Monkhai/shwipe-server.git/pkg/restaurant"

const (
	ERROR_MESSAGE_TYPE           = "error"
	RESTAURANT_LIST_MESSAGE_TYPE = "restaurant_list"
)

type BaseServerMessage struct {
	Type string `json:"type"`
}

type ErrorMessage struct {
	BaseServerMessage
	Error string `json:"error"`
}

func NewErrorMessage(error string) ErrorMessage {
	return ErrorMessage{
		BaseServerMessage: BaseServerMessage{Type: ERROR_MESSAGE_TYPE},
		Error:             error,
	}
}

//==================================

type RestaurantListMessage struct {
	BaseServerMessage
	Restaurants []restaurant.Restaurant `json:"restaurants"`
}

func NewRestaurantListMessage(restaurants []restaurant.Restaurant) RestaurantListMessage {
	return RestaurantListMessage{
		BaseServerMessage: BaseServerMessage{Type: RESTAURANT_LIST_MESSAGE_TYPE},
		Restaurants:       restaurants,
	}
}
