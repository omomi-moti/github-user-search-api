package main

import "net/http"

// routes はGitHub APIと同じパスとハンドラの対応を返す。main とテストの両方で使う
func routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /search/users", handleSearchUsers)
	mux.HandleFunc("GET /users/{username}", handleUserDetail)
	mux.HandleFunc("GET /users/{username}/repos", handleRepos)
	return mux
}

// TODO: GitHub APIと同じ形のJSONを返す
func handleSearchUsers(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

// TODO: GitHub APIと同じ形のJSONを返す
func handleUserDetail(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

// TODO: GitHub APIと同じ形のJSONを返す
func handleRepos(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
