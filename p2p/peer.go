package p2p

import (
	"chat.com/utils"
	"encoding/json"
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
	Key    string
	Author string
	Conn   *websocket.Conn
}

// Struct 내부 필드가 전역이 아니면 WriteMessage 에 넘겼을 때 사용을 못함
type chatMessage struct {
	Author  string `json:"Author"`
	Message string `json:"Message"`
}

func (peer *Peer) Read() {
	for {
		messageType, payload, err := peer.Conn.ReadMessage()
		if err != nil {
			log.Printf("conn.ReadMessage: %v", err)
			return
		}
		var chat chatMessage
		utils.HandleErr(json.Unmarshal([]byte(payload), &chat))

		byteChat, err := json.Marshal(chat)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// peer의 conn가 실행이 안되면 삭제 작업
		for _, p := range Peers.V {
			if err := p.Conn.WriteMessage(messageType, byteChat); err != nil {
				log.Printf("conn.WriteMessage: %v", err)
				continue
			}
		}
	}
}
