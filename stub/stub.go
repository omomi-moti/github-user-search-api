package main

import (
	_ "embed"
	"log"
	"net/http"
	// "strconv" // 2ページ目の処理を戻すときに使う
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
	q := r.URL.Query().Get("q")
	if q == "" {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}
	if status, ok := errorStatusFor(q); ok {
		http.Error(w, http.StatusText(status), status)
		return
	}
	writeJSON(w, searchUsersJSON)
}

func handleUserDetail(w http.ResponseWriter, r *http.Request) {
	if status, ok := errorStatusFor(r.PathValue("username")); ok {
		http.Error(w, http.StatusText(status), status)
		return
	}
	writeJSON(w, userDetailJSON)
}

func handleRepos(w http.ResponseWriter, r *http.Request) {
	if status, ok := errorStatusFor(r.PathValue("username")); ok {
		http.Error(w, http.StatusText(status), status)
		return
	}
	// iOSのUserDetailViewModelは、1ページが30件未満なら次のページを取りに来ない。
	// repos.jsonは1件なので2ページ目のリクエストは来ないため、処理を無効にしている。
	// repos.jsonを30件以上に増やしたときは、以下と import の "strconv" を戻す。
	// page, err := strconv.Atoi(r.URL.Query().Get("page"))
	// if err == nil && page >= 2 {
	// 	writeJSON(w, []byte("[]"))
	// 	return
	// }
	writeJSON(w, reposJSON)
}

func errorStatusFor(name string) (int, bool) {
	switch name {
	case "notfound":
		return http.StatusNotFound, true
	case "ratelimited":
		return http.StatusForbidden, true
	case "servererror":
		return http.StatusInternalServerError, true
	}
	return 0, false
}

// writeJSON はJSONをステータス200で返す
func writeJSON(w http.ResponseWriter, body []byte) {
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(body); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
