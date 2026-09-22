package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRoutes(t *testing.T) {
	test := []struct {
		name       string
		path       string
		wantStatus int
	}{
		{name: "検索できる", path: "/search/users?q=octocat", wantStatus: http.StatusOK},
		{name: "検索キーワードがからなら422", path: "/search/users?q=", wantStatus: http.StatusUnprocessableEntity},
		{name: "検索キーワードがservererrorなら500", path: "/search/users?q=servererror", wantStatus: http.StatusInternalServerError},
		{name: "ユーザー詳細を取得できる", path: "/users/swiftlang", wantStatus: http.StatusOK},
		{name: "ユーザーがratelimitedなら403", path: "/users/ratelimited", wantStatus: http.StatusForbidden},
		{name: "ユーザーがnotfoundなら404", path: "/users/notfound", wantStatus: http.StatusNotFound},
		{name: "リポジトリ一覧を取得できる", path: "/users/swiftlang/repos", wantStatus: http.StatusOK},
		{name: "リポジトリ一覧でservererrorなら500", path: "/users/servererror/repos", wantStatus: http.StatusInternalServerError},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()
			routes().ServeHTTP(w, req)
			if w.Code != tt.wantStatus {
				t.Errorf("want status %d, got %d", tt.wantStatus, w.Code)
			}
		})
	}
}
