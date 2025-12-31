package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3" // Importación necesaria para el driver
)

var db *sql.DB

func init() {
	// Inicializamos la base de datos al arrancar
	var err error
	db, err = sql.Open("sqlite3", "./kiwilab.db")
	if err != nil {
		log.Fatal(err)
	}

	// Creamos la tabla si no existe
	statement, _ := db.Prepare(`CREATE TABLE IF NOT EXISTS bookmarks (
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		title TEXT, 
		url TEXT, 
		x INTEGER, 
		y INTEGER, 
		w INTEGER, 
		h INTEGER
	)`)
	statement.Exec()
}

func main() {
	fs := http.FileServer(http.Dir("ui/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("ui/templates/dashboard.html")
		if err != nil {
			http.Error(w, "Error cargando plantilla", 500)
			return
		}
		tmpl.Execute(w, nil)
	})

	log.Println("🚀 kiwiLAB con SQLite activo en http://localhost:8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
