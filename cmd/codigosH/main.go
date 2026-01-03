package main

import (
	"net/http"

	"github.com/kiwinh0/CodigoSH/internal/config"
	"github.com/kiwinh0/CodigoSH/internal/db"
	"github.com/kiwinh0/CodigoSH/internal/handlers"
	"github.com/kiwinh0/CodigoSH/internal/middleware"
	"github.com/sirupsen/logrus"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			logrus.Fatal("Panic in main: ", r)
		}
	}()

	logrus.Info("🚀 Iniciando CodigoSH v0.1.9-Beta...")

	// Cargar configuración
	cfg, err := config.LoadConfig()
	if err != nil {
		logrus.Fatalf("Error cargando configuración: %v", err)
	}

	// Configurar logging
	level, err := logrus.ParseLevel(cfg.Logging.Level)
	if err != nil {
		logrus.Warn("Nivel de log inválido, usando 'info'")
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)

	// Inicializar DB
	db.InitDB(cfg.Database.Path)
	logrus.Info("✅ Base de datos inicializada")

	// Crear handlers
	h := handlers.NewHandler(db.DB)
	mux := http.NewServeMux()

	// Archivos estáticos (públicos)
	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", middleware.LogStaticFiles(fs)))

	// Rutas públicas
	mux.HandleFunc("/login", h.HandleLogin)
	mux.HandleFunc("/setup", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.HandleSetup(w, r)
		} else if r.Method == http.MethodPost {
			h.HandleSetupSubmit(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Aplicar middleware a rutas protegidas
	authHandler := middleware.AuthMiddleware(logrus.StandardLogger())
	loggedHandler := middleware.LoggingMiddleware(logrus.StandardLogger())

	// Rutas protegidas
	mux.Handle("/dashboard", authHandler(loggedHandler(http.HandlerFunc(h.HandleIndex))))
	mux.Handle("/settings", authHandler(loggedHandler(http.HandlerFunc(h.HandleSettings))))
	mux.Handle("/about", authHandler(loggedHandler(http.HandlerFunc(h.HandleAbout))))
	mux.Handle("/add-bookmark", authHandler(loggedHandler(http.HandlerFunc(h.HandleAdd))))
	mux.Handle("/edit-bookmark", authHandler(loggedHandler(http.HandlerFunc(h.HandleEdit))))
	mux.Handle("/delete-bookmark", authHandler(loggedHandler(http.HandlerFunc(h.HandleDelete))))
	mux.Handle("/reorder-bookmarks", authHandler(loggedHandler(http.HandlerFunc(h.HandleReorder))))
	mux.Handle("/logout", authHandler(loggedHandler(http.HandlerFunc(h.HandleLogout))))
	mux.Handle("/update-profile", authHandler(loggedHandler(http.HandlerFunc(h.HandleUpdateProfile))))
	mux.Handle("/export-data", authHandler(loggedHandler(http.HandlerFunc(h.HandleExportData))))
	mux.Handle("/import-data", authHandler(loggedHandler(http.HandlerFunc(h.HandleImportData))))
	mux.Handle("/check-updates", authHandler(loggedHandler(http.HandlerFunc(h.HandleCheckUpdates))))
	mux.Handle("/perform-update", authHandler(loggedHandler(http.HandlerFunc(h.HandlePerformUpdate))))

	// Ruta raíz "/" debe ir AL FINAL como catch-all
	mux.HandleFunc("/", h.HandleRoot)

	logrus.Info("✅ Rutas configuradas")
	logrus.Infof("🚀 CodigoSH activo en http://%s:%s", cfg.Server.Host, cfg.Server.Port)

	// Iniciar servidor
	if err = http.ListenAndServe(cfg.Server.Host+":"+cfg.Server.Port, mux); err != nil {
		logrus.Fatalf("Error al iniciar servidor: %v", err)
	}
}
