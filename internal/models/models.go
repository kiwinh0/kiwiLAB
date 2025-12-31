package models

import "time"

// User representa a los administradores o usuarios limitados
type User struct {
	ID       int
	Username string
	Password string // Hash Bcrypt
	Role     string // "admin" o "user"
}

// Board es cada dashboard dinámico
type Board struct {
	ID        int
	UserID    int
	Name      string
	IsDefault bool
	Settings  string // JSON con colores, columnas, etc.
}

// Bookmark es cada acceso directo
type Bookmark struct {
	ID        int
	BoardID   int
	SectionID int
	Title     string
	URL       string
	Icon      string
	Position  int    // Para el "Arrastrar y Soltar"
	Size      string // "small", "medium", "large"
}
