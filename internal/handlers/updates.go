package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
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
	cacheExpiration = 1 * time.Hour
)

// InvalidateUpdateCache fuerza la próxima verificación a ir contra GitHub
func InvalidateUpdateCache() {
	updateCacheMu.Lock()
	defer updateCacheMu.Unlock()
	updateCache = nil
	lastCheckTime = time.Time{}
}

// CheckForUpdates verifica si hay una nueva versión en GitHub
func CheckForUpdates() *UpdateInfo {
	updateCacheMu.RLock()
	if updateCache != nil && time.Since(lastCheckTime) < cacheExpiration {
		info := updateCache
		updateCacheMu.RUnlock()
		logrus.WithFields(logrus.Fields{
			"source":            "cache",
			"current_version":   Version,
			"available_version": info.AvailableVersion,
			"update_available":  info.UpdateAvailable,
		}).Debug("Update check served from cache")
		return info
	}
	updateCacheMu.RUnlock()

	// Obtener la última release (incluye prereleases, excluye drafts) para no saltar betas
	resp, err := http.Get("https://api.github.com/repos/kiwinh0/CodigoSH/releases?per_page=1")
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

	var releases []struct {
		TagName string `json:"tag_name"`
		Body    string `json:"body"`
		Date    string `json:"published_at"`
	}

	if err := json.Unmarshal(body, &releases); err != nil {
		logrus.WithError(err).Warn("Could not parse GitHub response")
		return &UpdateInfo{UpdateAvailable: false}
	}

	if len(releases) == 0 {
		logrus.Warn("No releases found in GitHub response")
		return &UpdateInfo{UpdateAvailable: false}
	}

	release := releases[0]

	// Extraer versión del tag (ej: v0.2.0 -> 0.2.0)
	latestVersion := strings.TrimPrefix(release.TagName, "v")

	// Comparar versiones
	cmpResult := compareVersions(Version, latestVersion)
	updateAvailable := cmpResult < 0
	logrus.WithFields(logrus.Fields{
		"source":           "github",
		"current_version":  Version,
		"latest_version":   latestVersion,
		"compare_result":   cmpResult,
		"update_available": updateAvailable,
	}).Info("Update check evaluated")

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

// compareVersions compara dos versiones semánticas (1.0.0 vs 0.9.0)
// Retorna: -1 si v1 < v2, 0 si son iguales, 1 si v1 > v2
func compareVersions(v1, v2 string) int {
	// Limpiar sufijos (ej: -Beta, -alpha, etc)
	v1Clean := strings.Split(v1, "-")[0]
	v2Clean := strings.Split(v2, "-")[0]

	// Dividir en partes numéricas
	v1Parts := strings.Split(v1Clean, ".")
	v2Parts := strings.Split(v2Clean, ".")

	// Comparar cada parte
	maxLen := len(v1Parts)
	if len(v2Parts) > maxLen {
		maxLen = len(v2Parts)
	}

	for i := 0; i < maxLen; i++ {
		var num1, num2 int

		if i < len(v1Parts) {
			num1, _ = strconv.Atoi(v1Parts[i])
		}
		if i < len(v2Parts) {
			num2, _ = strconv.Atoi(v2Parts[i])
		}

		if num1 < num2 {
			return -1
		}
		if num1 > num2 {
			return 1
		}
	}

	return 0
}

// parseChanges extrae los cambios principales del body del release
func parseChanges(body string) []string {
	if body == "" {
		return []string{"Varias mejoras internas"}
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
		changes = []string{"Varias mejoras internas"}
	}

	return changes
}

// HandleCheckUpdates es el handler para verificar actualizaciones desde el frontend
func (h *Handler) HandleCheckUpdates(w http.ResponseWriter, r *http.Request) {
	// Permitir forzar invalidación de caché con ?force=true
	if r.URL.Query().Get("force") == "true" {
		InvalidateUpdateCache()
	}

	info := CheckForUpdates()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

// HandlePerformUpdate descarga e instala el binario pre-compilado desde GitHub
func (h *Handler) HandlePerformUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Verificar que hay una actualización disponible
	updateInfo := CheckForUpdates()
	if !updateInfo.UpdateAvailable {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "No hay actualizaciones disponibles",
		})
		return
	}

	logrus.Info("Iniciando descarga de actualización...")

	// Obtener ruta actual del ejecutable
	currentExecutable, err := os.Executable()
	if err != nil {
		logrus.WithError(err).Error("Error obteniendo ruta ejecutable")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Error en la actualización",
		})
		return
	}

	// Crear script de actualización que se ejecutará de forma independiente
	updateScript := "/tmp/update_codigosh.sh"
	releaseVersion := "v" + updateInfo.AvailableVersion

	scriptContent := `#!/bin/bash

# Redirigir TODO el output al log (stderr y stdout)
exec > >(tee -a /tmp/codigosh_update.log) 2>&1

# Función para logging con timestamp
log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $*"
}

log "=== Iniciando actualización de CodigoSH a ` + releaseVersion + ` ==="

RELEASE_VERSION="` + releaseVersion + `"
BINARY_PATH="` + currentExecutable + `"
TEMP_DIR="/tmp/codigosh-update-$$"

mkdir -p "$TEMP_DIR"
cd "$TEMP_DIR" || exit 1

# Descargar binario pre-compilado desde GitHub Release
log "Descargando binario desde GitHub Release $RELEASE_VERSION..."
DOWNLOAD_URL="https://github.com/kiwinh0/CodigoSH/releases/download/$RELEASE_VERSION/codigosH"

if ! curl -L -f -o codigosH "$DOWNLOAD_URL"; then
    log "ERROR: No se pudo descargar el binario desde $DOWNLOAD_URL"
    log "Esto puede indicar que el release no tiene binario adjunto"
    exit 1
fi

log "Binario descargado correctamente ($(du -h codigosH | cut -f1))"

# Verificar que el binario descargado es válido
if [ ! -f codigosH ] || [ ! -s codigosH ]; then
    log "ERROR: El archivo descargado está vacío o no existe"
    exit 1
fi

# Hacer ejecutable
chmod +x codigosH
log "Permisos de ejecución aplicados"

# Backup del binario actual
log "Creando backup del binario actual..."
cp "$BINARY_PATH" "$BINARY_PATH.backup" || {
    log "ERROR: No se pudo crear backup"
    exit 1
}
log "Backup creado: $BINARY_PATH.backup"

# Reemplazar binario
log "Instalando nuevo binario..."
cp codigosH "$BINARY_PATH" || {
    log "ERROR: No se pudo copiar el nuevo binario, restaurando backup..."
    cp "$BINARY_PATH.backup" "$BINARY_PATH"
    exit 1
}
chmod +x "$BINARY_PATH"

# Verificar instalación
if [ ! -f "$BINARY_PATH" ] || [ ! -x "$BINARY_PATH" ]; then
    log "ERROR: La instalación falló, restaurando backup..."
    cp "$BINARY_PATH.backup" "$BINARY_PATH"
    exit 1
fi

log "Binario instalado correctamente"

# Limpiar archivos temporales
log "Limpiando archivos temporales..."
cd /tmp
rm -rf "$TEMP_DIR"

log "=== Actualización completada exitosamente ==="
log "Esperando 2 segundos antes de reiniciar el servicio..."
sleep 2

# Reiniciar servicio
log "Reiniciando servicio..."
if command -v systemctl &> /dev/null && systemctl is-active --quiet codigosH; then
    log "Reiniciando con systemctl..."
    systemctl restart codigosH && log "Servicio reiniciado con systemctl"
elif command -v service &> /dev/null; then
    log "Reiniciando con service..."
    service codigosH restart && log "Servicio reiniciado con service"
else
    log "systemctl/service no disponible. Intentando reinicio manual..."
    # Matar proceso actual
    CURRENT_PID=$(pgrep -f "codigosH" | head -1)
    if [ -n "$CURRENT_PID" ]; then
        log "Matando proceso actual (PID: $CURRENT_PID)..."
        kill -9 "$CURRENT_PID" 2>/dev/null || true
        sleep 1
    fi
    # Reiniciar binario directamente en segundo plano
    log "Iniciando nuevo binario en background..."
    nohup "$BINARY_PATH" > /dev/null 2>&1 &
    NEW_PID=$!
    log "Binario reiniciado directamente (PID: $NEW_PID)"
fi

log "=== Proceso de actualización finalizado ==="

# Eliminar el script al final
rm -f /tmp/update_codigosh.sh
`

	if err := os.WriteFile(updateScript, []byte(scriptContent), 0755); err != nil {
		logrus.WithError(err).Error("Error creando script de actualización")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Error preparando actualización",
		})
		return
	}

	logrus.WithField("script_path", updateScript).Info("Script de actualización creado")

	// Ejecutar el script usando bash directamente con nohup
	// &  al final lo pone en background, pero queremos que persista después del cierre del servidor
	cmd := exec.Command("bash", "-c", fmt.Sprintf("nohup bash %s > /dev/null 2>&1 &", updateScript))
	
	if err := cmd.Run(); err != nil {
		logrus.WithError(err).Error("Error iniciando script de actualización")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Error iniciando actualización",
		})
		return
	}

	logrus.Info("Script de actualización lanzado exitosamente en background")

	// Invalidar caché para próxima verificación
	InvalidateUpdateCache()

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Actualización iniciada. El servicio se reiniciará en breve.",
	})
}
