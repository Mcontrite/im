package controller

import (
	"im/model"
	"im/service"
	"im/utils"
	"net/http"
)

//添加好友
func AddFriend(w http.ResponseWriter, req *http.Request) {
	// request.ParseForm()
	// mobile := request.PostForm.Get("mobile")
	// passwd := request.PostForm.Get("passwd")
	var page model.Page
	utils.Bind(req, &page)
	err := service.AddFriend(page.UserID, page.ObjectID)
	if err != nil {
		utils.RespFail(w, err.Error())
	} else {
		utils.RespOK(w, nil, "好友添加成功")
	}
}

//加载个人的好友
func LoadFriend(w http.ResponseWriter, req *http.Request) {
	var page model.Page
	utils.Bind(req, &page)
	users := service.SearchFriend(page.UserID)
	utils.RespOkList(w, users, len(users))
}
