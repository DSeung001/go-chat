package p2p

import (
	"chat.com/utils"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{}

func Upgrade(rw http.ResponseWriter, r *http.Request) {
	openPort := r.URL.Query().Get("openPort")
	ip := utils.Splitter(r.RemoteAddr, ":", 0)

	// Port :3000 will upgrade the reqeust form :4000
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	fmt.Printf("%s wants an upgrade\n", openPort)
	conn, err := upgrader.Upgrade(rw, r, nil)
	utils.HandleErr(err)
	initPeer(conn, ip, openPort)
}
