package service

import (
	"errors"
	"im/model"
	"time"
)

//查找好友
func SearchFriend(userid int64) []model.User {
	friends := make([]model.Friend, 0)
	ids := make([]int64, 0)
	DB.Where("userid=?", userid).Find(&friends)
	for _, v := range friends {
		ids = append(ids, v.User2ID)
	}
	users := make([]model.User, 0)
	if len(ids) == 0 {
		return users
	}
	DB.In("id", ids).Find(&users)
	return users
}

//添加好友
func AddFriend(userid, user2id int64) error {
	if userid == user2id {
		return errors.New("不能添加自己为好友")
	}
	friend := model.Friend{}
	DB.Where("userid=?", userid).And("user2id=?", user2id).Get(&friend)
	if friend.ID > 0 {
		return errors.New("该用户已经被添加过啦")
	}
	session := DB.NewSession() //事务
	session.Begin()
	_, err2 := session.InsertOne(model.Friend{ //插入自己的数据
		UserID:   userid,
		User2ID:  user2id,
		CreateAt: time.Now(),
	})
	_, err3 := session.InsertOne(model.Friend{ //插入对方的数据
		UserID:   user2id,
		User2ID:  userid,
		CreateAt: time.Now(),
	})
	if err2 == nil && err3 == nil {
		session.Commit()
		return nil
	} else {
		session.Rollback()
		if err2 != nil {
			return err2
		} else {
			return err3
		}
	}
}
