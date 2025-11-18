package backend

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type DataComment struct {
	Username  string
	Content   string
	CreatedAt string
}
type Datapost struct {
	Title    string
	Content  string
	IdPost   int
	Comments []DataComment
	Likes    string
	Dislikes string
	Username string
}

type Message_Error struct {
	Status  int
	Message string
}

var CategoriesId = make(map[string]int)

func tableExists(DB *sql.DB, tableName string) bool {
	query := `SELECT name FROM sqlite_master WHERE type='table' AND name=?;`
	row := DB.QueryRow(query, tableName)
	var name string
	err := row.Scan(&name)
	return err == nil
}

func InsertCategorie() {
	categories := []string{"Technology", "Science", "Education", "Engineering", "Entertainment"}
	i := 1
	for _, categorie := range categories {
		CategoriesId[categorie] = i
		i++
	}
}

func InsertNamePost(DB *sql.DB) []string {
	names := []string{"title", "content", "category_ids"}
	insertcategorie := `INSERT INTO method_post(name) VALUES (?)`

	for _, name := range names {
		stmt, err := DB.Prepare(insertcategorie)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		_, err = stmt.Exec(name)
		if err != nil {
			fmt.Println(err)
			return nil
		}
	}
	return names
}

func WriteCategories(DB *sql.DB) {
	categories := []string{"Technology", "Science", "Education", "Engineering", "Entertainment"}
	insertcategorie := `INSERT INTO categories(categorie) VALUES (?)`

	for _, catcategorie := range categories {
		stmt, err := DB.Prepare(insertcategorie)
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = stmt.Exec(catcategorie)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func checkuser(DB *sql.DB, userid int64) bool {
	var token string
	err := DB.QueryRow(`SELECT token FROM sessions WHERE user_id = ? `, userid).Scan(&token)
	if err == sql.ErrNoRows {
		return true
	}
	_, err = DB.Exec(`DELETE FROM sessions WHERE user_id = ? `, userid)
	if err != nil {
		log.Fatal(err)
	}
	return false
}

func InsertCategoriId(DB *sql.DB, post_id int64, categories []string) {
	var categorie_id int
	for _, categorie := range categories {
		fmt.Println("Inserting category:", categorie, "for post ID:", post_id)
		err := DB.QueryRow(`SELECT id FROM categories WHERE categorie = ?`, categorie).Scan(&categorie_id)
		if err != nil {
			return
		}
		_, err = DB.Exec("INSERT INTO post_categories (post_id,category_id) VALUES (?,?)", post_id, categorie_id)
		if err != nil {
			return
		}
	}
}

func GetPost(DB *sql.DB, category, username string, UserId int64) []Datapost {
	posts := []Datapost{}
	fmt.Println("Fetching posts for category:", category, "and username:", UserId)
	Categorie_Id := CategoriesId[category]

	var row *sql.Rows
	var err error
	if category == "" {
		row, err = DB.Query(`SELECT title,content,id FROM posts`)
	} else if category == username {
		row, err = DB.Query(`SELECT title,content,id FROM posts WHERE user_id=?`, UserId)
	} else if category == "liked" {

		row, err = DB.Query(`SELECT posts.title,posts.content,posts.id
	FROM posts
	JOIN likes ON likes.post_id=posts.id
	WHERE likes.kind=1 AND likes.user_id=?
	`, UserId)
	}else{

		row, err = DB.Query(`SELECT posts.title,posts.content,posts.id
	FROM posts
	JOIN post_categories ON post_categories.post_id=posts.id
	WHERE post_categories.category_id=?
	`, Categorie_Id)
	}

	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer row.Close()
	for row.Next() {
		var post Datapost

		if err := row.Scan(&post.Title, &post.Content, &post.IdPost); err != nil {
			log.Fatal(err)
			return nil
		}
		post.Username = username
		post.Comments = GetComment(DB, post.IdPost)
		post.Likes, post.Dislikes = GetCountLike(DB, post.IdPost)
		posts = append(posts, post)

	}
	if err = row.Err(); err != nil {
		log.Fatal(err)
		return nil
	}
	return posts
}

func GetPostById(DB *sql.DB, PostId int) []Datapost {
	posts := []Datapost{}
	row, err := DB.Query(`SELECT title,content,id FROM posts WHERE id =?`, PostId)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer row.Close()
	for row.Next() {
		var post Datapost
		if err := row.Scan(&post.Title, &post.Content, &post.IdPost); err != nil {
			log.Fatal(err)
			return nil
		}

		posts = append(posts, post)

	}
	if err = row.Err(); err != nil {
		log.Fatal(err)
		return nil
	}
	return posts
}

func GetComment(DB *sql.DB, PostId int) []DataComment {
	Comments := []DataComment{}
	rows, err := DB.Query(`
			SELECT u.username, c.comment, c.created_at
			FROM comments c
			JOIN users u ON u.id = c.user_id
			WHERE c.post_id = ?
			ORDER BY c.created_at DESC`, PostId)
	if err != nil {
		fmt.Println(err)
		// http.Error(w, "Erreur base de données", http.StatusInternalServerError)
		return nil
	}
	defer rows.Close()
	// On renvoie du HTML directement pour mettre à jour dynamiquement

	for rows.Next() {
		var DataComments DataComment
		rows.Scan(&DataComments.Username, &DataComments.Content, &DataComments.CreatedAt)
		Comments = append(Comments, DataComments)
	}

	return Comments
}

func GetCountLike(DB *sql.DB, PostId int) (string, string) {
	// Compte les likes et dislikes
	var likesCount, dislikesCount int
	DB.QueryRow("SELECT COUNT(*) FROM likes WHERE post_id = ? AND kind = 1", PostId).Scan(&likesCount)
	DB.QueryRow("SELECT COUNT(*) FROM likes WHERE post_id = ? AND kind = -1", PostId).Scan(&dislikesCount)
	likes := strconv.Itoa(likesCount)
	dislikes := strconv.Itoa(dislikesCount)
	return likes, dislikes
}

func Render(w http.ResponseWriter, status int) {
	// Parse the error template file
	tmp, err := template.ParseFiles("templates/errorpage.html")
	// Set the HTTP status code in the response
	w.WriteHeader(status)
	// If there is an error loading the template, show a simple error message
	if err != nil {
		http.Error(w, "page not found", http.StatusNotFound)
		return
	}
	// Prepare the error message based on the status code
	message := ""
	switch status {
	case 400:
		message = "Bad Request."
	case 404:
		message = "Not Found."
	case 405:
		message = "Status Method Not Allowed."
	case 403:
		message = "Access denied: you don’t have permission to view this resource."
	default:
		message = "Status Internal Server Error"
	}
	// Create a struct with status and message to pass to the template
	mes := Message_Error{
		Status:  status,
		Message: message,
	}
	// Execute the template and display the error page
	tmp.Execute(w, mes)
}

func CheckDataPost(DB *sql.DB, r *http.Request, errorMsg string) PostPageData {
	var post []Datapost
	title := r.FormValue("title")
	content := r.FormValue("content")
	userid := GetUserIDFromRequest(DB, r)
	username := ""
	if userid != 0 {
		err := DB.QueryRow("SELECT username FROM users WHERE id = ?", userid).Scan(&username)
		if err != nil {
			fmt.Print(err)
			// return nil
		}
	}

	LastPath, err := r.Cookie("LastPath")
	if err != nil {
	}

	if LastPath.Value != "/post" {
		lastCategories := strings.Split(LastPath.Value, "=")
		post = GetPost(DB, lastCategories[len(lastCategories)-1], username, userid)
	} else {
		post = GetPost(DB, "", username, userid)
	}

	PageData := &PostPageData{
		Error:         errorMsg,
		Popup:         true,
		Posts:         post,
		Username:      username,
		Cachetitle:    title,
		Cacheconetent: content,
		Categories:    []string{"Technology", "Science", "Education", "Engineering", "Entertainment"},
		Path:          LastPath.Value,
	}
	return *PageData
}

func CheckFiltere(query string, username string) bool {
	Filtre := strings.Split(query, "=")
	if Filtre[0] != "Categories" {
		return false
	}
	categories := []string{"liked", "Technology", "Science", "Education", "Engineering", "Entertainment", username}
	for _, categorie := range categories {
		if categorie == Filtre[len(Filtre)-1] {
			return true
		}
	}
	return false
}
