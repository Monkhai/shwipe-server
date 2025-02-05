package clientmessages

import "github.com/Monkhai/shwipe-server.git/pkg/protocol"

const (
	UPDATE_INDEX_MESSAGE_TYPE                = "update_index"
	UPDATE_LOCATION_MESSAGE_TYPE             = "update_location"
	START_SESSION_MESSAGE_TYPE               = "start_session"
	CREATE_SESSION_MESSAGE_TYPE              = "create_session"
	CREATE_SESSION_WITH_FRIENDS_MESSAGE_TYPE = "create_session_with_friends"
	JOIN_SESSION_MESSAGE_TYPE                = "join_session"
	LEAVE_SESSION_MESSAGE_TYPE               = "leave_session"
)

type BaseClientMessage struct {
	Type    string `json:"type"`
	TokenID string `json:"token_id"`
}

type IndexUpdateMessage struct {
	BaseClientMessage
	Index     int    `json:"index"`
	SessionId string `json:"session_id"`
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

type LeaveSessionMessage struct {
	BaseClientMessage
	SessionId string `json:"session_id"`
}

type CreateSessionWithFriendsMessage struct {
	BaseClientMessage
	FriendIds []string `json:"friend_ids"`
}
