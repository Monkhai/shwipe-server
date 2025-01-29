package server

import (
	"encoding/json"
	"net/http"
)

func (s *Server) GetSessions(w http.ResponseWriter, r *http.Request) {
	tokenID := r.Header.Get("Authorization")
	if tokenID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	_, err := s.app.AuthenticateUser(tokenID)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	sessions := s.SessionManager.GetAllSessions()
	sessionList := make([]string, 0, len(sessions))
	for _, session := range sessions {
		sessionList = append(sessionList, session.ID)
	}
	json.NewEncoder(w).Encode(sessionList)
}
