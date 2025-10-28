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

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/register" {
		return
	} else if r.Method != http.MethodGet {
		return
	}
	tmpl, err := template.ParseFiles("tamplates/register.html")
	if err != nil {
		return
	}
	tmpl.Execute(w, nil)
}

func HandlerDataRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}
	user := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	insertsql := `INSERT INTO users (username, email, password) VALUES (?, ?, ?)`
	stmt, err := db.Prepare(insertsql)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(user, email, password)
	if err != nil {
		return
	}
}
