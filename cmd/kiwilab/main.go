package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type Bookmark struct {
	ID    int
	Title string
	URL   string
	Icon  string
}

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./kiwilab.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Asegurar que la tabla existe y tiene el campo position para el ordenado
	db.Exec(`CREATE TABLE IF NOT EXISTS bookmarks (
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		title TEXT, 
		url TEXT, 
		icon TEXT, 
		position INTEGER DEFAULT 0
	)`)

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/add-bookmark", handleAdd)
	http.HandleFunc("/edit-bookmark", handleEdit)
	http.HandleFunc("/delete-bookmark", handleDelete)
	http.HandleFunc("/reorder-bookmarks", handleReorder)

	log.Println("🚀 kiwiLAB activo en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	rows, _ := db.Query("SELECT id, title, url, icon FROM bookmarks ORDER BY position ASC")
	defer rows.Close()
	var bms []Bookmark
	for rows.Next() {
		var b Bookmark
		rows.Scan(&b.ID, &b.Title, &b.URL, &b.Icon)
		bms = append(bms, b)
	}
	tmpl, _ := template.ParseFiles("ui/templates/dashboard.html")
	tmpl.Execute(w, bms)
}

func handleReorder(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		ids := strings.Split(r.FormValue("ids"), ",")
		for i, id := range ids {
			db.Exec("UPDATE bookmarks SET position=? WHERE id=?", i, id)
		}
		w.WriteHeader(http.StatusOK)
	}
}

func handleAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		db.Exec("INSERT INTO bookmarks (title, url, icon) VALUES (?, ?, ?)", 
			r.FormValue("title"), r.FormValue("url"), r.FormValue("icon"))
		http.Redirect(w, r, "/", 303)
	}
}

func handleEdit(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		db.Exec("UPDATE bookmarks SET title=?, url=?, icon=? WHERE id=?", 
			r.FormValue("title"), r.FormValue("url"), r.FormValue("icon"), r.FormValue("id"))
		http.Redirect(w, r, "/", 303)
	}
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		db.Exec("DELETE FROM bookmarks WHERE id=?", r.FormValue("id"))
		w.WriteHeader(http.StatusOK)
	}
}
