package p2p

import (
	"fmt"
	"github.com/gorilla/websocket"
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
	key     string
	address string
	port    string
	conn    *websocket.Conn
	inbox   chan []byte
}

// Receiver Method
func (p *peer) close() {
	// data race 보호를 위한 코드 추가

	Peers.m.Lock()
	defer func() {
		//	// 20초 안에 port:4000으로 시도시 락이 안풀려서 에러 발생!
		//	time.Sleep(20 * time.Second)
		Peers.m.Unlock()
	}()
	p.conn.Close()
	delete(Peers.v, p.key)
}

func (p *peer) read() {
	// 에러 발생 시 메서드가 종료되니 defer 사용
	defer p.close()
	for {
		m := Message{}
		err := p.conn.ReadJSON(&m)
		if err != nil {
			break
		}
		handleMsg(&m, p)
	}
}

func (p *peer) write() {
	defer p.close()
	for {
		m, ok := <-p.inbox
		if !ok {
			break
		}
		p.conn.WriteMessage(websocket.TextMessage, m)
	}
}

func initPeer(conn *websocket.Conn, address, port string) *peer {
	Peers.m.Lock()
	defer Peers.m.Unlock()
	key := fmt.Sprintf("%s:%s", address, port)
	p := &peer{
		conn:    conn,
		inbox:   make(chan []byte),
		address: address,
		key:     key,
		port:    port,
	}
	go p.read()
	go p.write()
	Peers.v[key] = p
	return p
}
