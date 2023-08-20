package route

import (
	"chat.com/p2p"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
)

func Start(port int) {
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/ws", socketHandler)

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
	log.Printf("[socketHandler]")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrader.Upgrade: %v", err)
		return
	}

	key := strconv.Itoa(len(p2p.Peers.V))

	p := &p2p.Peer{
		Conn: conn,
		Key:  key,
	}
	p2p.Peers.V[key] = p

	go p.Read()
}
