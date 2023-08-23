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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type userRequestBody struct {
	name string
}

func Start(port int) {

	router := mux.NewRouter()

	router.HandleFunc("/login", loginHandler).Methods("POST")
	router.HandleFunc("/ws", socketHandler).Methods("POST", "GET")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
	log.Printf("Listening on port %d", port)
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

func loginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("loginHandler")

	var user userRequestBody
	user = userRequestBody{name: r.PostFormValue("name")}

	var admissionChat = p2p.ChatMessage{
		Author:  "admin",
		Message: fmt.Sprintf("%s님이 참가했습니다.", user.name),
		Type:    strconv.Itoa(0),
	}

	byteChat := utils.StructToBytes(admissionChat)

	for _, p := range p2p.Peers.V {
		if err := p.Conn.WriteMessage(websocket.TextMessage, byteChat); err != nil {
			log.Printf("conn.WriteMessage: %v", err)
			continue
		}
	}

	// 결과 반환
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(user.name))
}
