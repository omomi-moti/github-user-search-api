package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
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
		{name: "すでに存在するなら409で追加されない", body: `{"username":"omomi-moti"}`, wantStatus: http.StatusConflict, wantAdded: 0},
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

func TestGetFavoriteName(t *testing.T) {
	tests := []struct {
		name     string
		username string
		wantName string // レスポンスのJSONに含まれているはずの文字列
	}{
		{name: "名前が設定されていればその名前を返す", username: "omomi-moti", wantName: `"name":"鈴木聖也"`},
		{name: "名前が未設定ならnullを返す", username: "onevcat", wantName: `"name":null`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newServer()
			req := httptest.NewRequest(http.MethodGet, "/favorites/"+tt.username, nil)
			rec := httptest.NewRecorder()
			s.routes().ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
			}

			// デコードすると null とキーなしの区別がつかないため、JSONの文字列のまま確認する
			if body := rec.Body.String(); !strings.Contains(body, tt.wantName) {
				t.Errorf("body = %s, want to contain %s", body, tt.wantName)
			}
		})
	}
}

func TestPostFavoritesSavedAt(t *testing.T) {
	s := newServer()
	req := httptest.NewRequest(http.MethodPost, "/favorites", strings.NewReader(`{"username":"swift"}`))
	rec := httptest.NewRecorder()
	s.routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusCreated)
	}

	var got Favorite
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if got.Username != "swift" {
		t.Errorf("got.Username = %s, want %s", got.Username, "swift")
	}

	added := s.favorites[len(s.favorites)-1]
	if got := added.SavedAt.Nanosecond(); got != 0 {
		t.Errorf("SavedAt.Nanosecond() = %d, want 0", got)
	}
	if got := added.SavedAt.Location(); got != time.UTC {
		t.Errorf("SavedAt.Location() = %v, want UTC", got)
	}
}

func TestDeleteFavorite(t *testing.T) {
	tests := []struct {
		name        string
		username    string
		wantStatus  int
		wantRemoved int // 減る件数
	}{
		{name: "登録済みなら204で削除される", username: "omomi-moti", wantStatus: http.StatusNoContent, wantRemoved: 1},
		{name: "登録されていなければ404で何も消えない", username: "nonexistent", wantStatus: http.StatusNotFound, wantRemoved: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newServer()
			before := len(s.favorites)

			req := httptest.NewRequest(http.MethodDelete, "/favorites/"+tt.username, nil)
			rec := httptest.NewRecorder()
			s.routes().ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if removed := before - len(s.favorites); removed != tt.wantRemoved {
				t.Errorf("removed = %d, want %d", removed, tt.wantRemoved)
			}
		})
	}
}

func TestDeleteAllFavoritesReturnsEmptyArray(t *testing.T) {
	s := newServer()
	handler := s.routes()

	// デモデータの2件を両方削除する
	for _, username := range []string{"omomi-moti", "onevcat"} {
		req := httptest.NewRequest(http.MethodDelete, "/favorites/"+username, nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusNoContent {
			t.Fatalf("DELETE %s: status = %d, want %d", username, rec.Code, http.StatusNoContent)
		}
	}

	// 0件になった一覧を取得する
	req := httptest.NewRequest(http.MethodGet, "/favorites", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	// iOSの[ServerFavorite]はnullを読めないので、[]で返ることを確かめる
	if body := strings.TrimSpace(rec.Body.String()); body != "[]" {
		t.Errorf("body = %s, want []", body)
	}
}
