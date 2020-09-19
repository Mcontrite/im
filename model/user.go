package model

import "time"

type User struct {
	ID       int64     `xorm:"pk autoincr bigint(20)" form:"id" json:"id"`
	Username string    `xorm:"varchar(20)" form:"username" json:"username"`
	Password string    `xorm:"varchar(40)" form:"password" json:"password"`
	Phone    string    `xorm:"varchar(20)" form:"phone" json:"phone"`
	Avatar   string    `xorm:"varchar(150)" form:"avatar" json:"avatar"`
	Sex      string    `xorm:"varchar(2)" form:"sex" json:"sex"`
	IsOnline int       `xorm:"int(10)" form:"online" json:"online"`   //是否在线
	Salt     string    `xorm:"varchar(10)" form:"salt" json:"-"`      //加盐随机字符串
	Token    string    `xorm:"varchar(40)" form:"token" json:"token"` //前端鉴权因子
	CreateAt time.Time `xorm:"datetime" form:"createat" json:"createat"`
}
