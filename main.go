package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/ws", socketHandler)

	port := 8080
	log.Printf("Listening on port %d", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatal(err)
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func socketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	defer conn.Close()
	if err != nil {
		log.Printf("upgrader.Upgrade: %v", err)
		return
	}
	for {
		messageType, p, err := conn.ReadMessage()
		fmt.Println(string(p))

		if err != nil {
			log.Printf("conn.ReadMessage: %v", err)
			return
		}

		// 메세지 생성 => 메세지 발송한 부분에서만 메시지를 받을 수 있고
		// 제 3자에 경우는 메시지를 받지 못함
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Printf("conn.WriteMessage: %v", err)
			return
		}
	}

}
