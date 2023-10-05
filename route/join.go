package route

import (
	"net/http"
)

func joinGetHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "."+publicPath+"/join.html")
}

func joinPostHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
}
