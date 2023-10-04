package route

import (
	"chat.com/p2p"
	"chat.com/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-session/session/v3"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
)

// upgrader : http 연결을 websocket 연결로 업그레이드하는 데 사용, 안에 내용은 버퍼 사이즈 정의
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Peers map 에서 사용할 키 값
var peerKey = 1

var publicPath = "/public"

// Start : 각 요청별 핸들러 함수 정의
func Start(port int) {
	staticHandler := http.FileServer(http.Dir("." + publicPath))
	http.Handle(publicPath+"/", http.StripPrefix(publicPath, staticHandler))

	// index
	http.HandleFunc("/", indexHandler)
	// login
	http.HandleFunc("/login", loginHandler)
	// join
	http.HandleFunc("/join", joinHandler)
	// 현재 유저 리스트 가져올 때 사용
	http.HandleFunc("/getUsers", getUsersHandler)
	// js에서 Websocket 객체 만들때 사용
	http.HandleFunc("/ws", socketHandler)

	log.Printf("Listening on localhost:%d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

// socketHandler : 앞단에서 js로 Websocket 객체를 만들 때 한번 실행 (처음 연결 후 계속 연결이 유지됨)
func socketHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fallthrough
	case "POST":
		// upgrader로 http 연결을 websocket 연결 객체로 변경
		conn, err := upgrader.Upgrade(w, r, nil)
		utils.HandleErr(err)

		// 새로운 연결 추가
		p := &p2p.Peer{
			Conn: conn,
			Key:  strconv.Itoa(peerKey),
		}
		p2p.Peers.V[strconv.Itoa(peerKey)] = p
		peerKey++

		// peer의 Read 함수를 고루틴으로 실행
		go p.Read()
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		store, err := session.Start(context.Background(), w, r)
		utils.HandleErr(err)

		_, ok := store.Get("user")
		if ok {
			http.ServeFile(w, r, "."+publicPath)
			return
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	default:
		fmt.Fprintf(w, "Sorry, only GET methods are supported.")
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		loginGetHandler(w, r)
	case "POST":
		loginPostHandler(w, r)
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func joinHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":

	case "POST":

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

// getUsersHandler : 전체 유저 리스트 요청 핸들러 함수
func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		var userNames []string

		// 전체 peer를 돌며 이름을 []string에 저장
		for _, p := range p2p.Peers.V {
			userNames = append(userNames, p.Name)
		}

		// json 으로 변환
		jsonUserNames, err := json.Marshal(userNames)
		utils.HandleErr(err)

		// http 성공 코드 및 json 데이터를 반환
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonUserNames)
	default:
		fmt.Fprintf(w, "Sorry, only GET methods are supported.")
	}
}

func loginGetHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "."+publicPath+"/login.html")
}

// loginPostHandler : 로그인 핸들러 함수
func loginPostHandler(w http.ResponseWriter, r *http.Request) {
	// post 데이터로 name을 받음
	userName := r.PostFormValue("name")

	// 이름이 존재한다면 에러 반환
	for _, p := range p2p.Peers.V {
		if p.Name == userName {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	// 이름이 존재하지 않는 peer 에 해당 이름을 매칭
	for _, p := range p2p.Peers.V {
		if p.Name == "" {
			p.Name = userName
			break
		}
	}

	// 참가 메시지 생성, Type : 0은 관리자 메시지
	var admissionChat = p2p.ChatMessage{
		Author:  "admin",
		Message: fmt.Sprintf("%s님이 참가했습니다.", userName),
		Type:    strconv.Itoa(0),
	}

	// []byte로 변경 후 전체 peer 에게 메세지 발송
	byteChat := utils.StructToBytes(admissionChat)
	p2p.SendMessageToPeers(websocket.TextMessage, byteChat)

	// http 성공 코드 반환
	w.WriteHeader(http.StatusCreated)
}

func joinGetHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "."+publicPath+"/join.html")
}

func joinPostHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
}
