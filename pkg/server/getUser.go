package server

import (
	"encoding/json"
	"log"
	"net/http"
)

func (s *Server) GetUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	user, err := s.db.GetUser(id)
	if err != nil {
		log.Println(err, "from getUser.go")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(user)
}
