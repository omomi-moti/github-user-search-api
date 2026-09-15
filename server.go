package main

import (
	"net/http"
	"sync"
	"time"
)

type server struct {
	mu        sync.Mutex
	favorites []Favorite
}

func newServer() *server {
	return &server{
		favorites: []Favorite{
			{
				Username:  "omomi-moti",
				AvatarURL: "https://avatars.githubusercontent.com/u/1?v=4",
				Name:      "鈴木聖也",
				SavedAt:   time.Date(2026, 9, 13, 10, 0, 0, 0, time.UTC),
			},
			{
				Username:  "onevcat",
				AvatarURL: "https://avatars.githubusercontent.com/u/2?v=4",
				Name:      "",
				SavedAt:   time.Date(2026, 9, 13, 11, 30, 0, 0, time.UTC),
			},
		},
	}
}

func (s *server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /favorites", s.handleGetFavorites)
	mux.HandleFunc("POST /favorites", s.handlePostFavorites)
	mux.HandleFunc("GET /favorites/{username}", s.handleGetFavorite)
	return mux
}
