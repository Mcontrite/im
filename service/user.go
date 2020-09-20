package service

import (
	"errors"
	"fmt"
	"im/model"
	"im/utils"
	"math/rand"
	"strings"
	"time"
)

func Register(phone, password, username, avatar, sex string) (user model.User, err error) {
	_, err = DB.Where("phone=?", phone).Get(&user)
	if err != nil {
		return
	}
	if user.Id > 0 {
		return user, errors.New("手机号已经注册")
	}
	user.Phone = phone
	user.Avatar = avatar
	user.Username = username
	user.Sex = sex
	user.Salt = fmt.Sprintf("%06d", rand.Int31n(10000))
	user.Password = utils.MD5Password(password, user.Salt)
	user.Token = fmt.Sprintf("%08d", rand.Int31())
	user.CreateAt = time.Now()
	_, err = DB.InsertOne(&user)
	//前端恶意插入特殊字符?
	//数据库连接操作失败?
	return
}

func Login(phone, password string) (user model.User, err error) {
	_, err = DB.Where("phone=?", phone).Get(&user)
	if user.Id == 0 {
		return user, errors.New("该用户不存在")
	}
	if !utils.ValidatePassword(password, user.Salt, user.Password) {
		return user, errors.New("密码不正确")
	}
	//刷新token
	str := fmt.Sprintf("%d", time.Now().Unix())
	token := strings.ToUpper(utils.MD5(str))
	user.Token = token
	_, err = DB.Where("id=?", user.Id).Cols("token").Update(&user)
	return
}

//查找某个用户
func GetUserByID(id int64) (user model.User) {
	DB.Where("id=?", id).Get(&user)
	return
}
