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
	route.Start(port)
}

func setSessionKey() {
	var sessionKey = securecookie.GenerateRandomKey(32)
	utils.HandleErr(os.Setenv("SESSION_KEY", string(sessionKey)))
}
