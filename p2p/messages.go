package p2p

import (
	"fmt"
)

type MessageKind int

const (
	// iota 를 사용하는 순간 1,2,3 로 차례대로 매핑
	MessageNewestBlock MessageKind = iota
	MessageAllBlocksRequest
	MessageAllBlocksResponse
	MessageNewBlockNotify
	MessageNewTxNotify
	MessageNewPeerNotify
)

type Message struct {
	Kind    MessageKind
	Payload []byte
}

func handleMsg(m *Message, p *peer) {
	switch m.Kind {

	default:
		fmt.Printf("Peer: %s, Sent a message with kind of: %d", p.key, m.Kind)
	}
}
