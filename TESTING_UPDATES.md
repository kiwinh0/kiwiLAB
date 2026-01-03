# Prueba del Sistema de Actualizaciones Automáticas

Este documento describe cómo probar el nuevo sistema de actualizaciones que compila desde fuente.

## 📋 Descripción Técnica

### Flujo de Actualización

```
1. Usuario hace clic en "Actualizar"
   ↓
2. Frontend envía POST a /perform-update
   ↓
3. Backend inicia script en background
   ├─ Crea /tmp/codigosh-update/
   ├─ Git clone desde main branch
   ├─ go mod download
   ├─ go build
   ├─ Backup binario actual
   ├─ Reemplaza binario
   ├─ systemctl restart codigosH
   └─ Logs en /tmp/codigosh_update.log
   ↓
4. Frontend recibe respuesta de éxito
   ↓
5. Usuario espera 30-60 segundos
   ↓
6. Servicio se reinicia y refresca la página
```

### Archivos Modificados

- **internal/handlers/updates.go**: Handler `HandlePerformUpdate` - Nuevo flujo de compilación
- **web/templates/dashboard.html**: Función `performUpdate()` - Mensajes mejorados
- **README.md**: Sección nueva sobre actualizaciones

## 🧪 Cómo Probar Localmente (macOS/Linux)

### Prerequisitos

```bash
# Go 1.24+
go version

# Git
git --version

# Compilar la versión local
cd /Users/kiwinho/Proyectos/CodigoSH
go build -o bin/codigosH ./cmd/codigosH/main.go
```

### Pasos de Prueba

1. **Crear rama de prueba para simular nueva versión:**
```bash
cd /Users/kiwinho/Proyectos/CodigoSH
git tag v0.1.4-Beta
git push origin v0.1.4-Beta
```

2. **Ejecutar la aplicación localmente:**
```bash
./bin/codigosH
# O: go run ./cmd/codigosH/main.go
```

3. **Acceder a http://localhost:8080**

4. **Verificar que se detecta actualización:**
   - Debería aparecer badge rojo en el avatar
   - Menú usuario → "Actualización disponible" en rojo

5. **Hacer clic en "Actualizar":**
   - Botón muestra "Compilando..."
   - Espera 30-60 segundos
   - Botón muestra "✓ Actualización iniciada"

6. **Verificar logs de actualización:**
```bash
tail -f /tmp/codigosh_update.log
```

Debería ver algo como:
```
[2024-01-15 10:30:45] Iniciando actualización de CodigoSH
[2024-01-15 10:30:46] Clonando repositorio desde GitHub...
[2024-01-15 10:30:52] Descargando dependencias...
[2024-01-15 10:31:15] Compilando binario...
[2024-01-15 10:31:32] Backup del binario actual...
[2024-01-15 10:31:32] Instalando nuevo binario...
[2024-01-15 10:31:32] Limpiando archivos temporales...
[2024-01-15 10:31:33] Reiniciando servicio codigosH...
[2024-01-15 10:31:33] Actualización completada exitosamente
```

## 🚀 Prueba en Producción (Debian/Ubuntu)

### Setup Inicial

```bash
# Instalar CodigoSH
curl -sSL "https://raw.githubusercontent.com/kiwinh0/CodigoSH/main/scripts/install.sh" | sudo bash

# Verificar que está corriendo
sudo systemctl status codigosH
```

### Acceso

```
URL: http://IP_DEL_SERVIDOR:8080
```

### Prueba de Actualización

1. **Verificar versión actual:**
   - Click en usuario → "Acerca de"
   - Mostrar version actual (ej: v0.1.3-Beta)

2. **Crear nueva tag en GitHub:**
```bash
git tag v0.1.4-Beta
git push origin v0.1.4-Beta
```

3. **En la web, acceder a actualización:**
   - Debería aparecer "Actualización disponible"
   - Click en actualizar
   - Esperar a que compile

4. **Verificar actualización exitosa:**
```bash
# Ver logs del sistema
sudo journalctl -u codigosH -n 50 -f

# Ver logs de actualización
cat /tmp/codigosh_update.log

# Verificar nuevo binario
ls -lah /opt/CodigoSH/codigosH
ls -lah /opt/CodigoSH/codigosH.backup
```

5. **Acceder a CodigoSH de nuevo:**
   - Debería estar en nueva versión
   - Click en usuario → "Acerca de" mostrar nueva versión

## ⚠️ Solución de Problemas

### Error: "No hay actualizaciones disponibles"
- Verificar que hay nueva tag en GitHub (`git tag v0.1.4-Beta`)
- Verificar que versión en VERSION file es más antigua que tag

### Error: "Error en la actualización"
- Revisar `/tmp/codigosh_update.log`
- Verificar que Git está instalado: `which git`
- Verificar que Go está instalado: `go version`

### Script no ejecuta
- Verificar permisos: `chmod +x /tmp/update_codigosh_build.sh`
- Revisar logs: `cat /tmp/codigosh_update.log`

### Compilación muy lenta
- Normal la primera vez (descarga dependencias)
- Posteriores son más rápidas (cache de go mod)
- Timeout típico: 60 segundos en servidor 2-core

### Binario antiguo después de actualizar
- Revisar que binario fue reemplazado:
```bash
file /opt/CodigoSH/codigosH
/opt/CodigoSH/codigosH: ELF 64-bit LSB executable...
```
- Revisar fecha de modificación:
```bash
ls -la /opt/CodigoSH/codigosH
# Debe tener fecha reciente
```

## 📊 Checklist de Verificación

- [ ] Aplicación inicia correctamente
- [ ] Badge de actualización aparece cuando hay versión nueva
- [ ] Click en "Actualización disponible" abre modal
- [ ] Modal muestra versión actual vs nueva
- [ ] Modal muestra cambios/changelog
- [ ] Click en "Actualizar" inicia proceso
- [ ] Botón cambia a "Compilando..."
- [ ] Espera 30-60 segundos
- [ ] Botón cambia a "✓ Actualización iniciada"
- [ ] Logs aparecen en /tmp/codigosh_update.log
- [ ] Servicio se reinicia automáticamente
- [ ] Nueva versión visible después de recargar
- [ ] Backup creado en binario.backup
- [ ] No hay errores en systemd journal

## 🔍 Monitoreo Durante Actualización

```bash
# Terminal 1: Ver logs del sistema
sudo journalctl -u codigosH -f

# Terminal 2: Ver logs de actualización
tail -f /tmp/codigosh_update.log

# Terminal 3: Monitoreo de archivos
watch -n 1 'ls -lah /opt/CodigoSH/ | grep -E "codigosh|backup"'

# Terminal 4: Monitoreo de proceso
watch -n 1 'ps aux | grep -E "codigosH|go"'
```

## 📈 Mejoras Futuras

- [ ] Rollback automático si falla la compilación
- [ ] Notificaciones vía webhook después de actualización
- [ ] Actualizaciones programadas en horario específico
- [ ] Descarga precompilada en GitHub Actions
- [ ] Delta updates (solo cambios incrementales)
