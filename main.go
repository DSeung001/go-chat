package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type peers struct {
	v map[string]*peer
	m sync.Mutex // Mutex를 넣어야 unlock/lock 가능
}

var Peers peers = peers{
	v: make(map[string]*peer),
}

type peer struct {
	key  string
	conn *websocket.Conn
}

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
	if err != nil {
		log.Printf("upgrader.Upgrade: %v", err)
		return
	}

	key := strconv.Itoa(len(Peers.v))

	p := &peer{
		conn: conn,
		key:  key,
	}
	Peers.v[key] = p

	go p.read()
}

func (peer *peer) read() {
	for {
		messageType, payload, err := peer.conn.ReadMessage()
		fmt.Println(string(payload))

		if err != nil {
			log.Printf("conn.ReadMessage: %v", err)
			return
		}

		// 이벤트?
		for _, p := range Peers.v {
			if err := p.conn.WriteMessage(messageType, payload); err != nil {
				log.Printf("conn.WriteMessage: %v", err)
				return
			}

		}
	}
}
