package p2p

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

type peers struct {
	V map[string]*Peer
	m sync.Mutex // Mutex를 넣어야 unlock/lock 가능
}

var Peers = peers{
	V: make(map[string]*Peer),
}

type Peer struct {
	Key  string
	Conn *websocket.Conn
}

func (peer *Peer) Read() {
	for {
		messageType, payload, err := peer.Conn.ReadMessage()
		fmt.Println(string(payload))

		if err != nil {
			log.Printf("conn.ReadMessage: %v", err)
			return
		}

		for _, p := range Peers.V {
			if err := p.Conn.WriteMessage(messageType, payload); err != nil {
				log.Printf("conn.WriteMessage: %v", err)
				continue
			}
		}
	}
}
