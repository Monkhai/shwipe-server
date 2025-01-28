package servermessages

import (
	"github.com/Monkhai/shwipe-server.git/pkg/restaurant"
)

const (
	ERROR_MESSAGE_TYPE             = "error"
	RESTAURANT_LIST_MESSAGE_TYPE   = "restaurant_list"
	RESTAURANT_UPDATE_MESSAGE_TYPE = "restaurant_update"
	SESSION_START_MESSAGE_TYPE     = "session_start"
	SESSION_CREATED_MESSAGE_TYPE   = "session_create"
)

type SAFE_SessionUser struct {
	PhotoURL string `json:"photo_url"`
	UserName string `json:"user_name"`
}

type BaseServerMessage struct {
	Type string `json:"type"`
}

type SessionStartMessage struct {
	BaseServerMessage
	SessionId   string                  `json:"session_id"`
	Users       []SAFE_SessionUser      `json:"users"`
	Restaurants []restaurant.Restaurant `json:"restaurants"`
}

func NewSessionStartMessage(sessionId string, users []SAFE_SessionUser, restaurants []restaurant.Restaurant) SessionStartMessage {
	return SessionStartMessage{
		BaseServerMessage: BaseServerMessage{Type: SESSION_START_MESSAGE_TYPE},
		SessionId:         sessionId,
		Users:             users,
		Restaurants:       restaurants,
	}
}

type SessionCreatedMessage struct {
	BaseServerMessage
	SessionId string `json:"session_id"`
}

func NewSessionCreatedMessage(sessionId string) SessionCreatedMessage {
	return SessionCreatedMessage{
		BaseServerMessage: BaseServerMessage{Type: SESSION_CREATED_MESSAGE_TYPE},
		SessionId:         sessionId,
	}
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

type RestaurantUpdateMessage struct {
	BaseServerMessage
	Restaurants []restaurant.Restaurant `json:"restaurants"`
}

func NewRestaurantUpdateMessage(restaurants []restaurant.Restaurant) RestaurantUpdateMessage {
	return RestaurantUpdateMessage{
		BaseServerMessage: BaseServerMessage{Type: RESTAURANT_UPDATE_MESSAGE_TYPE},
		Restaurants:       restaurants,
	}
}
