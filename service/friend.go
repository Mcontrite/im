package service

import (
	"errors"
	"im2/model"
	"time"
)

func SearchFriends(userid int) (users []model.User) {
	friends := make([]model.Friend, 0)
	ids := make([]int, 0)
	DB.Where("user_id=?", userid).Find(&friends)
	for _, v := range friends {
		ids = append(ids, v.User2ID)
	}
	if len(ids) == 0 {
		return
	}
	DB.In("id", ids).Find(&users)
	return
}

func AddFriend(userid, user2id int) error {
	if userid == user2id {
		return errors.New("Can not add yourself as friend")
	}
	friend := model.Friend{}
	DB.Where("user_id=?", userid).And("user2_id=?", user2id).Get(&friend)
	if friend.ID > 0 {
		return errors.New("You have add this friend already")
	}
	session := DB.NewSession()
	session.Begin()
	_, err2 := session.InsertOne(model.Friend{
		UserID:    userid,
		User2ID:   user2id,
		CreatedAt: time.Now(),
	})
	_, err3 := session.InsertOne(model.Friend{
		UserID:    user2id,
		User2ID:   userid,
		CreatedAt: time.Now(),
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
