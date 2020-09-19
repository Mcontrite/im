package controller

import (
	"im/model"
	"im/service"
	"im/utils"
	"net/http"
)

func CreateGroup(w http.ResponseWriter, req *http.Request) {
	var group model.Group
	utils.Bind(req, &group)
	newgroup, err := service.CreateGroup(group)
	if err != nil {
		utils.RespFail(w, err.Error())
	} else {
		utils.RespOK(w, newgroup, "")
	}
}

func JoinGroup(w http.ResponseWriter, req *http.Request) {
	var page model.Page
	utils.Bind(req, &page)
	err := service.JoinGroup(page.UserID, page.ObjectID)
	//todo 刷新用户的群组信息
	AddGroupID(page.UserID, page.ObjectID)
	if err != nil {
		utils.RespFail(w, err.Error())
	} else {
		utils.RespOK(w, nil, "")
	}
}

//加载他的群
func LoadGroup(w http.ResponseWriter, req *http.Request) {
	var page model.Page
	utils.Bind(req, &page)
	groups := service.SearchGroup(page.UserID)
	utils.RespOkList(w, groups, len(groups))
}
