package handlers

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kiwinh0/CodigoSH/internal/middleware"
	"github.com/kiwinh0/CodigoSH/internal/models"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// Version de CodigoSH
const Version = "0.3.0-Beta"

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// detectMimeType detects the MIME type of image data
func detectMimeType(data []byte) string {
	if len(data) < 4 {
		return "image/png" // default fallback
	}

	// Check for PNG
	if data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 {
		return "image/png"
	}

	// Check for JPEG
	if data[0] == 0xFF && data[1] == 0xD8 {
		return "image/jpeg"
	}

	// Check for GIF
	if data[0] == 0x47 && data[1] == 0x49 && data[2] == 0x46 {
		return "image/gif"
	}

	// Check for WebP
	if len(data) >= 12 && string(data[0:4]) == "RIFF" && string(data[8:12]) == "WEBP" {
		return "image/webp"
	}

	return "image/png" // default fallback
}

type Handler struct {
	DB *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	if db == nil {
		logrus.Fatal("Database is nil")
	}
	return &Handler{DB: db}
}

// HandleRoot handles the root path and redirects based on authentication status
func (h *Handler) HandleRoot(w http.ResponseWriter, r *http.Request) {
	logrus.WithFields(logrus.Fields{
		"method": r.Method,
		"path":   r.URL.Path,
	}).Info("HandleRoot called")

	// Check if this is first-time setup (no users in database)
	var userCount int
	err := h.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount)
	if err != nil {
		logrus.WithError(err).Error("Error checking user count")
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// If no users exist, redirect to setup wizard
	if userCount == 0 {
		logrus.Info("No users found, redirecting to setup wizard")
		http.Redirect(w, r, "/setup", http.StatusFound)
		return
	}

	// Check if user has a valid session cookie
	cookie, err := r.Cookie("token")
	if err == nil && cookie.Value != "" {
		// Try to validate the token
		token, err := jwt.ParseWithClaims(cookie.Value, &middleware.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return middleware.JWTSecret, nil
		})

		if err == nil && token.Valid {
			// Valid session, redirect to dashboard
			logrus.Info("Valid session found, redirecting to dashboard")
			http.Redirect(w, r, "/dashboard", http.StatusFound)
			return
		}
	}

	// No valid session, redirect to login
	logrus.Info("No valid session, redirecting to login")
	http.Redirect(w, r, "/login", http.StatusFound)
}

