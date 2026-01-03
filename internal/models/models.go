package models

import "time"

// User representa a los administradores o usuarios limitados
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"` // Hash Bcrypt
	Avatar    string    `json:"avatar"`   // Base64 encoded image
	Role      string    `json:"role"`     // "admin" o "user"
	Language  string    `json:"language"` // Código de idioma ISO (en, es, fr, etc.)
	Theme     string    `json:"theme"`    // "dark" o "light"
	CreatedAt time.Time `json:"created_at"`
}

// Bookmark representa un marcador
type Bookmark struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	URL      string `json:"url"`
	Icon     string `json:"icon"`
	Position int    `json:"position,omitempty"`
}
