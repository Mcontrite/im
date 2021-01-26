package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func InitTemplate() {
	tpl, err := template.ParseGlob("view/**/*")
	if err != nil {
		fmt.Println("Init template err: ", err)
		return
	}
	for _, v := range tpl.Templates() {
		tplname := v.Name()
		http.HandleFunc(tplname, func(w http.ResponseWriter, req *http.Request) {
			err := tpl.ExecuteTemplate(w, tplname, nil)
			if err != nil {
				fmt.Println("Excute template err: ", err)
				return
			}
		})
	}
}

func main() {
	http.Handle("/static/", http.FileServer(http.Dir(".")))
	http.Handle("/upload/", http.FileServer(http.Dir(".")))
	InitTemplate()

	http.ListenAndServe(":3030", nil)
	fmt.Println("IM running at 3333...")
}
