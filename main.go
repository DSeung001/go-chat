package main

import (
	"chat.com/rest"
	"chat.com/route"
)

func main() {
	port := 8080

	route.Start(port)
	rest.Start(port)
}
