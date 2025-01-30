package servermessages

import (
	"github.com/Monkhai/shwipe-server.git/pkg/restaurant"
)

const (
	ERROR_MESSAGE_TYPE                = "error"
	SESSION_START_MESSAGE_TYPE        = "session_start"
	SESSION_CREATED_MESSAGE_TYPE      = "session_create"
	JOINT_SESSION_MESSAGE_TYPE        = "joint_session"
	UPDATE_RESTAURANTS_MESSAGE_TYPE   = "update_restaurants"
	UPDATE_USER_LIST_MESSAGE_TYPE     = "update_user_list"
	SESSION_CLOSED_MESSAGE_TYPE       = "session_closed"
	REMOVED_FROM_SESSION_MESSAGE_TYPE = "removed_from_session"
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
	SessionId string             `json:"session_id"`
	Users     []SAFE_SessionUser `json:"users"`
}

func NewSessionCreatedMessage(sessionId string, users []SAFE_SessionUser) SessionCreatedMessage {
	return SessionCreatedMessage{
		BaseServerMessage: BaseServerMessage{Type: SESSION_CREATED_MESSAGE_TYPE},
		SessionId:         sessionId,
		Users:             users,
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

type RestaurantUpdateMessage struct {
	BaseServerMessage
	Restaurants []restaurant.Restaurant `json:"restaurants"`
}

func NewRestaurantUpdateMessage(restaurants []restaurant.Restaurant) RestaurantUpdateMessage {
	return RestaurantUpdateMessage{
		BaseServerMessage: BaseServerMessage{Type: UPDATE_RESTAURANTS_MESSAGE_TYPE},
		Restaurants:       restaurants,
	}
}

type JointSessionMessage struct {
	BaseServerMessage
	SessionID   string                  `json:"session_id"`
	Users       []SAFE_SessionUser      `json:"users"`
	Restaurants []restaurant.Restaurant `json:"restaurants"`
	IsStarted   bool                    `json:"is_started"`
}

func NewJointSessionMessage(sessionID string, restaurants []restaurant.Restaurant, users []SAFE_SessionUser, isStarted bool) JointSessionMessage {
	return JointSessionMessage{
		BaseServerMessage: BaseServerMessage{Type: JOINT_SESSION_MESSAGE_TYPE},
		SessionID:         sessionID,
		Users:             users,
		Restaurants:       restaurants,
		IsStarted:         isStarted,
	}
}

type UpdateUserListMessage struct {
	BaseServerMessage
	Users     []SAFE_SessionUser `json:"users"`
	SessionID string             `json:"session_id"`
}

func NewUpdateUserListMessage(users []SAFE_SessionUser, sessionID string) UpdateUserListMessage {
	return UpdateUserListMessage{
		BaseServerMessage: BaseServerMessage{Type: UPDATE_USER_LIST_MESSAGE_TYPE},
		Users:             users,
		SessionID:         sessionID,
	}
}

type SessionClosedMessage struct {
	BaseServerMessage
}

func NewSessionClosedMessage() SessionClosedMessage {
	return SessionClosedMessage{
		BaseServerMessage: BaseServerMessage{Type: SESSION_CLOSED_MESSAGE_TYPE},
	}
}

type RemovedFromSessionMessage struct {
	BaseServerMessage
}

func NewRemovedFromSessionMessage(sessionID string) RemovedFromSessionMessage {
	return RemovedFromSessionMessage{
		BaseServerMessage: BaseServerMessage{Type: REMOVED_FROM_SESSION_MESSAGE_TYPE},
	}
}
