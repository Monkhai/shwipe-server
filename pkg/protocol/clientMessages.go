package protocol

const (
	INDEX_UPDATE_MESSAGE_TYPE    = "index_update"
	UPDATE_LOCATION_MESSAGE_TYPE = "update_location"
)

type BaseClientMessage struct {
	Type    string `json:"type"`
	ID      string `json:"id"`
	TokenID string `json:"token_id"`
}

type IndexUpdateMessage struct {
	BaseClientMessage
	Index int `json:"index"`
}

type Location struct {
	Lat string
	Lng string
}

type UpdateLocationMessage struct {
	BaseClientMessage
	Location Location `json:"location"`
}
