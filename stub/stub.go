package main

import (
	_ "embed"
	"log"
	"net/http"
)

// 返すJSONは、iOSのAPIClientTestsのテスト用JSONと同じ内容にしている

//go:embed fixtures/search_users.json
var searchUsersJSON []byte

//go:embed fixtures/user_detail.json
var userDetailJSON []byte

//go:embed fixtures/repos.json
var reposJSON []byte

// routes はGitHub APIと同じパスとハンドラの対応を返す。main とテストの両方で使う
func routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /search/users", handleSearchUsers)
	mux.HandleFunc("GET /users/{username}", handleUserDetail)
	mux.HandleFunc("GET /users/{username}/repos", handleRepos)
	return mux
}

func handleSearchUsers(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, searchUsersJSON)
}

func handleUserDetail(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, userDetailJSON)
}

func handleRepos(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, reposJSON)
}

// writeJSON はJSONをステータス200で返す
func writeJSON(w http.ResponseWriter, body []byte) {
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(body); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
