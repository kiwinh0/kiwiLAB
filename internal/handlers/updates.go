package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// UpdateInfo contiene información de la actualización
type UpdateInfo struct {
	AvailableVersion string   `json:"available_version"`
	ReleaseDate      string   `json:"release_date"`
	Changes          []string `json:"changes"`
	UpdateAvailable  bool     `json:"update_available"`
}

// UpdateChecker gestiona la verificación de actualizaciones
var (
	updateCache     *UpdateInfo
	updateCacheMu   sync.RWMutex
	lastCheckTime   time.Time
	cacheExpiration = 24 * time.Hour
)

// CheckForUpdates verifica si hay una nueva versión en GitHub
func CheckForUpdates() *UpdateInfo {
	updateCacheMu.RLock()
	if updateCache != nil && time.Since(lastCheckTime) < cacheExpiration {
		defer updateCacheMu.RUnlock()
		return updateCache
	}
	updateCacheMu.RUnlock()

	// Obtener información de la última versión desde GitHub
	resp, err := http.Get("https://api.github.com/repos/kiwinh0/CodigoSH/releases/latest")
	if err != nil {
		logrus.WithError(err).Warn("Could not check for updates")
		return &UpdateInfo{UpdateAvailable: false}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &UpdateInfo{UpdateAvailable: false}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.WithError(err).Warn("Could not read GitHub response")
		return &UpdateInfo{UpdateAvailable: false}
	}

	var release struct {
		TagName string `json:"tag_name"`
		Body    string `json:"body"`
		Date    string `json:"published_at"`
	}

	if err := json.Unmarshal(body, &release); err != nil {
		logrus.WithError(err).Warn("Could not parse GitHub response")
		return &UpdateInfo{UpdateAvailable: false}
	}

	// Extraer versión del tag (ej: v0.2.0 -> 0.2.0)
	latestVersion := strings.TrimPrefix(release.TagName, "v")

	// Comparar versiones
	updateAvailable := compareVersions(Version, latestVersion) < 0

	// Parsear cambios desde el body
	changes := parseChanges(release.Body)

	info := &UpdateInfo{
		AvailableVersion: latestVersion,
		ReleaseDate:      release.Date,
		Changes:          changes,
		UpdateAvailable:  updateAvailable,
	}

	// Actualizar caché
	updateCacheMu.Lock()
	updateCache = info
	lastCheckTime = time.Now()
	updateCacheMu.Unlock()

	return info
}

// compareVersions compara dos versiones (1.0.0 vs 0.9.0)
// Retorna: -1 si v1 < v2, 0 si son iguales, 1 si v1 > v2
func compareVersions(v1, v2 string) int {
	// Simplificar: solo comparación de strings lexicográfica para esta beta
	if v1 == v2 {
		return 0
	}
	if v1 < v2 {
		return -1
	}
	return 1
}

// parseChanges extrae los cambios principales del body del release
func parseChanges(body string) []string {
	if body == "" {
		return []string{"Nueva actualización disponible"}
	}

	lines := strings.Split(body, "\n")
	var changes []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "- ") {
			change := strings.TrimPrefix(line, "- ")
			change = strings.TrimPrefix(change, "* ")
			if change != "" && len(changes) < 5 {
				changes = append(changes, change)
			}
		}
	}

	if len(changes) == 0 {
		changes = []string{"Nueva actualización disponible"}
	}

	return changes
}

// HandleCheckUpdates es el handler para verificar actualizaciones desde el frontend
func (h *Handler) HandleCheckUpdates(w http.ResponseWriter, r *http.Request) {
	info := CheckForUpdates()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}
