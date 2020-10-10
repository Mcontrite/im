package controller

import (
	"im/model"
	"im/service"
	"im/utils"
	"net/http"
)

func UserRegister(w http.ResponseWriter, req *http.Request) {
	//req.ParseForm()
	//phone := req.PostForm.Get("phone")
	//plainpwd := req.PostForm.Get("password")
	//username := fmt.Sprintf("user%06d",rand.Int31())
	//avatar :=""
	//sex := model.SEX_UNKNOW
	//有了数据绑定方法,不需要其他的啦
	var user model.User
	utils.Bind(req, &user)
	user, err := service.Register(
		user.Phone,
		user.Password,
		user.Username,
		user.Avatar,
		user.Sex)
	if err != nil {
		utils.RespFail(w, err.Error())
	} else {
		utils.RespOK(w, user, "")
	}
}

func UserLogin(w http.ResponseWriter, req *http.Request) {
	//restapi json/xml返回
	req.ParseForm()
	phone := req.PostForm.Get("phone")
	password := req.PostForm.Get("password")
	user, err := service.Login(phone, password)
	if err != nil {
		utils.RespFail(w, err.Error())
	} else {
		utils.RespOK(w, user, "")
	}
}

func FindUserByID(w http.ResponseWriter, req *http.Request) {
	var user model.User
	utils.Bind(req, &user)
	user = service.GetUserByID(user.Id)
	if user.Id == 0 {
		utils.RespFail(w, "该用户不存在")
	} else {
		utils.RespOK(w, user, "")
	}
}
