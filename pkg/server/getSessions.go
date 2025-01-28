package server

import (
	"encoding/json"
	"net/http"
)

func (s *Server) GetSessions(w http.ResponseWriter, r *http.Request) {
	sessions := s.SessionManager.GetAllSessions()
	json.NewEncoder(w).Encode(sessions)
}
