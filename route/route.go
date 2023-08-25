package route

import (
	"chat.com/p2p"
	"chat.com/utils"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
)

type userRequestBody struct {
	name string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var peerKey = 1

func Start(port int) {
	router := mux.NewRouter()

	router.HandleFunc("/login", loginHandler).Methods("POST")
	router.HandleFunc("/ws", socketHandler).Methods("POST", "GET")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	log.Printf("Listening on localhost:%d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}

// socketHandler : 웹소캣은 일반적으로 처음에 한번 요청을 보냄
func socketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	utils.HandleErr(err)

	p := &p2p.Peer{
		Conn: conn,
		Key:  strconv.Itoa(peerKey),
	}
	p2p.Peers.V[strconv.Itoa(peerKey)] = p
	peerKey++
	go p.Read()
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var user userRequestBody
	user = userRequestBody{name: r.PostFormValue("name")}

	for _, p := range p2p.Peers.V {
		if p.Name == user.name {
			w.WriteHeader(http.StatusBadRequest)
		}
	}

	var admissionChat = p2p.ChatMessage{
		Author:  "admin",
		Message: fmt.Sprintf("%s님이 참가했습니다.", user.name),
		Type:    strconv.Itoa(0),
	}

	byteChat := utils.StructToBytes(admissionChat)
	p2p.SendMessageToPeers(websocket.TextMessage, byteChat)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(user.name))
}