func (h *Handler) HandleIndex(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Error("Panic in HandleIndex: ", r)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}()

	logrus.WithFields(logrus.Fields{
		"method": r.Method,
		"url":    r.URL.Path,
	}).Info("HandleIndex called")

	// Get username from context
	username := r.Context().Value("username")
	if username == nil {
		logrus.Error("No username in context - redirecting to login")
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	usernameStr := username.(string)
	logrus.WithField("username", usernameStr).Info("Username obtained from context")

	// Get user data from database
	var user models.User
	logrus.WithField("username", usernameStr).Info("Querying user data from database")
	err := h.DB.QueryRow("SELECT id, username, COALESCE(avatar, ''), COALESCE(language, ''), COALESCE(theme, '') FROM users WHERE username = ?", usernameStr).Scan(&user.ID, &user.Username, &user.Avatar, &user.Language, &user.Theme)
	if err != nil {
		logrus.WithError(err).WithField("username", usernameStr).Error("Database error in dashboard - user not found")
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	logrus.WithField("user_id", user.ID).Info("User data retrieved successfully")

	logrus.Info("Querying bookmarks from database")
	rows, err := h.DB.Query("SELECT id, title, url, icon, position FROM bookmarks ORDER BY position ASC")
	if err != nil {
		logrus.WithError(err).Error("Error querying bookmarks")
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	logrus.Info("Database query executed successfully")
	var bms []models.Bookmark
	for rows.Next() {
		var b models.Bookmark
		if err := rows.Scan(&b.ID, &b.Title, &b.URL, &b.Icon, &b.Position); err != nil {
			http.Error(w, "Scan error", http.StatusInternalServerError)
			return
		}
		bms = append(bms, b)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, "Rows error", http.StatusInternalServerError)
		return
	}
	logrus.WithField("bookmark_count", len(bms)).Info("Bookmarks processed successfully")

	// Create data structure for template
	data := struct {
		Bookmarks []models.Bookmark
		User      models.User
		Version   string
	}{
		Bookmarks: bms,
		User:      user,
		Version:   Version,
	}

	logrus.Info("Parsing dashboard template")
	tmpl, err := template.ParseFiles("web/templates/dashboard.html")
	if err != nil {
		logrus.WithError(err).Error("Error parsing dashboard template")
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
	logrus.Info("Executing dashboard template")
	err = tmpl.Execute(w, data)
	if err != nil {
		logrus.WithError(err).Error("Error executing dashboard template")
		return
	}
	logrus.Info("Dashboard rendered successfully")
}

func (h *Handler) HandleReorder(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	idsStr := r.FormValue("ids")
	if idsStr == "" {
		http.Error(w, "IDs are required", http.StatusBadRequest)
		return
	}
	ids := strings.Split(idsStr, ",")
	stmt, err := h.DB.Prepare("UPDATE bookmarks SET position=? WHERE id=?")
	if err != nil {
		http.Error(w, "Prepare error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()
	for i, id := range ids {
		if id == "" {
			continue
		}
		_, err = stmt.Exec(i, id)
		if err != nil {
			http.Error(w, "Update error", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) HandleAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	title := r.FormValue("title")
	url := r.FormValue("url")
	icon := r.FormValue("icon")
	if title == "" || url == "" {
		http.Error(w, "Title and URL are required", http.StatusBadRequest)
		return
	}
	stmt, err := h.DB.Prepare("INSERT INTO bookmarks (title, url, icon) VALUES (?, ?, ?)")
	if err != nil {
		http.Error(w, "Prepare error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(title, url, icon)
	if err != nil {
		http.Error(w, "Insert error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (h *Handler) HandleEdit(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := r.FormValue("id")
	title := r.FormValue("title")
	url := r.FormValue("url")
	icon := r.FormValue("icon")
	if id == "" || title == "" || url == "" {
		http.Error(w, "ID, title and URL are required", http.StatusBadRequest)
		return
	}
	stmt, err := h.DB.Prepare("UPDATE bookmarks SET title=?, url=?, icon=? WHERE id=?")
	if err != nil {
		http.Error(w, "Prepare error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(title, url, icon, id)
	if err != nil {
		http.Error(w, "Update error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (h *Handler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}
	stmt, err := h.DB.Prepare("DELETE FROM bookmarks WHERE id=?")
	if err != nil {
		http.Error(w, "Prepare error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(id)
	if err != nil {
		http.Error(w, "Delete error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	logrus.WithFields(logrus.Fields{
		"method": r.Method,
		"url":    r.URL.Path,
	}).Info("HandleLogin called")

	if r.Method == "GET" {
		logrus.Info("Serving login page")

		// Check if this is first-time setup (no users in database)
		var userCount int
		err := h.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount)
		if err != nil {
			logrus.WithError(err).Error("Error checking user count")
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		// If no users exist, redirect to setup wizard
		if userCount == 0 {
			logrus.Info("No users found in login page, redirecting to setup wizard")
			http.Redirect(w, r, "/setup", http.StatusFound)
			return
		}

		// Check for error or message in URL parameters
		errorMsg := r.URL.Query().Get("error")
		successMsg := r.URL.Query().Get("message")

		// Try to get language from cookie, default to 'es'
		language := "es"
		if langCookie, err := r.Cookie("currentLanguage"); err == nil {
			language = langCookie.Value
		}

		data := map[string]interface{}{
			"Error":   errorMsg,
			"Message": successMsg,
			"User": models.User{
				Language: language, // Use saved language or default
			},
		}

		tmpl, err := template.ParseFiles("web/templates/login.html")
		if err != nil {
			logrus.WithError(err).Error("Error parsing login template")
			http.Error(w, "Template error", http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, data)
		if err != nil {
			logrus.WithError(err).Error("Error executing login template")
			http.Error(w, "Template execution error", http.StatusInternalServerError)
			return
		}
		return
	}

	if r.Method == "POST" {
		logrus.Info("Processing login POST request")

		// Parse form data
		err := r.ParseForm()
		if err != nil {
			logrus.WithError(err).Error("Error parsing login form")
			http.Redirect(w, r, "/login?error=Error al procesar el formulario", http.StatusSeeOther)
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == "" || password == "" {
			logrus.Warn("Username or password empty")
			http.Redirect(w, r, "/login?error=Usuario y contraseña son requeridos", http.StatusSeeOther)
			return
		}

		// Check "Remember Me" checkbox
		rememberMe := r.FormValue("remember_me") == "on"
		logrus.WithField("remember_me", rememberMe).Info("Login attempt")

		var user models.User
		err = h.DB.QueryRow("SELECT id, username, password, role, COALESCE(language, 'es') FROM users WHERE username=?", username).Scan(&user.ID, &user.Username, &user.Password, &user.Role, &user.Language)
		if err != nil {
			logrus.WithError(err).Warn("User not found")
			http.Redirect(w, r, "/login?error=Usuario o contraseña incorrectos", http.StatusSeeOther)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			logrus.Warn("Invalid password")
			http.Redirect(w, r, "/login?error=Usuario o contraseña incorrectos", http.StatusSeeOther)
			return
		}

		// Determine session duration based on "Remember Me"
		var sessionDuration time.Duration
		if rememberMe {
			sessionDuration = middleware.SessionDurationLong
			logrus.Info("Long session (30 days) created")
		} else {
			sessionDuration = middleware.SessionDurationShort
			logrus.Info("Short session (12 hours) created")
		}

		// Generate JWT
		claims := &middleware.Claims{
			Username:   user.Username,
			Role:       user.Role,
			RememberMe: rememberMe,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(sessionDuration)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				NotBefore: jwt.NewNumericDate(time.Now()),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(middleware.JWTSecret)
		if err != nil {
			logrus.WithError(err).Error("Token generation error")
			http.Error(w, "Token generation error", http.StatusInternalServerError)
			return
		}

		// Set cookie with appropriate duration
		cookieMaxAge := int(sessionDuration.Seconds())
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    tokenString,
			Path:     "/",
			HttpOnly: true,
			Secure:   false, // Set to true in production with HTTPS
			SameSite: http.SameSiteStrictMode,
			MaxAge:   cookieMaxAge,
		})

		logrus.WithFields(logrus.Fields{
			"username":      user.Username,
			"remember_me":   rememberMe,
			"cookie_maxage": cookieMaxAge / 3600,
		}).Info("Login successful")

		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
}

func (h *Handler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
	http.Redirect(w, r, "/login", http.StatusFound) // 302
}

func (h *Handler) HandleSettings(w http.ResponseWriter, r *http.Request) {
	// Get current user info from context
	username := r.Context().Value("username")
	if username == nil {
		logrus.Error("No username in context - redirecting to login")
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	usernameStr := username.(string)
	logrus.WithField("username", usernameStr).Info("Accessing settings page")

	// Get user data from database
	var user models.User
	err := h.DB.QueryRow("SELECT id, username, COALESCE(avatar, ''), COALESCE(language, ''), COALESCE(theme, '') FROM users WHERE username = ?", usernameStr).Scan(&user.ID, &user.Username, &user.Avatar, &user.Language, &user.Theme)
	if err != nil {
		logrus.WithError(err).WithField("username", usernameStr).Error("Database error in settings - user not found")
		// If user not found, redirect to login (token might be invalid)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	logrus.WithFields(logrus.Fields{
		"user":            user.Username,
		"avatar_length":   len(user.Avatar),
		"avatar_has_data": user.Avatar != "",
		"avatar_prefix": func() string {
			if len(user.Avatar) > 50 {
				return user.Avatar[:50] + "..."
			}
			return user.Avatar
		}(),
	}).Info("Settings loaded successfully for user")

	// Add URL to context for template
	data := struct {
		models.User
		URL interface{}
	}{
		User: user,
		URL:  r.URL,
	}

	tmpl, err := template.ParseFiles("web/templates/settings.html")
	if err != nil {
		logrus.WithError(err).Error("Template parsing error in settings")
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		logrus.WithError(err).Error("Template execution error in settings")
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	logrus.Info("Settings page rendered successfully")
}

func (h *Handler) HandleUpdateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get current user info from context
	username := r.Context().Value("username").(string)

	// Log incoming request for debugging
	contentType := r.Header.Get("Content-Type")
	logrus.WithFields(logrus.Fields{
		"username":     username,
		"content_type": contentType,
		"method":       r.Method,
	}).Info("HandleUpdateProfile request received")

	// Parse form based on content type
	var err error
	if strings.Contains(contentType, "multipart/form-data") {
		err = r.ParseMultipartForm(32 << 20) // 32MB max for file uploads
		if err != nil {
			logrus.WithError(err).Error("Multipart form parse error")
			http.Error(w, "Form parse error", http.StatusBadRequest)
			return
		}
	} else {
		err = r.ParseForm()
		if err != nil {
			logrus.WithError(err).Error("Regular form parse error")
			http.Error(w, "Form parse error", http.StatusBadRequest)
			return
		}
	}

	newUsername := strings.TrimSpace(r.FormValue("username"))
	newPassword := r.FormValue("password")
	passwordConfirm := r.FormValue("password_confirm")
	currentPassword := r.FormValue("current_password")
	newTheme := strings.TrimSpace(r.FormValue("theme"))
	newLanguage := strings.TrimSpace(r.FormValue("language"))

	logrus.WithFields(logrus.Fields{
		"theme":    newTheme,
		"language": newLanguage,
		"username": newUsername,
	}).Info("Form values received")

	// Check if password change is requested
	passwordChangeRequested := newPassword != ""

	// Validate current password only if changing password
	if passwordChangeRequested {
		if currentPassword == "" {
			http.Redirect(w, r, "/settings?error=current_password_required", http.StatusSeeOther)
			return
		}

		var hashedPassword string
		err = h.DB.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&hashedPassword)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(currentPassword))
		if err != nil {
			// Redirect back with error
			http.Redirect(w, r, "/settings?error=invalid_current_password", http.StatusSeeOther)
			return
		}

		// Check password confirmation
		if newPassword != passwordConfirm {
			http.Redirect(w, r, "/settings?error=password_mismatch", http.StatusSeeOther)
			return
		}
	}

	// Update username if provided and different
	if newUsername != "" && newUsername != username {
		// Check if username already exists
		var count int
		h.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ? AND username != ?", newUsername, username).Scan(&count)
		if count > 0 {
			http.Redirect(w, r, "/settings?error=username_exists", http.StatusSeeOther)
			return
		}

		_, err = h.DB.Exec("UPDATE users SET username = ? WHERE username = ?", newUsername, username)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		// Username changed - redirect to login since JWT contains old username
		http.Redirect(w, r, "/login?message=Usuario actualizado correctamente. Por favor, inicia sesión nuevamente.", http.StatusSeeOther)
		return
	}

	// Update password if provided
	if passwordChangeRequested {
		hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Password hash error", http.StatusInternalServerError)
			return
		}

		_, err = h.DB.Exec("UPDATE users SET password = ? WHERE username = ?", string(hashedNewPassword), username)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
	}

	// Handle avatar upload
	file, header, err := r.FormFile("avatar")
	if err == nil {
		defer file.Close()

		logrus.WithFields(logrus.Fields{
			"filename":     header.Filename,
			"size":         header.Size,
			"content_type": header.Header.Get("Content-Type"),
		}).Info("Processing avatar upload")

		// Read file content first to detect MIME type
		buffer := make([]byte, header.Size)
		_, err = file.Read(buffer)
		if err != nil {
			logrus.WithError(err).Error("File read error during avatar upload")
			http.Error(w, "File read error", http.StatusInternalServerError)
			return
		}

		// Detect MIME type from file content
		mimeType := detectMimeType(buffer)
		logrus.WithFields(logrus.Fields{
			"detected_mime": mimeType,
			"header_mime":   header.Header.Get("Content-Type"),
		}).Info("MIME type detection")

		// Validate file type
		if !strings.Contains(mimeType, "image/") {
			logrus.Error("Invalid file type for avatar upload")
			http.Redirect(w, r, "/settings?error=invalid_file_type", http.StatusSeeOther)
			return
		}

		// Save to database as base64 only (without data:image prefix)
		avatarData := base64.StdEncoding.EncodeToString(buffer)
		logrus.WithField("avatar_length", len(avatarData)).Info("Avatar data prepared for database")

		_, err = h.DB.Exec("UPDATE users SET avatar = ? WHERE username = ?", avatarData, username)
		if err != nil {
			logrus.WithError(err).Error("Database error saving avatar")
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		logrus.Info("Avatar saved successfully")
	} else if err != http.ErrMissingFile {
		logrus.WithError(err).Error("Error processing avatar file")
	}

	// Update theme and language preferences (with safe defaults)
	if newTheme == "" {
		newTheme = "dark"
	}
	if newLanguage == "" {
		newLanguage = "es"
	}

	logrus.WithFields(logrus.Fields{
		"theme":    newTheme,
		"language": newLanguage,
		"username": username,
	}).Info("About to UPDATE theme and language")

	result, err := h.DB.Exec("UPDATE users SET theme = ?, language = ? WHERE username = ?", newTheme, newLanguage, username)
	if err != nil {
		logrus.WithError(err).Error("Database error updating theme/language")
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	logrus.WithFields(logrus.Fields{
		"rows_affected": rowsAffected,
		"theme":         newTheme,
		"language":      newLanguage,
	}).Info("UPDATE executed successfully")

	// Set language cookie for i18n synchronization
	languageCookieExpiration := time.Now().Add(365 * 24 * time.Hour) // 1 year
	http.SetCookie(w, &http.Cookie{
		Name:     "currentLanguage",
		Value:    newLanguage,
		Expires:  languageCookieExpiration,
		Path:     "/",
		HttpOnly: false, // Need to be readable by JavaScript
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, "/settings?success=profile_updated&lang_updated=true", http.StatusSeeOther)
}
func (h *Handler) HandleAbout(w http.ResponseWriter, r *http.Request) {
	logrus.Info("About page requested")

	// Get username from context
	username := r.Context().Value("username")
	var user models.User

	if username != nil {
		// User is authenticated - get their preferences
		usernameStr := username.(string)
		err := h.DB.QueryRow("SELECT id, username, COALESCE(avatar, ''), COALESCE(language, 'es'), COALESCE(theme, 'dark') FROM users WHERE username = ?", usernameStr).Scan(&user.ID, &user.Username, &user.Avatar, &user.Language, &user.Theme)
		if err != nil {
			logrus.WithError(err).Warn("Could not get user data, using defaults")
			user.Language = "es"
		}
	} else {
		// User not authenticated - use default language
		user.Language = "es"
	}

	tmpl, err := template.ParseFiles("web/templates/about.html")
	if err != nil {
		logrus.WithError(err).Error("Template parsing error in about")
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, map[string]interface{}{
		"User":    user,
		"Version": Version,
	})
	if err != nil {
		logrus.WithError(err).Error("Template execution error in about")
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	logrus.Info("About page rendered successfully")
}

// HandleSetup shows the setup wizard for first-time installation
func (h *Handler) HandleSetup(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Setup wizard requested")

	// Check if setup is still needed
	var userCount int
	err := h.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount)
	if err != nil {
		logrus.WithError(err).Error("Error checking user count in setup")
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// If users already exist, redirect to login
	if userCount > 0 {
		logrus.Info("Users already exist, redirecting to login")
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// Serve setup wizard
	tmpl, err := template.ParseFiles("web/templates/setup.html")
	if err != nil {
		logrus.WithError(err).Error("Template parsing error in setup")
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"User": models.User{
			Language: "en", // Default language for setup page
		},
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		logrus.WithError(err).Error("Template execution error in setup")
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	logrus.Info("Setup wizard rendered successfully")
}

// HandleSetupSubmit processes the setup wizard form
func (h *Handler) HandleSetupSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	logrus.Info("Processing setup form submission")

	// Check if setup is still needed
	var userCount int
	err := h.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount)
	if err != nil {
		logrus.WithError(err).Error("Error checking user count in setup submit")
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if userCount > 0 {
		logrus.Warn("Setup attempted but users already exist")
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// Parse form with file upload
	err = r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		logrus.WithError(err).Error("Error parsing multipart form in setup")
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Get form values
	username := r.FormValue("username")
	password := r.FormValue("password")
	passwordConfirm := r.FormValue("password_confirm")
	theme := r.FormValue("theme")
	language := r.FormValue("language")

	logrus.WithFields(logrus.Fields{
		"username": username,
		"theme":    theme,
		"language": language,
	}).Info("Setup form values received")

	// Validate inputs
	if username == "" || password == "" {
		logrus.Error("Username or password missing in setup")
		http.Redirect(w, r, "/setup?error=required_fields", http.StatusSeeOther)
		return
	}

	if password != passwordConfirm {
		logrus.Error("Password mismatch in setup")
		http.Redirect(w, r, "/setup?error=password_mismatch", http.StatusSeeOther)
		return
	}

	if len(password) < 6 {
		logrus.Error("Password too short in setup")
		http.Redirect(w, r, "/setup?error=password_length", http.StatusSeeOther)
		return
	}

	// Set defaults if not provided
	if theme == "" {
		theme = "dark"
	}
	if language == "" {
		language = "en"
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logrus.WithError(err).Error("Error hashing password in setup")
		http.Error(w, "Error processing password", http.StatusInternalServerError)
		return
	}

	// Handle avatar upload (optional)
	var avatarData string
	file, header, err := r.FormFile("avatar")
	if err == nil {
		defer file.Close()
		logrus.WithFields(logrus.Fields{
			"filename": header.Filename,
			"size":     header.Size,
		}).Info("Processing avatar upload in setup")

		// Read file
		buffer := make([]byte, header.Size)
		_, err = file.Read(buffer)
		if err != nil {
			logrus.WithError(err).Error("Error reading avatar file in setup")
		} else {
			// Validate MIME type
			mimeType := detectMimeType(buffer)
			if strings.Contains(mimeType, "image/") {
				avatarData = base64.StdEncoding.EncodeToString(buffer)
				logrus.Info("Avatar processed successfully in setup")
			} else {
				logrus.Warn("Invalid avatar file type in setup")
			}
		}
	} else if err != http.ErrMissingFile {
		logrus.WithError(err).Warn("Error with avatar file in setup")
	}

	// Insert new user
	_, err = h.DB.Exec(
		"INSERT INTO users (username, password, avatar, role, language, theme) VALUES (?, ?, ?, ?, ?, ?)",
		username, string(hashedPassword), avatarData, "admin", language, theme,
	)
	if err != nil {
		logrus.WithError(err).Error("Error creating user in setup")
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	logrus.WithField("username", username).Info("First user created successfully via setup wizard")

	// Seed example bookmarks only on first setup
	var bmCount int
	if err := h.DB.QueryRow("SELECT COUNT(*) FROM bookmarks").Scan(&bmCount); err != nil {
		logrus.WithError(err).Warn("Error checking bookmarks count during setup seed")
	} else if bmCount == 0 {
		samples := []struct {
			title string
			url   string
			icon  string
		}{
			{"CodigoSH", "https://github.com/kiwinh0/CodigoSH", "git"},
			{"Selfh-st", "https://selfh.st/icons/", "selfh-st"},
			{"Helper Scripts", "https://community-scripts.github.io/ProxmoxVE/", "proxmox-helper-scripts"},
		}

		stmt, prepErr := h.DB.Prepare("INSERT INTO bookmarks (title, url, icon, position) VALUES (?, ?, ?, ?)")
		if prepErr != nil {
			logrus.WithError(prepErr).Warn("Could not prepare seed insert for bookmarks")
		} else {
			defer stmt.Close()
			for idx, s := range samples {
				if _, execErr := stmt.Exec(s.title, s.url, s.icon, idx); execErr != nil {
					logrus.WithError(execErr).Warn("Failed inserting sample bookmark")
				}
			}
			logrus.Info("Seeded sample bookmarks after setup")
		}
	}

	// Generate JWT token for automatic login
	expirationTime := time.Now().Add(12 * time.Hour) // 12 hour session
	claims := &middleware.Claims{
		Username: username,
		Role:     "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(middleware.JWTSecret)
	if err != nil {
		logrus.WithError(err).Error("Error generating JWT token in setup")
		http.Error(w, "Error generating session", http.StatusInternalServerError)
		return
	}

	// Set cookie with JWT token
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationTime,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	// Set language cookie for i18n
	languageCookieExpiration := time.Now().Add(365 * 24 * time.Hour) // 1 year
	http.SetCookie(w, &http.Cookie{
		Name:     "currentLanguage",
		Value:    language,
		Expires:  languageCookieExpiration,
		Path:     "/",
		HttpOnly: false, // Need to be readable by JavaScript
		SameSite: http.SameSiteLaxMode,
	})

	logrus.WithField("username", username).Info("Session cookie set for setup completion")

	// Redirect to dashboard
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// HandleExportData exports all user data and bookmarks as JSON
func (h *Handler) HandleExportData(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get username from context
	username := r.Context().Value("username")
	if username == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	usernameStr := username.(string)
	logrus.WithField("username", usernameStr).Info("Exporting data")

	// Get user data
	var user struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
		Avatar   string `json:"avatar"`
		Language string `json:"language"`
		Theme    string `json:"theme"`
	}

	err := h.DB.QueryRow("SELECT id, username, COALESCE(avatar, ''), COALESCE(language, 'en'), COALESCE(theme, 'dark') FROM users WHERE username = ?", usernameStr).
		Scan(&user.ID, &user.Username, &user.Avatar, &user.Language, &user.Theme)
	if err != nil {
		logrus.WithError(err).Error("Error querying user for export")
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Get all bookmarks
	rows, err := h.DB.Query("SELECT id, title, url, icon, position FROM bookmarks ORDER BY position ASC")
	if err != nil {
		logrus.WithError(err).Error("Error querying bookmarks for export")
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var bookmarks []struct {
		ID       int    `json:"id"`
		Title    string `json:"title"`
		URL      string `json:"url"`
		Icon     string `json:"icon"`
		Position int    `json:"position"`
	}

	for rows.Next() {
		var b struct {
			ID       int    `json:"id"`
			Title    string `json:"title"`
			URL      string `json:"url"`
			Icon     string `json:"icon"`
			Position int    `json:"position"`
		}
		if err := rows.Scan(&b.ID, &b.Title, &b.URL, &b.Icon, &b.Position); err != nil {
			logrus.WithError(err).Error("Error scanning bookmark for export")
			continue
		}
		bookmarks = append(bookmarks, b)
	}

	// Create export structure
	exportData := struct {
		Version    string `json:"version"`
		ExportedAt string `json:"exported_at"`
		User       struct {
			Username string `json:"username"`
			Avatar   string `json:"avatar"`
			Language string `json:"language"`
			Theme    string `json:"theme"`
		} `json:"user"`
		Bookmarks []struct {
			ID       int    `json:"id"`
			Title    string `json:"title"`
			URL      string `json:"url"`
			Icon     string `json:"icon"`
			Position int    `json:"position"`
		} `json:"bookmarks"`
	}{
		Version:    Version,
		ExportedAt: time.Now().Format(time.RFC3339),
		User: struct {
			Username string `json:"username"`
			Avatar   string `json:"avatar"`
			Language string `json:"language"`
			Theme    string `json:"theme"`
		}{
			Username: user.Username,
			Avatar:   user.Avatar,
			Language: user.Language,
			Theme:    user.Theme,
		},
		Bookmarks: bookmarks,
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename=codigosh-backup.json")

	if err := json.NewEncoder(w).Encode(exportData); err != nil {
		logrus.WithError(err).Error("Error encoding export data")
		http.Error(w, "Error encoding data", http.StatusInternalServerError)
		return
	}

	logrus.WithFields(logrus.Fields{
		"username":        usernameStr,
		"bookmarks_count": len(bookmarks),
	}).Info("Data exported successfully")
}

// HandleImportData imports user data and bookmarks from JSON
func (h *Handler) HandleImportData(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get username from context
	username := r.Context().Value("username")
	if username == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	usernameStr := username.(string)
	logrus.WithField("username", usernameStr).Info("Importing data")

	// Parse multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB max
		logrus.WithError(err).Error("Error parsing multipart form")
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		logrus.WithError(err).Error("Error getting file from form")
		http.Error(w, "Error reading file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Parse JSON
	var importData struct {
		Version string `json:"version"`
		User    struct {
			Avatar   string `json:"avatar"`
			Language string `json:"language"`
			Theme    string `json:"theme"`
		} `json:"user"`
		Bookmarks []struct {
			Title    string `json:"title"`
			URL      string `json:"url"`
			Icon     string `json:"icon"`
			Position int    `json:"position"`
		} `json:"bookmarks"`
	}

	if err := json.NewDecoder(file).Decode(&importData); err != nil {
		logrus.WithError(err).Error("Error decoding JSON")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON format"})
		return
	}

	// Start transaction
	tx, err := h.DB.Begin()
	if err != nil {
		logrus.WithError(err).Error("Error starting transaction")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Database error"})
		return
	}
	defer tx.Rollback()

	// Update user preferences
	_, err = tx.Exec("UPDATE users SET avatar = ?, language = ?, theme = ? WHERE username = ?",
		importData.User.Avatar, importData.User.Language, importData.User.Theme, usernameStr)
	if err != nil {
		logrus.WithError(err).Error("Error updating user preferences")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error updating user preferences"})
		return
	}

	// Delete existing bookmarks
	_, err = tx.Exec("DELETE FROM bookmarks")
	if err != nil {
		logrus.WithError(err).Error("Error deleting existing bookmarks")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error clearing bookmarks"})
		return
	}

	// Insert imported bookmarks
	stmt, err := tx.Prepare("INSERT INTO bookmarks (title, url, icon, position) VALUES (?, ?, ?, ?)")
	if err != nil {
		logrus.WithError(err).Error("Error preparing bookmark insert")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error preparing import"})
		return
	}
	defer stmt.Close()

	for _, b := range importData.Bookmarks {
		if _, err := stmt.Exec(b.Title, b.URL, b.Icon, b.Position); err != nil {
			logrus.WithError(err).WithField("bookmark", b.Title).Error("Error inserting bookmark")
			continue
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		logrus.WithError(err).Error("Error committing transaction")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error saving data"})
		return
	}

	logrus.WithFields(logrus.Fields{
		"username":           usernameStr,
		"bookmarks_imported": len(importData.Bookmarks),
	}).Info("Data imported successfully")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("Datos importados: %d marcadores", len(importData.Bookmarks)),
	})
}
