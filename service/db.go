package service

import (
	"errors"
	"fmt"
	"im2/model"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var DB *xorm.Engine

func init() {
	err := errors.New("NoError")
	DB, err = xorm.NewEngine("mysql", "root:123456@(127.0.0.1:3306)/im2?charset=utf8")
	if err != nil {
		fmt.Println("Connect DB error: ", err)
		return
	}
	DB.ShowSQL(true)
	DB.SetMaxOpenConns(10)
	DB.Sync2(
		new(model.User),
		new(model.Group),
		new(model.Friend),
		new(model.UserGroup),
	)
	fmt.Println("Init DB OK...")
}
