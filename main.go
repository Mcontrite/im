package main

import (
	"fmt"
	"html/template"
	ctr "im/controller"
	"log"
	"net/http"
)

//注册模板
func RegisterView() {
	tpl, err := template.ParseGlob("view/**/*")
	if nil != err {
		log.Fatal(err)
	}
	//通过for循环做好映射
	for _, v := range tpl.Templates() {
		tplname := v.Name()
		fmt.Println("HandleFunc     " + v.Name())
		http.HandleFunc(tplname, func(w http.ResponseWriter, request *http.Request) {
			fmt.Println("parse     " + v.Name() + "==" + tplname)
			err := tpl.ExecuteTemplate(w, tplname, nil)
			if err != nil {
				log.Fatal(err.Error())
			}
		})
	}
}

//注册首页自动跳转
func RegisterIndex() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/user/login.shtml", http.StatusFound) //跳转到百度
	})
}

func main() {
	//绑定请求和处理函数
	http.HandleFunc("/user/register", ctr.UserRegister)
	http.HandleFunc("/user/login", ctr.UserLogin)
	http.HandleFunc("/user/find", ctr.FindUserById)
	http.HandleFunc("/group/create", ctr.CreateGroup)
	http.HandleFunc("/group/join", ctr.JoinGroup)
	http.HandleFunc("/group/load", ctr.LoadGroup)
	http.HandleFunc("/friend/add", ctr.Addfriend)
	http.HandleFunc("/friend/load", ctr.LoadFriend)
	http.HandleFunc("/chat", ctr.Chat)
	http.HandleFunc("/attach/upload", ctr.Upload)

	// 指定静态文件目录
	http.Handle("/asset/", http.FileServer(http.Dir(".")))
	http.Handle("/upload/", http.FileServer(http.Dir(".")))

	RegisterView()
	RegisterIndex()

	fmt.Println("run at :3030")
	http.ListenAndServe(":3030", nil)
}
