package session

import (
	clientmessages "github.com/Monkhai/shwipe-server.git/pkg/protocol/clientMessages"
)

func (s *Session) handleUpdateLocationMessage(msg clientmessages.UpdateLocationMessage) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.Location = msg.Location
}
