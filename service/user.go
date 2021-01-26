package service

import (
	"errors"
	"fmt"
	"im2/model"
	"im2/utils"
	"math/rand"
	"strings"
	"time"
)

func Register(username, password string) (user model.User, err error) {
	_, err = DB.Where("username=?", username).Get(&user)
	if err != nil {
		return
	}
	if user.ID > 0 {
		return
	}
	user.Username = username
	user.Salt = fmt.Sprintf("%06d", rand.Int31n(10000))
	user.Password = utils.MD5(password + user.Salt)
	user.Token = fmt.Sprintf("%08d", rand.Int31())
	user.CreatedAt = time.Now()
	_, err = DB.InsertOne(&user)
	return
}

func Login(username, password string) (user model.User, err error) {
	_, err = DB.Where("username=?", username).Get(&user)
	if user.ID == 0 {
		return user, errors.New("The User is not exist")
	}
	if utils.MD5(password+user.Salt) != user.Password {
		return user, errors.New("Password is error")
	}
	str := fmt.Sprintf("%d", time.Now().Unix())
	user.Token = strings.ToUpper(utils.MD5(str))
	_, err = DB.Where("id=?", user.ID).Cols("token").Update(&user)
	return
}

func GetUserByID(id int) (user model.User) {
	DB.Where("id=?", id).Get(&user)
	return
}
