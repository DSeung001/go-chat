package route

import (
	"chat.com/p2p"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
)

func Start(port int) {

	router := mux.NewRouter()

	router.HandleFunc("/login", LoginHandler).Methods("POST")
	router.HandleFunc("/ws", socketHandler).Methods("POST", "GET")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
	log.Printf("Listening on port %d", port)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type userRequestBody struct {
	name string
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

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LoginHandler")

	var user userRequestBody
	user = userRequestBody{name: r.PostFormValue("name")}

	// 결과 반환
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(user.name))
}
