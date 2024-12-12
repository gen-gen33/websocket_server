package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/ws", handleConnections)

	go handleMessages()

	fmt.Println("WebSocketサーバーが起動しました。ポート: 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("サーバー起動エラー: ", err)
	}
}
