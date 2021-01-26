package service

import (
	"errors"
	"im2/model"
	"time"
)

func CreateGroup(group model.Group) (ret model.Group, err error) {
	if len(group.Groupname) == 0 {
		err = errors.New("Groupname can not be empty")
		return
	}
	if group.LeaderID == 0 {
		err = errors.New("Please login first")
		return
	}
	oldgroup := model.Group{LeaderID: group.LeaderID}
	count, err := DB.Count(&oldgroup)
	if count > 5 {
		err = errors.New("Each user can create up to 5 groups")
		return oldgroup, err
	}
	group.CreatedAt = time.Now()
	session := DB.NewSession()
	session.Begin()
	_, err = session.InsertOne(&group)
	if err != nil {
		session.Rollback()
		return oldgroup, err
	}
	_, err = session.InsertOne(model.UserGroup{
		UserID:    group.LeaderID,
		GroupID:   group.ID,
		CreatedAt: time.Now(),
	})
	if err != nil {
		session.Rollback()
	} else {
		session.Commit()
	}
	return oldgroup, err
}

func SearchGroupID(userid int) (ids []int) {
	usergroups := make([]model.UserGroup, 0)
	DB.Where("user_id=?", userid).Find(&usergroups)
	for _, v := range usergroups {
		ids = append(ids, v.GroupID)
	}
	return
}

func SearchGroups(userid int) (groups []model.Group) {
	ids := SearchGroupsID(userid)
	if len(ids) == 0 {
		return
	}
	DB.In("id", ids).Find(&groups)
	return
}

func JoinGroup(userid, groupid int) error {
	usergroup := model.UserGroup{
		UserID:  userid,
		GroupID: groupid,
	}
	DB.Get(&usergroup)
	if usergroup.ID == 0 {
		usergroup.CreatedAt = time.Now()
		_, err := DB.InsertOne(usergroup)
		return err
	}
	return errors.New("You are already in this group")
}
