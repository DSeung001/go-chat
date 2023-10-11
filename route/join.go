package route

import (
	"chat.com/db"
	"chat.com/model"
	"chat.com/utils"
	"fmt"
	"net/http"
)

func joinGetHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "."+publicPath+"/join.html")
}

func joinPostHandler(w http.ResponseWriter, r *http.Request) {
	utils.HandleErr(r.ParseForm())
	postData := r.PostForm

	user := model.User{
		Name:            postData.Get("name"),
		Email:           postData.Get("email"),
		Password:        postData.Get("password"),
		ConfirmPassword: postData.Get("confirm-password"),
	}

	// 유효성 검사
	if user.Password != user.ConfirmPassword {
		fmt.Println("실패")
		http.Error(w, "비밀번호 불일치", http.StatusBadRequest)
		return
	}

	fmt.Println(user)
	// go에서는 그냥 instert? db? 파일에서? 여기서?
	db.Create("INSERT INTO user(name, email, password) VALUES (?,?,?)", user.Name, user.Email, user.Password)
	fmt.Println("성공")
	//w.WriteHeader(http.StatusCreated)
	http.Redirect(w, r, "/login", http.StatusSeeOther)

}
