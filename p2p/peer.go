package p2p

import (
	"chat.com/utils"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"strconv"
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
	Name string
	Conn *websocket.Conn
}

// Struct 내부 필드가 전역이 아니면 WriteMessage 에 넘겼을 때 사용을 못함
type ChatMessage struct {
	Author  string `json:"Author"`
	Message string `json:"Message"`
	Type    string `json:"Type"`
}

func (peer *Peer) Read() {
	defer peer.close()

	for {
		messageType, payload, err := peer.Conn.ReadMessage()
		var chat ChatMessage
		var byteChat []byte
		var isDuplication bool

		isDuplication = peerNameDuplicationCheck(peer)

		if err != nil {
			log.Printf("conn.ReadMessage: %v", err)
			return
		}

		utils.HandleErr(json.Unmarshal(payload, &chat))
		chat.Type = strconv.Itoa(1)

		// Peer 에 이름 추가
		if !isDuplication {
			peer.Name = chat.Author
		}

		byteChat = utils.StructToBytes(chat)
		SendMessageToPeers(messageType, byteChat)
	}
}

func (p *Peer) close() {
	// data race 보호를 위한 코드 추가
	Peers.m.Lock()
	defer func() {
		Peers.m.Unlock()
	}()
	p.Conn.Close()

	delete(Peers.V, p.Key)

	// 퇴장 메시지 발생
	var leaveChat ChatMessage

	leaveChat.Author = "admin"
	leaveChat.Message = fmt.Sprintf("%s님이 나갔습니다.", p.Name)
	leaveChat.Type = strconv.Itoa(0)

	byteChat := utils.StructToBytes(leaveChat)
	SendMessageToPeers(websocket.TextMessage, byteChat)
}

// peerNameDuplicationCheck : 파라미터로 온 값이 Peers 에 존재 여부 반환
func peerNameDuplicationCheck(peer *Peer) bool {
	fmt.Println("peerNameDuplicationCheck")
	if peer.Name == "" {
		for _, p := range Peers.V {
			if p.Name != "" && p.Name == peer.Name {
				return true
			}
		}
	}
	return false
}

// SendMessageToPeers : Peers 에 메세지 전달
func SendMessageToPeers(messageType int, byteChat []byte) {
	for _, p := range Peers.V {
		if err := p.Conn.WriteMessage(messageType, byteChat); err != nil {
			log.Printf("conn.WriteMessage: %v", err)
			continue
		}
	}
}
