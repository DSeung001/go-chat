package main

import (
	"chat.com/route"
	"chat.com/utils"
	"github.com/gorilla/securecookie"
	"os"
)

func main() {
	port := 8080

	setSessionKey()

	// router가 두개 면 안된다
	// 아래 코드의 라우터 병합이 필요
	route.Start(port)
}

func setSessionKey() {
	var sessionKey = securecookie.GenerateRandomKey(32)
	utils.HandleErr(os.Setenv("SESSION_KEY", string(sessionKey)))
}
