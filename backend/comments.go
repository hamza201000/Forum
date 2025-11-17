package backend

import (
	"database/sql"
	"fmt"
	"net/http"
)

func HandleAddComment(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Vérifie la méthode
		if r.Method != http.MethodPost {
			http.Error(w, "/login", http.StatusMethodNotAllowed)
			return
		}

		// Vérifie la session utilisateur
		userID, err := GetUserIDFromSession(r, db)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Récupère les champs du formulaire
		postID := r.FormValue("post_id")
		content := r.FormValue("comment")
		if postID == "" || content == "" {
			http.Error(w, "Champs manquants", http.StatusBadRequest)
			return
		}
		var dummy int // or matching columns
		err = db.QueryRow("SELECT id FROM posts WHERE id = ?", postID).Scan(&dummy)

		if err == sql.ErrNoRows {
			http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
			return
		}
		// Insérer le commentaire
		_, err = db.Exec("INSERT INTO comments (post_id, user_id, comment) VALUES (?, ?, ?)", postID, userID, content)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Erreur base de données", http.StatusInternalServerError)
			return
		}

		// Récupère la liste mise à jour des commentaires
		rows, err := db.Query(`
			SELECT u.username, c.comment, c.created_at
			FROM comments c
			JOIN users u ON u.id = c.user_id
			WHERE c.post_id = ?
			ORDER BY c.created_at DESC`, postID)
		if err != nil {
			http.Error(w, "Erreur base de données", http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		http.Redirect(w, r, r.Referer()+"#"+postID, http.StatusSeeOther)
	}
}
