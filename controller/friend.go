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
	// phone := request.PostForm.Get("phone")
	// password := request.PostForm.Get("password")
	var page model.Page
	utils.Bind(req, &page)
	err := service.AddFriend(page.UserId, page.ObjectId)
	if err != nil {
		utils.RespFail(w, err.Error())
	} else {
		utils.RespOK(w, nil, "好友添加成功")
	}
}

//加载全部好友
func LoadFriends(w http.ResponseWriter, req *http.Request) {
	var page model.Page
	utils.Bind(req, &page)
	users := service.SearchFriends(page.UserId)
	utils.RespOKList(w, users, len(users))
}
