package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetFavorites(t *testing.T) {
	s := newServer() //構造体を作成している(変数名がよくないかも)

	//偽のリクエストを作成する
	req := httptest.NewRequest("GET", "/favorites", nil)
	//レスポンスを描く混む入れ物を作成
	rec := httptest.NewRecorder()
	//routesになんのHandleを設定して、どこに書き込むかを指定している
	s.routes().ServeHTTP(rec, req)

	//status codeがちゃんと帰るか検証
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d , want %d", rec.Code, http.StatusOK)
	}
	//Content-Typeがちゃんと帰るか検証
	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want %q", got, "application/json")
	}

	var got []Favorite

	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if len(got) != 2 {
		t.Errorf("len(got) = %d, want %d", len(got), 2)
	}
}

func TestGetFavorite(t *testing.T) {
	tests := []struct {
		name       string
		username   string
		wantStatus int
	}{
		{name: "存在するusernameなら200", username: "omomi-moti", wantStatus: http.StatusOK},
		{name: "存在しないusernameなら404", username: "nobody", wantStatus: http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newServer()
			req := httptest.NewRequest(http.MethodGet, "/favorites/"+tt.username, nil)
			rec := httptest.NewRecorder()
			s.routes().ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestPostFavorites(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantAdded  int // POSTで増える件数
	}{
		{name: "正しいJSONなら201で追加される", body: `{"username":"swift"}`, wantStatus: http.StatusCreated, wantAdded: 1},
		{name: "壊れたJSONなら400で追加されない", body: `{"username":`, wantStatus: http.StatusBadRequest, wantAdded: 0},
		{name: "usernameが空なら400で追加されない", body: `{"username":""}`, wantStatus: http.StatusBadRequest, wantAdded: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newServer()
			before := len(s.favorites)

			req := httptest.NewRequest(http.MethodPost, "/favorites", strings.NewReader(tt.body))
			rec := httptest.NewRecorder()
			s.routes().ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d ,want %d", rec.Code, tt.wantStatus)
			}

			if after := len(s.favorites) - before; after != tt.wantAdded {
				t.Errorf("len(s.favorites) - before = %d, want %d", after, tt.wantAdded)
			}
		})
	}
}
