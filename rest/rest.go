package rest

import (
	"chat.com/p2p"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var port string

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL)
		next.ServeHTTP(rw, r)
	})
}

func Start(portNumber int) {
	port = fmt.Sprintf(":%d", portNumber)

	// Gorilla 의 Router 기능 사용
	router := mux.NewRouter()

	// 미들웨어
	router.Use(jsonContentTypeMiddleware, loggerMiddleware)

	router.HandleFunc("/ws", p2p.Upgrade).Methods("GET")

	fmt.Printf("[REST} Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
