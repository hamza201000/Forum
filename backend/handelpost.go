package backend

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
)

var lastCategories string

type PostPageData struct {
	Popup         bool
	Username      string
	Posts         []Datapost
	Error         string
	Cachetitle    string
	Cacheconetent string
	Categories    []string
	Path          string
}

func Handler(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			return
		} else if r.Method != http.MethodGet {
			return
		}
		http.Redirect(w, r, "/post", http.StatusSeeOther)
	}
}

func HandlePost(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmp, err := template.ParseFiles("templates/post.html")
		Categories := r.URL.Query().Get("Categories")

		http.SetCookie(w, &http.Cookie{
			Name:  "LastPath",
			Value: r.RequestURI,
			Path:  "/post",
		})
		IdPst, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			IdPst = 0
		}
		if r.URL.Path != "/post" {
			return
		}
		if r.Method == http.MethodGet {
			Data := &PostPageData{}
			userid := GetUserIDFromRequest(DB, r)
			username := ""
			if userid != 0 {
				err := DB.QueryRow("SELECT username FROM users WHERE id = ?", userid).Scan(&username)
				if err != nil {
					fmt.Print(err)
					return
				}
			}
			if IdPst != 0 {
				// fmt.Println("ok")
				post := GetPostById(DB, IdPst)
				if len(post) == 0 {
					Render(w, 404)
					return
				}

				Data = &PostPageData{
					Username:   username,
					Posts:      post,
					Categories: []string{"Technology", "Science", "Education", "Engineering", "Entertainment"},
				}

			} else {

				post := GetPost(DB, Categories, username, userid)
				lastCategories = Categories

				Data = &PostPageData{
					Username:   username,
					Posts:      post,
					Categories: []string{"Technology", "Science", "Education", "Engineering", "Entertainment"},
				}

			}

			if err = tmp.Execute(w, Data); err != nil {
				fmt.Println(err)
				return
			}
			return

		}
		if r.Method == http.MethodPost {
			userId := GetUserIDFromRequest(DB, r)

			if userId == 0 {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			} else {
				title, titleOK := r.Form["title"]
				content, contentOK := r.Form["content"]

				if !titleOK || !contentOK {
					http.Error(w, "Missing required fields", http.StatusBadRequest)
					return
				}
				// category := r.FormValue("category_ids")
				if err := r.ParseForm(); err != nil {
					http.Error(w, "Error parsing form", http.StatusBadRequest)
				}
				if len(title) == 0 {
					errorMsg := "⚠️ Your post needs a title. Please enter one."
					RenderTemplate(w, "post.html", CheckDataPost(DB, r, errorMsg))
					return
				} else if len(content) == 0 {
					errorMsg := "⚠️ Your post needs some content. Please type something to continue."
					RenderTemplate(w, "post.html", CheckDataPost(DB, r, errorMsg))
					return
				}

				// var category []string
				category := r.Form["category_ids"]
				if len(category) == 0 {
					errorMsg := "⚠️ You must choose one category or more"

					RenderTemplate(w, "post.html", CheckDataPost(DB, r, errorMsg))
					return
				}

				insrtpost := `INSERT INTO posts (title,content,user_id) VALUES (?,?,?)`
				stmt, err := DB.Prepare(insrtpost)
				if err != nil {
					fmt.Println(err)
					return
				}

				defer stmt.Close()
				res, err := stmt.Exec(title, content, userId)
				if err != nil {
					fmt.Println(err)
					return
				}
				IdPost, err := res.LastInsertId()
				if err != nil {
					fmt.Println("Error getting last insert ID:", err)
					return
				}
				InsertCategoriId(DB, IdPost, category)
				http.Redirect(w, r, "/post", http.StatusSeeOther)
			}
		}
	}
}

func HandlerStatic(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only allow GET requests
		if r.Method != http.MethodGet {
			Render(w, 405)
			return
		} else {
			// Check if the requested file exists and is not a directory
			info, err := os.Stat(r.URL.Path[1:])

			if err != nil {
				return
			} else if info.IsDir() {
				Render(w, 403)
				return
			} else {
				// Serve the static file
				http.ServeFile(w, r, r.URL.Path[1:])
			}
		}
	}
}
