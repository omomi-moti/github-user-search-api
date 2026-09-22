# github-user-search-api

iOSアプリ「GitHub User Search」向けの、Goで書かれたローカル開発用サーバー。

- **お気に入りAPI**（ルート）: アプリが登録したお気に入りユーザーを、メモリ上で保持して返す
- **GitHub APIスタブ**（`stub/`）: GitHub APIと同じパス・同じ形のJSONを固定で返す。レート制限やネットワークに左右されずにアプリを動かすため

どちらも外部依存なし（標準ライブラリのみ）。データは**メモリ上のみ**で、再起動すると初期状態に戻る。

## 必要なもの

- Go 1.27.1 以上（`go.mod` を参照）

## 起動

お気に入りAPI（:8080）

```bash
go run .
```

GitHub APIスタブ（:8081）

```bash
go run ./stub
```

2つは独立しているので、必要な方だけ起動してよい。

## お気に入りAPI（:8080）

### エンドポイント

| メソッド | パス | 成功時 | 説明 |
| --- | --- | --- | --- |
| GET | `/favorites` | 200 | お気に入り一覧を返す |
| GET | `/favorites/{username}` | 200 | 1件を返す。無ければ404 |
| POST | `/favorites` | 201 | 1件追加し、追加したものを返す |
| DELETE | `/favorites/{username}` | 204 | 1件削除する。無ければ404 |

### レスポンスの形

```json
{
  "username": "octocat",
  "avatarURL": "https://avatars.githubusercontent.com/u/583231?v=4",
  "name": "The Octocat",
  "savedAt": "2026-09-13T10:00:00Z"
}
```

- `name` は未設定なら `null`（キー自体は必ず返す）
- `savedAt` はUTCで、秒より下は切り捨てる
- 一覧は0件でも `null` ではなく `[]` を返す

### エラー

| ステータス | 起きるとき |
| --- | --- |
| 400 | JSONが壊れている / `username` が空 |
| 404 | 指定した `username` が登録されていない |
| 409 | すでに同じ `username` が登録されている |

### 動作確認の例

```bash
curl -s localhost:8080/favorites
```

```bash
curl -s -X POST localhost:8080/favorites -d '{"username":"swift"}'
```

```bash
curl -s -o /dev/null -w '%{http_code}\n' -X DELETE localhost:8080/favorites/swift
```

### 起動時の初期データ

デモ用のお気に入りが2件入った状態で起動する。中身は `server.go` の `newServer` を参照。

## GitHub APIスタブ（:8081）

本物のGitHub APIと同じパスで、`stub/fixtures/` のJSONをそのまま返す。返す内容はiOS側の `APIClientTests` のテスト用JSONと揃えてある。

| メソッド | パス |
| --- | --- |
| GET | `/search/users?q={keyword}` |
| GET | `/users/{username}` |
| GET | `/users/{username}/repos` |

### エラーを起こす入れ方

`q` や `username` に次の値を渡すと、対応するエラーを返す。アプリ側のエラー表示を確認するため。

| 値 | 返るステータス |
| --- | --- |
| `notfound` | 404 |
| `ratelimited` | 403 |
| `servererror` | 500 |
| （`q` が空） | 422 |

```bash
curl -s -o /dev/null -w '%{http_code}\n' 'localhost:8081/users/ratelimited'
```

### iOSアプリからつなぐ

アプリのベースURLを `https://api.github.com` から `http://localhost:8081` に差し替える。

## テスト

```bash
go test -race -shuffle=on ./...
```

CI（GitHub Actions）では、mainへのpushとPRのたびに `go build` / `go vet` / 上記のテストを実行している。

## 構成

| ファイル | 役割 |
| --- | --- |
| `main.go` | お気に入りAPIの起動 |
| `server.go` | `server` 構造体、初期データ、ルーティング |
| `favorites.go` | `Favorite` 型と各ハンドラ |
| `favorites_test.go` | ハンドラのテスト |
| `stub/stub.go` | スタブのルーティングとハンドラ |
| `stub/fixtures/` | スタブが返すJSON |

## 制限

- データはメモリ上のみ。永続化しない
- 認証なし。ローカル開発専用で、公開して使うことは想定していない
- お気に入りの更新（PUT/PATCH）は未実装
