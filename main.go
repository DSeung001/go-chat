package main

import (
	"chat.com/db"
	"chat.com/route"
)

func main() {
	db.Start()
	port := 8080
	route.Start(port)
}
