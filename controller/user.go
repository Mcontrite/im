package controller

import (
	"im/model"
	"im/service"
	"im/utils"
	"net/http"
)

func UserRegister(writer http.ResponseWriter, request *http.Request) {
	//request.ParseForm()
	//phone := request.PostForm.Get("phone")
	//plainpwd := request.PostForm.Get("password")
	//nickname := fmt.Sprintf("user%06d",rand.Int31())
	//avatar :=""
	//sex := model.SEX_UNKNOW
	//有了数据绑定方法,不需要其他的啦
	var user model.User
	utils.Bind(request, &user)
	user, err := service.Register(
		user.Phone,
		user.Password,
		user.Username,
		user.Avatar,
		user.Sex)
	if err != nil {
		utils.RespFail(writer, err.Error())
	} else {
		utils.RespOK(writer, user, "")
	}
}

func UserLogin(writer http.ResponseWriter, request *http.Request) {
	//restapi json/xml返回
	request.ParseForm()
	phone := request.PostForm.Get("phone")
	password := request.PostForm.Get("password")
	user, err := service.Login(phone, password)
	if err != nil {
		utils.RespFail(writer, err.Error())
	} else {
		utils.RespOK(writer, user, "")
	}
}

//解析一下
func FindUserById(writer http.ResponseWriter, request *http.Request) {
	var user model.User
	utils.Bind(request, &user)
	user = service.GetUserByID(user.ID)
	if user.ID == 0 {
		utils.RespFail(writer, "该用户不存在")
	} else {
		utils.RespOK(writer, user, "")
	}
}
