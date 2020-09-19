package service

import (
	"errors"
	"fmt"
	"im/model"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

const (
	dbType  = "mysql"
	dsName  = "root:123456@(127.0.0.1:3306)/im?charset=utf8"
	showSQL = true
	maxCons = 10
	noError = "noerror" //没有错误
)

var DB *xorm.Engine

//初始化数据库
func init() {
	err := errors.New(noError)
	DB, err = xorm.NewEngine(dbType, dsName)
	if err != nil && err.Error() != noError {
		log.Fatal(err.Error())
	}
	DB.ShowSQL(showSQL)         //是否显示SQL语句
	DB.SetMaxOpenConns(maxCons) //最大打开的连接数
	//自动建表
	DB.Sync2(
		new(model.User),
		new(model.Group),
		new(model.Friend),
		new(model.UserGroup),
	)
	fmt.Println("Init Database OK")
}
