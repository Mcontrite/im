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

//添加新的群ID到用户的groupset中
func AddGroupID(userId, gid int64) {
	//取得node
	rwlocker.Lock()
	node, ok := userNodeMap[userId]
	if ok {
		node.GroupSets.Add(gid)
	}
	//clientMap[userId] = node
	rwlocker.Unlock()
}

func JoinGroup(w http.ResponseWriter, req *http.Request) {
	var page model.Page
	utils.Bind(req, &page)
	err := service.JoinGroup(page.UserId, page.ObjectId)
	//todo 刷新用户的群组信息
	AddGroupID(page.UserId, page.ObjectId)
	if err != nil {
		utils.RespFail(w, err.Error())
	} else {
		utils.RespOK(w, nil, "")
	}
}

//加载他的群
func LoadGroups(w http.ResponseWriter, req *http.Request) {
	var page model.Page
	utils.Bind(req, &page)
	groups := service.SearchGroups(page.UserId)
	utils.RespOKList(w, groups, len(groups))
}
