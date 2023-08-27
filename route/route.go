package route

import (
	"chat.com/p2p"
	"chat.com/utils"
	"encoding/json"
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

var peerKey = 1

func Start(port int) {
	router := mux.NewRouter()

	router.HandleFunc("/login", loginHandler).Methods("POST")
	router.HandleFunc("/getUsers", getUsersHandler).Methods("GET")
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
	userName := r.PostFormValue("name")

	// + 이름 중복 테스트하려면 피어에 이름이 있어야함
	for _, p := range p2p.Peers.V {
		if p.Name == userName {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	for _, p := range p2p.Peers.V {
		if p.Name == "" {
			p.Name = userName
			break
		}
	}

	var admissionChat = p2p.ChatMessage{
		Author:  "admin",
		Message: fmt.Sprintf("%s님이 참가했습니다.", userName),
		Type:    strconv.Itoa(0),
	}

	byteChat := utils.StructToBytes(admissionChat)
	p2p.SendMessageToPeers(websocket.TextMessage, byteChat)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(userName))
}

// 아래꺼가 빈 struct를 가진 슬라이스를 빈게옴
func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	var userNames []string

	for _, p := range p2p.Peers.V {
		userNames = append(userNames, p.Name)
	}

	jsonUserNames, err := json.Marshal(userNames)
	utils.HandleErr(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonUserNames)
}
