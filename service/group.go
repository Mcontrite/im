package service

import (
	"errors"
	"im/model"
	"time"
)

func CreateGroup(group model.Group) (ret model.Group, err error) {
	if len(group.Groupname) == 0 {
		err = errors.New("群名称不能为空")
		return
	}
	if group.LeaderId == 0 {
		err = errors.New("请先登录")
		return
	}
	oldgroup := model.Group{LeaderId: group.LeaderId}
	count, err := DB.Count(&oldgroup)
	if count > 5 {
		err = errors.New("一个用户最多创5个群")
		return oldgroup, err
	}
	group.CreateAt = time.Now()
	session := DB.NewSession()
	session.Begin()
	_, err = session.InsertOne(&group)
	if err != nil {
		session.Rollback()
		return oldgroup, err
	}
	_, err = session.InsertOne(
		model.UserGroup{
			UserId:   group.LeaderId,
			GroupId:  group.Id,
			CreateAt: time.Now(),
		})
	if err != nil {
		session.Rollback()
	} else {
		session.Commit()
	}
	return oldgroup, err
}

func SearchGroupsIDs(userid int64) (ids []int64) {
	usergroups := make([]model.UserGroup, 0)
	DB.Where("user_id=?", userid).Find(&usergroups)
	for _, v := range usergroups {
		ids = append(ids, v.GroupId)
	}
	return
}

func SearchGroups(userid int64) (groups []model.Group) {
	ids := SearchGroupsIDs(userid)
	if len(ids) == 0 {
		return
	}
	DB.In("id", ids).Find(&groups)
	return
}

func JoinGroup(userid, groupid int64) error {
	usergroup := model.UserGroup{
		UserId:  userid,
		GroupId: groupid,
	}
	DB.Get(&usergroup)
	if usergroup.Id == 0 {
		usergroup.CreateAt = time.Now()
		_, err := DB.InsertOne(usergroup)
		return err
	}
	return errors.New("已加入这个群")
}
