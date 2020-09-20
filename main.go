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
		http.HandleFunc(tplname, func(w http.ResponseWriter, req *http.Request) {
			fmt.Println("parse     " + v.Name() + "==" + tplname)
			err := tpl.ExecuteTemplate(w, tplname, nil)
			if err != nil {
				log.Fatal(err.Error())
			}
		})
	}
}

func main() {
	//绑定请求和处理函数
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/user/login.html", http.StatusFound)
	})
	http.HandleFunc("/user/register", ctr.UserRegister)
	http.HandleFunc("/user/login", ctr.UserLogin)
	http.HandleFunc("/user/find", ctr.FindUserByID)
	http.HandleFunc("/group/create", ctr.CreateGroup)
	http.HandleFunc("/group/join", ctr.JoinGroup)
	http.HandleFunc("/group/load", ctr.LoadGroups)
	http.HandleFunc("/friend/add", ctr.AddFriend)
	http.HandleFunc("/friend/load", ctr.LoadFriends)
	http.HandleFunc("/im", ctr.Im)
	http.HandleFunc("/attach/upload", ctr.Upload)

	// 指定静态文件目录
	http.Handle("/static/", http.FileServer(http.Dir(".")))
	http.Handle("/upload/", http.FileServer(http.Dir(".")))
	RegisterView()

	fmt.Println("run at :3333")
	http.ListenAndServe(":3333", nil)
}
