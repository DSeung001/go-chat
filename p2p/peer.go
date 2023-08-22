package p2p

import (
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
	Key  string
	Conn *websocket.Conn
}

// Struct 내부 필드가 전역이 아니면 WriteMessage 에 넘겼을 때 사용을 못함
type chatMessage struct {
	Author  string `json:"Author"`
	Content string `json:"Content"`
}

func (peer *Peer) Read(userName string) {
	for {
		messageType, payload, err := peer.Conn.ReadMessage()
		if err != nil {
			log.Printf("conn.ReadMessage: %v", err)
			return
		}

		// 페이로드에 작성자 값 추가!
		fmt.Printf("%s: %s \n", userName, string(payload))

		chat := chatMessage{
			Author:  userName,
			Content: string(payload),
		}

		byteChat, err := json.Marshal(chat)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		for _, p := range Peers.V {
			if err := p.Conn.WriteMessage(messageType, byteChat); err != nil {
				log.Printf("conn.WriteMessage: %v", err)
				continue
			}
		}
	}
}
