package route

import (
	"chat.com/p2p"
	"chat.com/utils"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
)

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
