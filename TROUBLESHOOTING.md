# 🆘 Guía de Troubleshooting - CodigoSH

## Problema: El servicio no inicia después de la instalación

### Diagnóstico rápido

```bash
# Ver estado del servicio
sudo systemctl status codigosH

# Ver últimos logs
sudo journalctl -u codigosH -n 50 --no-pager

# Ver logs en tiempo real
sudo journalctl -u codigosH -f
```

### Soluciones por error

#### ❌ "command not found: go"

**Causa:** Go no está instalado o no está en PATH

**Solución:**
```bash
# Verificar instalación
go version

# Si no funciona, instalar Go 1.21+
sudo apt update
sudo apt install -y golang-go

# Verificar nuevamente
go version
```

---

#### ❌ "no such file or directory" para el binario

**Causa:** La compilación falló o el binario no se generó

**Solución:**
```bash
# Ir al directorio
cd /opt/CodigoSH

# Compilar manualmente
export CGO_ENABLED=1
go mod download
go build -o codigosH ./cmd/codigosH

# Verificar que existe
ls -la codigosH

# Si falta sqlite, instalar:
sudo apt install -y sqlite3 build-essential
```

---

#### ❌ "Error: can't find package"

**Causa:** Dependencias de Go no descargadas correctamente

**Solución:**
```bash
cd /opt/CodigoSH

# Limpiar y descargar de nuevo
go clean -modcache
go mod download
go mod tidy

# Compilar
export CGO_ENABLED=1
go build -o codigosH ./cmd/codigosH
```

---

#### ❌ "permission denied" al ejecutar

**Causa:** El binario no tiene permisos de ejecución

**Solución:**
```bash
chmod +x /opt/CodigoSH/codigosH
sudo systemctl restart codigosH
```

---

#### ❌ "config file not found" o "YAML parse error"

**Causa:** Falta el archivo config.yaml o está mal formateado

**Solución:**
```bash
# Crear directorio de configuración
sudo mkdir -p /opt/CodigoSH/configs

# Crear archivo config.yaml
sudo tee /opt/CodigoSH/configs/config.yaml > /dev/null <<'EOF'
server:
  host: "0.0.0.0"
  port: "8080"

database:
  path: "/opt/CodigoSH/codigosH.db"

logging:
  level: "info"
EOF

# Reiniciar servicio
sudo systemctl restart codigosH
```

---

#### ❌ "database locked" o "database is locked"

**Causa:** Múltiples instancias accediendo a la BD o archivo corrupto

**Solución:**
```bash
# Parar el servicio
sudo systemctl stop codigosH

# Eliminar BD corrupta
sudo rm -f /opt/CodigoSH/codigosH.db

# Iniciar nuevamente (se creará una BD nueva)
sudo systemctl start codigosH

# Esperar a que se inicialice
sleep 3

# Verificar estado
sudo systemctl status codigosH
```

---

#### ❌ Puerto 8080 ya está en uso

**Causa:** Otra aplicación ocupa el puerto 8080

**Solución:**
```bash
# Ver qué proceso usa el puerto
sudo netstat -tlnp | grep 8080
# o
sudo ss -tlnp | grep 8080

# Opción 1: Matar el proceso conflictivo
sudo kill -9 <PID>

# Opción 2: Cambiar puerto en config.yaml
sudo nano /opt/CodigoSH/configs/config.yaml
# Cambiar "port: 8080" a "port: 9090"

sudo systemctl restart codigosH
```

---

#### ❌ "Connection refused" al acceder a http://localhost:8080

**Causa:** El servicio no está corriendo en la interfaz correcta

**Solución:**
```bash
# Verificar si el servicio está activo
sudo systemctl is-active codigosH

# Ver si escucha en el puerto
sudo ss -tlnp | grep codigosH

# Si no escucha, ver logs de error
sudo journalctl -u codigosH -n 100 --no-pager

# Reiniciar servicio
sudo systemctl restart codigosH

# Verificar desde el host (si es LXC):
curl http://ip-del-lxc:8080/login
```

---

### Para LXC en Proxmox

