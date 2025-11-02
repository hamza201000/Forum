package forum

import (
	"database/sql"
	"fmt"
)

func tableExists(db *sql.DB, tableName string) bool {
	query := `SELECT name FROM sqlite_master WHERE type='table' AND name=?;`
	row := db.QueryRow(query, tableName)

	var name string
	err := row.Scan(&name)
	return err == nil // true if table exists
}
func WriteCategories() {
	categories := []string{"Technology", "Science", "Education", "Engineering", "Entertainment"}
	insertcategorie := `INSERT INTO categories(categorie) VALUES (?)`

	for _, catcategorie := range categories {
		stmt, err := db.Prepare(insertcategorie)
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
