package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3" // Necesitaremos descargar esto luego
	"log"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./kiwilab.db")
	if err != nil {
		log.Fatal(err)
	}

	// Crear tabla de ejemplo si no existe
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS bookmarks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		url TEXT,
		icon TEXT,
		x INTEGER,
		y INTEGER,
		w INTEGER,
		h INTEGER
	);`
	_, err = DB.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}