#### Problema: No puedo acceder desde el host

**Solución:**

```bash
# 1. Verificar que el servicio está corriendo dentro del LXC
sudo systemctl status codigosH

# 2. Verificar que escucha en 0.0.0.0 (no solo localhost)
sudo ss -tlnp | grep 8080
# Deberías ver: 0.0.0.0:8080

# 3. Si usa localhost, editar config.yaml
sudo nano /opt/CodigoSH/configs/config.yaml
# Cambiar host a "0.0.0.0"

# 4. Reiniciar
sudo systemctl restart codigosH

# 5. Desde el host, acceder con la IP del LXC:
curl http://192.168.x.x:8080/login

# 6. Si aún no funciona, verificar firewall del LXC
sudo ufw status
sudo ufw allow 8080/tcp
```

---

#### Problema: Scripts de instalación falla en Debian Bookworm (12)

**Solución:**
```bash
# Go package en Debian 12 puede ser muy antiguo
# Instalar Go desde repositorio oficial:

sudo apt remove golang-go
cd /tmp
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' | sudo tee -a /etc/profile
source /etc/profile
go version
```

---

### Usar el script de diagnóstico

```bash
# Descargar y ejecutar el diagnóstico
cd /opt/CodigoSH
sudo bash scripts/diagnose.sh
```

Este script verificará:
- ✓ Servicio systemd
- ✓ Binario compilado
- ✓ Configuración
- ✓ Base de datos
- ✓ Archivos web
- ✓ Estado del servicio
- ✓ Puertos en escucha
- ✓ Conectividad

---

### Logs completos para debugging

```bash
# Ver todos los logs desde el inicio
sudo journalctl -u codigosH --all

# Ver últimas 100 líneas
sudo journalctl -u codigosH -n 100 --no-pager

# Ver logs en tiempo real con follow
sudo journalctl -u codigosH -f

# Ver solo errores
sudo journalctl -u codigosH -p err --no-pager

# Exportar logs a archivo
sudo journalctl -u codigosH > /tmp/codigosH_logs.txt
cat /tmp/codigosH_logs.txt
```

---

### Reinicio limpio completo

Si todo falla, hacer un reset completo:

```bash
# 1. Parar el servicio
sudo systemctl stop codigosH

# 2. Limpiar datos
sudo rm -rf /opt/CodigoSH/codigosH.db
sudo rm -rf /opt/CodigoSH/configs/config.yaml

# 3. Reinstalar script
cd /opt/CodigoSH
sudo bash scripts/install.sh

# O compilar manualmente:
export CGO_ENABLED=1
go mod tidy
go build -o codigosH ./cmd/codigosH

# 4. Crear configuración
sudo mkdir -p /opt/CodigoSH/configs
sudo tee /opt/CodigoSH/configs/config.yaml > /dev/null <<'EOF'
server:
  host: "0.0.0.0"
  port: "8080"

database:
  path: "/opt/CodigoSH/codigosH.db"

logging:
  level: "info"
EOF

# 5. Reiniciar servicio
sudo systemctl restart codigosH
sleep 3
sudo systemctl status codigosH
```

---

### Verificación final

```bash
# 1. Verificar que el servicio está activo
sudo systemctl is-active codigosH
# Debería mostrar: active

# 2. Verificar que escucha en puerto 8080
sudo ss -tlnp | grep 8080
# Debería mostrar: LISTEN

# 3. Verificar que responde
curl -I http://localhost:8080/login
# Debería mostrar HTTP 200 o 302 (redirect)

# 4. Ver logs recientes
sudo journalctl -u codigosH -n 10 --no-pager
```

Si todo está ✅, puedes acceder en `http://IP-DEL-LXC:8080`

---

### Obtener ayuda

Si los problemas persisten:

```bash
# Recopilar información de diagnóstico
sudo bash scripts/diagnose.sh > /tmp/diagnostico.txt
sudo journalctl -u codigosH -n 100 >> /tmp/diagnostico.txt
cat /tmp/diagnostico.txt
```

Compartir el contenido de `/tmp/diagnostico.txt` en un issue de GitHub.
