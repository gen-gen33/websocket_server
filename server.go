package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 必要に応じてオリジンを制限
	},
}

// メッセージの構造
type Message struct {
	Type    string `json:"type"`    // メッセージの種類（例: "sync_time", "chat"）
	Data    string `json:"data"`    // メッセージの内容
	GroupID string `json:"groupId"` // グループID
}

var clients = make(map[*websocket.Conn]bool) // 接続中のクライアント
var broadcast = make(chan Message)           // メッセージチャネル

// WebSocket接続を処理
func handleConnections(w http.ResponseWriter, r *http.Request) {
	// WebSocket接続をアップグレード
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket接続エラー: %v", err)
		return
	}
	defer ws.Close()

	// クライアントを登録
	clients[ws] = true
	log.Printf("新しいクライアントが接続しました。現在のクライアント数: %d", len(clients))

	for {
		var msg Message
		// クライアントからメッセージを受信
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("メッセージ受信エラー: %v", err)
			delete(clients, ws) // クライアントを削除
			log.Printf("クライアントが切断されました。現在のクライアント数: %d", len(clients))
			break
		}
		// メッセージをブロードキャストチャネルに送信
		log.Printf("受信メッセージ: %+v", msg)
		broadcast <- msg
	}
}

// メッセージを全クライアントにブロードキャスト
func handleMessages() {
	for {
		msg := <-broadcast // チャネルからメッセージを取得
		log.Printf("ブロードキャスト中のメッセージ: %+v", msg)

		// すべての接続中のクライアントにメッセージを送信
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("メッセージ送信エラー: %v", err)
				client.Close()
				delete(clients, client) // エラーが発生したクライアントを削除
				log.Printf("エラーのためクライアントを削除しました。現在のクライアント数: %d", len(clients))
			}
		}
	}
}
