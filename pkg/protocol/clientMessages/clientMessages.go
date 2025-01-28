package clientmessages

import "github.com/Monkhai/shwipe-server.git/pkg/protocol"

const (
	INDEX_UPDATE_MESSAGE_TYPE    = "index_update"
	UPDATE_LOCATION_MESSAGE_TYPE = "update_location"

	START_SESSION_MESSAGE_TYPE  = "start_session"
	CREATE_SESSION_MESSAGE_TYPE = "create_session"
	JOIN_SESSION_MESSAGE_TYPE   = "join_session"
)

type BaseClientMessage struct {
	Type    string `json:"type"`
	TokenID string `json:"token_id"`
}

type IndexUpdateMessage struct {
	BaseClientMessage
	Index int `json:"index"`
}

type UpdateLocationMessage struct {
	BaseClientMessage
	Location protocol.Location `json:"location"`
}

type StartSessionMessage struct {
	BaseClientMessage
	SessionId string `json:"session_id"`
}

type CreateSessionMessage struct {
	BaseClientMessage
}

type JoinSessionMessage struct {
	BaseClientMessage
	SessionId string `json:"session_id"`
}
