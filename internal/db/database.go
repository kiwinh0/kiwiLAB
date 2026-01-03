package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(dbPath string) {
	log.Println("InitDB called with path:", dbPath)
	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return
	}

	log.Println("Database connection established")

	// Crear tabla users
	userStmt := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		avatar TEXT,
		role TEXT DEFAULT 'user',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = DB.Exec(userStmt)
	if err != nil {
		log.Printf("Error creating users table: %q: %s\n", err, userStmt)
	} else {
		log.Println("Users table created or already exists")
	}

	// Add avatar column if it doesn't exist (migration)
	_, err = DB.Exec(`ALTER TABLE users ADD COLUMN avatar TEXT`)
	if err != nil {
		log.Printf("Error adding avatar column (this is expected if column already exists): %v", err)
	} else {
		log.Println("Avatar column added or already exists")
	}

	// Add language preference column
	_, err = DB.Exec(`ALTER TABLE users ADD COLUMN language TEXT DEFAULT 'en'`)
	if err != nil {
		log.Printf("Error adding language column (this is expected if column already exists): %v", err)
	}

	// Add theme preference column
	_, err = DB.Exec(`ALTER TABLE users ADD COLUMN theme TEXT DEFAULT 'dark'`)
	if err != nil {
		log.Printf("Error adding theme column (this is expected if column already exists): %v", err)
	}

	// Crear tabla bookmarks
	bookmarkStmt := `
	CREATE TABLE IF NOT EXISTS bookmarks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		url TEXT,
		icon TEXT,
		position INTEGER DEFAULT 0
	);`
	_, err = DB.Exec(bookmarkStmt)
	if err != nil {
		log.Printf("Error creating bookmarks table: %q: %s\n", err, bookmarkStmt)
	} else {
		log.Println("Bookmarks table created or already exists")
	}

	// Crear usuario por defecto si no existe
	// NOTA: Este código se mantiene por retrocompatibilidad, pero el setup wizard
	// ahora es el método preferido para crear el primer usuario
	var count int
	err = DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		log.Printf("Error checking users count: %v", err)
	} else {
		log.Printf("Current users count: %d", count)
	}

	if err == nil && count == 0 {
		log.Println("No users found - setup wizard will be shown on first access")
	} else if err == nil && count > 0 {
		log.Printf("Users exist, normal login will be used")
	}

	log.Println("InitDB completed successfully")
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
