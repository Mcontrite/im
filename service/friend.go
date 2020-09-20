package service

import (
	"errors"
	"im/model"
	"time"
)

//查找好友
func SearchFriends(userid int64) (users []model.User) {
	friends := make([]model.Friend, 0)
	ids := make([]int64, 0)
	DB.Where("user_id=?", userid).Find(&friends)
	for _, v := range friends {
		ids = append(ids, v.User2Id)
	}
	if len(ids) == 0 {
		return
	}
	DB.In("id", ids).Find(&users)
	return
}

//添加好友
func AddFriend(userid, user2id int64) error {
	if userid == user2id {
		return errors.New("不能添加自己为好友")
	}
	friend := model.Friend{}
	DB.Where("user_id=?", userid).And("user2_id=?", user2id).Get(&friend)
	if friend.Id > 0 {
		return errors.New("已添加这个好友")
	}
	session := DB.NewSession() //事务
	session.Begin()
	_, err2 := session.InsertOne(model.Friend{ //插入自己的数据
		UserId:   userid,
		User2Id:  user2id,
		CreateAt: time.Now(),
	})
	_, err3 := session.InsertOne(model.Friend{ //插入对方的数据
		UserId:   user2id,
		User2Id:  userid,
		CreateAt: time.Now(),
	})
	if err2 == nil && err3 == nil {
		session.Commit()
		return nil
	}
	session.Rollback()
	if err2 != nil {
		return err2
	}
	return err3
}
