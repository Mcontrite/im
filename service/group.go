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
	if group.LeaderID == 0 {
		err = errors.New("请先登录")
		return
	}
	groups := model.Group{LeaderID: group.LeaderID}
	count, err := DB.Count(&groups)
	if count > 5 {
		err = errors.New("一个用户最多创5个群")
		return groups, err
	} else {
		group.CreateAt = time.Now()
		session := DB.NewSession()
		session.Begin()
		_, err = session.InsertOne(&group)
		if err != nil {
			session.Rollback()
			return groups, err
		}
		_, err = session.InsertOne(
			model.UserGroup{
				UserID:   group.LeaderID,
				GroupID:  group.ID,
				CreateAt: time.Now(),
			})
		if err != nil {
			session.Rollback()
		} else {
			session.Commit()
		}
		return groups, err
	}
}

func SearchGroup(userid int64) []model.Group {
	usergroups := make([]model.UserGroup, 0)
	ids := make([]int64, 0)
	DB.Where("userid=?", userid).Find(&usergroups)
	for _, v := range usergroups {
		ids = append(ids, v.GroupID)
	}
	groups := make([]model.Group, 0)
	if len(ids) == 0 {
		return groups
	}
	DB.In("id", ids).Find(&groups)
	return groups
}

func SearchGroupIds(userid int64) (ids []int64) {
	//todo 获取用户全部群ID
	usergroups := make([]model.UserGroup, 0)
	ids = make([]int64, 0)
	DB.Where("userid=?", userid).Find(&usergroups)
	for _, v := range usergroups {
		ids = append(ids, v.GroupID)
	}
	return ids
}

//加群
func JoinGroup(userid, groupid int64) error {
	usergroup := model.UserGroup{
		UserID:  userid,
		GroupID: groupid,
	}
	DB.Get(&usergroup)
	if usergroup.ID == 0 {
		usergroup.CreateAt = time.Now()
		_, err := DB.InsertOne(usergroup)
		return err
	} else {
		return errors.New("已加入这个群")
	}
}
