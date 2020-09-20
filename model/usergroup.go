package model

import "time"

type UserGroup struct {
	Id       int64     `xorm:"pk autoincr bigint(20)" form:"id" json:"id"`
	UserId   int64     `xorm:"bigint(20)" form:"userid" json:"userid"`   // 记录是谁的
	GroupId  int64     `xorm:"bigint(20)" form:"groupid" json:"groupid"` // 对端信息
	CreateAt time.Time `xorm:"datetime" form:"createat" json:"createat"` // 创建时间
}
