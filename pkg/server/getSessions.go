package server

import (
	"encoding/json"
	"net/http"
)

func (s *Server) GetSessions(w http.ResponseWriter, r *http.Request) {
	// tokenId := r.Header.Get("Authorization")

	sessions := s.SessionManager.GetAllSessions()
	json.NewEncoder(w).Encode(sessions)
}
