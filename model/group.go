package model

import "time"

type Group struct {
	Id        int64     `xorm:"pk autoincr bigint(20)" form:"id" json:"id"`
	LeaderId  int64     `xorm:"bigint(20)" form:"leaderid" json:"leaderid"` //群主ID
	Groupname string    `xorm:"varchar(30)" form:"groupname" json:"groupname"`
	Avatar    string    `xorm:"varchar(250)" form:"avatar" json:"avatar"` //群logo
	CreateAt  time.Time `xorm:"datetime" form:"createat" json:"createat"`
}
