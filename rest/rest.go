package rest

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func Start(portNumber int) {
	port := fmt.Sprintf(":%d", portNumber)

	// Gorilla 의 Router 기능 사용
	router := mux.NewRouter()

	// 미들웨어
	router.Use(jsonContentTypeMiddleware)

	router.HandleFunc("/login", loginHandler).Methods("GET")

	fmt.Printf("[REST} Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}

// jsonContentTypeMiddleware : Content-Type을 Json으로 바꾸는 미들 웨어
func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func loginHandler(rw http.ResponseWriter, r *http.Request) {
	// 로그인 후 결과값 반환

	// 해야할 것
	// 세션에 고유값 저장 후 그 고유값으로 peer map 에서 connection 찾게 수정
	/*
	1. 세션 생성 => 키 저장
	2. 키를 기준으로 peer 생성
	3. sockethandler에서 키 기준으로 사용

	*/
}
