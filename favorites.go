package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Favorite struct {
	Username  string    `json:"username"`
	AvatarURL string    `json:"avatarURL"`
	Name      *string   `json:"name"`
	SavedAt   time.Time `json:"savedAt"`
}

func (s *server) handleGetFavorites(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(s.favorites); err != nil {
		log.Printf("Error encoding favorites: %v", err)
	}
}

func (s *server) handleGetFavorite(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	username := r.PathValue("username")

	for _, f := range s.favorites {
		if f.Username == username {
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(f); err != nil {
				log.Printf("Error encoding favorites: %v", err)
			}
			return
		}
	}
	http.Error(w, "favorite not found", http.StatusNotFound)
}

func (s *server) handlePostFavorites(w http.ResponseWriter, r *http.Request) {

	var f Favorite
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if f.Username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, existing := range s.favorites {
		if existing.Username == f.Username {
			http.Error(w, "favorite already exists", http.StatusConflict)
			return
		}
	}

	f.SavedAt = time.Now().UTC().Truncate(time.Second)
	s.favorites = append(s.favorites, f)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(f); err != nil {
		log.Printf("Error encoding favorites: %v", err)
	}
}
