package route

import (
	"chat.com/p2p"
	"chat.com/utils"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"strconv"
)

func Start(port int) {

	router := mux.NewRouter()

	router.HandleFunc("/getUser", GetUserHandler).Methods("GET")
	router.HandleFunc("/setUser", SetUserHandler).Methods("POST")
	router.HandleFunc("/ws", socketHandler).Methods("POST")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	log.Printf("Listening on port %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

type userRequestBody struct {
	Name string
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

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetUserHandler")

	userSession, _ := store.Get(r, "user")
	json.NewEncoder(w).Encode(userSession.Values)
}

func SetUserHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "user")

	var userRequestBody userRequestBody
	json.NewDecoder(r.Body).Decode(&userRequestBody)

	session.Values["user"] = userRequestBody

	utils.HandleErr(session.Save(r, w))
	w.WriteHeader(http.StatusCreated)
}
