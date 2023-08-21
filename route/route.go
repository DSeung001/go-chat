package route

import (
	"chat.com/p2p"
	"chat.com/utils"
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

	//router.HandleFunc("/getUser", GetUserHandler).Methods("GET")
	router.HandleFunc("/login", LoginHandler).Methods("POST")
	router.HandleFunc("/ws", socketHandler).Methods("POST", "GET")
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
	session, _ := store.Get(r, "user")

	// 세션 저장
	var user userRequestBody
	user = userRequestBody{name: r.PostFormValue("name")}
	session.Values["user"] = user.name
	utils.HandleErr(session.Save(r, w))

	// 결과 반환
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(user.name))
}
