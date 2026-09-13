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
	Name      string    `json:"name"`
	SavedAt   time.Time `json:"savedAt"`
}

var favorites = []Favorite{ //でもデータを作成
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
}

func handleGetFavorites(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(favorites); err != nil {
		log.Printf("Error encoding favorites: %v", err)
	}
}

func handleGetFavorite(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")

	for _, f := range favorites {
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

func handlePostFavorites(w http.ResponseWriter, r *http.Request) {
	var f Favorite
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if f.Username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}

	f.SavedAt = time.Now()
	favorites = append(favorites, f)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(favorites); err != nil {
		log.Printf("Error encoding favorites: %v", err)
	}

}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /favorites", handleGetFavorites)
	mux.HandleFunc("POST /favorites", handlePostFavorites)
	mux.HandleFunc("GET /favorites/{username}", handleGetFavorite)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
