package forum

import (
	"net/http"
	"text/template"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		return
	} else if r.Method != http.MethodGet {
		return
	}
	tmpl, err := template.ParseFiles("tamplates/index.html")
	if err != nil {
		return
	}
	tmpl.Execute(w, nil)
}
func Handleregister(w http.ResponseWriter,r *http.Request){

	if r.URL.Path!="/register"{
		return
	}else if r.Method!=http.MethodGet{
		return
	}
	tmpl,err:=template.ParseFiles("tamplates/register.html")
	if err!=nil{
		return
	}
	tmpl.Execute(w,nil)


}
