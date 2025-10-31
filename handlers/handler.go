package forum

import (
	"database/sql"
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

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/login" {
		return
	}
	if r.Method == http.MethodGet {
		tmp, err := template.ParseFiles("tamplates/login.html")
		if err != nil {
			return
		}
		tmp.Execute(w, nil)

	} else if r.Method == http.MethodPost {
		user := r.FormValue("username")
		password := r.FormValue("password")
		var hashedPassword string

		err := db.QueryRow(`SELECT password FROM users WHERE username = ?`, user).Scan(&hashedPassword)
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		} else if err != nil {
			return
		}
		if hashedPassword != password {

			http.Redirect(w, r, "/login", http.StatusSeeOther)

		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		return
	}
}

func HandlePost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post" {
		return
	} else if r.Method != http.MethodGet {
		return
	}
	tmp, err := template.ParseFiles("tamplates/post.html")
	if err != nil {
		return
	}
	tmp.Execute(w, nil)
}
func HandleAddPost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/addpost" {
		return
	} else if r.Method != http.MethodPost {
		return
	}
	title := r.FormValue("title")
	content := r.FormValue("content")
	insrtpost := `INSERT INTO posts (title,content) VALUES (?, ?) `
	stmt, err := db.Prepare(insrtpost)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(title, content)
	if err != nil {
		return
	}
	tmp, err := template.ParseFiles("tamplates/addpost.html")
	if err != nil {
		return
	}
	tmp.Execute(w, nil)
}
