# Sistema de Actualización Automática

## 🚀 Descripción

CodigoSH incluye un sistema de actualización automática que descarga e instala nuevas versiones desde GitHub de forma rápida y confiable.

## ⚡ Características

- **Ultra-rápido**: 5-10 segundos (6x más rápido que antes)
- **Binarios pre-compilados**: No requiere compilar Go en producción
- **Automático**: Detección, descarga, instalación y reinicio sin intervención
- **Seguro**: Backup automático y rollback en caso de error
- **Logs detallados**: Seguimiento completo del proceso

## 🔧 Cómo Funciona

### Para Usuarios

1. Abre la aplicación y ve a **Dashboard**
2. Si hay una actualización disponible, verás un **badge rojo** en el menú
3. Haz clic en **"Actualización"** en el menú de usuario
4. Se abre un modal mostrando la nueva versión disponible
5. Haz clic en **"Actualizar Ahora"**
6. Espera 5-10 segundos mientras se descarga e instala
7. La página se recarga automáticamente con la nueva versión 🎉

### Desde la Página "Acerca de"

1. Ve a **Acerca de** (About)
2. Si hay actualización, verás un **banner verde pulsante**
3. Haz clic en el banner
4. Serás redirigido al dashboard con el modal abierto

## 🔍 Proceso Técnico

### 1. Detección (Automática)
```
- Verifica GitHub cada hora
- Compara versión actual vs última release
- Muestra notificación si hay actualización disponible
```

### 2. Descarga (2-3 segundos)
```bash
curl -L -f -o codigosH \
  "https://github.com/kiwinh0/CodigoSH/releases/download/v0.2.4-Beta/codigosH"
```

### 3. Instalación (1 segundo)
```bash
# Backup del binario actual
cp /path/to/codigosH /path/to/codigosH.backup

# Instalar nueva versión
cp /tmp/codigosH /path/to/codigosH
chmod +x /path/to/codigosH
```

### 4. Reinicio (2-3 segundos)
```bash
# Intenta systemctl primero
systemctl restart codigosH

# Fallback si systemctl no está disponible
nohup /path/to/codigosH &
```

## 📝 Logs

Todos los pasos se registran en:
```
/tmp/codigosh_update.log
```

Ejemplo de log:
```
[2026-01-03 19:33:21] === Iniciando actualización de CodigoSH a v0.2.4-Beta ===
[2026-01-03 19:33:21] Descargando binario desde GitHub Release v0.2.4-Beta...
[2026-01-03 19:33:23] Binario descargado correctamente (11M)
[2026-01-03 19:33:23] Permisos de ejecución aplicados
[2026-01-03 19:33:23] Creando backup del binario actual...
[2026-01-03 19:33:23] Backup creado: /path/to/codigosH.backup
[2026-01-03 19:33:23] Instalando nuevo binario...
[2026-01-03 19:33:24] Binario instalado correctamente
[2026-01-03 19:33:24] Limpiando archivos temporales...
[2026-01-03 19:33:24] === Actualización completada exitosamente ===
[2026-01-03 19:33:26] Servicio reiniciado con systemctl
[2026-01-03 19:33:26] === Proceso de actualización finalizado ===
```

## 🔒 Seguridad

### Verificaciones
- ✅ Verifica que el archivo descargado existe y no está vacío
- ✅ Valida que el binario tiene permisos de ejecución
- ✅ Comprueba la instalación antes de eliminar el backup

### Backup Automático
- Cada actualización crea un backup: `codigosH.backup`
- Si la instalación falla, se restaura automáticamente
- El backup se mantiene hasta la próxima actualización exitosa

### Rollback Manual
Si algo sale mal, puedes restaurar manualmente:
```bash
cp /path/to/codigosH.backup /path/to/codigosH
systemctl restart codigosH
```

## 🛠️ Desarrollo

### Crear una Nueva Release con Binario

1. **Actualizar versión** en 3 archivos:
   - `VERSION`
   - `internal/handlers/handlers.go` (const Version)
   - `scripts/diagnostico.sh`

2. **Commit y tag**:
   ```bash
   git add -A
   git commit -m "chore: bump version to v0.2.5-Beta"
   git tag v0.2.5-Beta
   git push && git push --tags
   ```

3. **GitHub Actions automático**:
   - El workflow `.github/workflows/release.yml` se activa
   - Compila el binario con Go 1.24
   - Lo sube automáticamente al release
   - ¡Listo! 🎉

### Workflow de GitHub Actions
```yaml
# .github/workflows/release.yml
- Trigger: Push de tag `v*`
- Compila: CGO_ENABLED=1 go build -ldflags="-s -w"
- Sube: Adjunta binario al release
```

## 📊 Comparación: Antes vs Ahora

| Aspecto | Antes (Compilación) | Ahora (Binario) |
|---------|---------------------|-----------------|
| **Tiempo** | 30-60 segundos | 5-10 segundos |
| **Recursos** | Alto (compilación) | Bajo (solo descarga) |
| **Dependencias** | Go toolchain | Ninguna |
| **Confiabilidad** | Media (puede fallar) | Alta (pre-testeado) |
| **Tamaño descarga** | ~50MB (repo) | ~11MB (binario) |

## ❓ Preguntas Frecuentes

### ¿Qué pasa si falla la actualización?
El sistema hace backup automático y puede restaurarse manualmente. Revisa los logs en `/tmp/codigosh_update.log`.

### ¿Puedo desactivar las actualizaciones automáticas?
La detección es automática, pero la instalación requiere confirmación del usuario.

### ¿Funciona en desarrollo?
Sí, pero se recomienda usar `make build` localmente en vez de actualizar desde GitHub.

### ¿Puedo actualizar manualmente?
Sí, puedes descargar el binario desde el release y reemplazarlo manualmente:
```bash
wget https://github.com/kiwinh0/CodigoSH/releases/download/v0.2.4-Beta/codigosH
chmod +x codigosH
sudo systemctl stop codigosH
sudo cp codigosH /path/to/codigosH
sudo systemctl start codigosH
```

## 🎯 Próximos Pasos

- [ ] Soporte para múltiples arquitecturas (ARM, macOS)
- [ ] Verificación de checksums SHA256
- [ ] Opción de actualización programada (cron)
- [ ] Notificaciones por email/Telegram

## 📞 Soporte

Si tienes problemas con el sistema de actualización:
1. Revisa los logs: `cat /tmp/codigosh_update.log`
2. Verifica el backup: `ls -lh /path/to/codigosH.backup`
3. Reporta el issue en GitHub con los logs adjuntos
