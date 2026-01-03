# 🚀 Guía de Instalación de CodigoSH en LXC (Proxmox)

## ✅ Pre-requisitos

Tu contenedor LXC debe tener:
- **OS**: Debian 11+ o Ubuntu 20.04+
- **CPU**: Mínimo 2 cores
- **RAM**: Mínimo 512MB (recomendado 1GB)
- **Almacenamiento**: Mínimo 2GB disponible
- **Red**: Acceso a internet (para descargar dependencias)

## 🔧 Paso 1: Verificar el Contenedor

Dentro de tu LXC, ejecuta:

```bash
# Verificar versión de Debian/Ubuntu
cat /etc/os-release

# Verificar conectividad
curl -I https://github.com

# Verificar espacio disponible
df -h

# Verificar memoria
free -h
```

## 📦 Paso 2: Instalación Automática

### Opción A: Descarga e Instalación Directa (Recomendado)

```bash
# 1. Conectar al LXC
pct shell <CONTAINER_ID>

# 2. Ejecutar el script de instalación
curl -sSL "https://raw.githubusercontent.com/kiwinh0/CodigoSH/main/scripts/install.sh" | sudo bash
```

### Opción B: Instalación Manual

```bash
# 1. Actualizar sistema
sudo apt-get update
sudo apt-get install -y git golang-go sqlite3 build-essential ca-certificates

# 2. Clonar repositorio
cd /opt
sudo git clone --depth=1 https://github.com/kiwinh0/CodigoSH.git

# 3. Compilar
cd /opt/CodigoSH
sudo bash -c 'export CGO_ENABLED=1 GOOS=linux && go mod download && go build -o codigosH ./cmd/codigosH/main.go'

# 4. Crear servicio systemd
sudo tee /etc/systemd/system/codigosH.service > /dev/null <<'EOF'
[Unit]
Description=CodigoSH Dashboard
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=root
WorkingDirectory=/opt/CodigoSH
ExecStart=/opt/CodigoSH/codigosH
Restart=on-failure
RestartSec=10
StandardOutput=append:/var/log/codigosH.log
StandardError=append:/var/log/codigosH.log

[Install]
WantedBy=multi-user.target
EOF

# 5. Iniciar servicio
sudo systemctl daemon-reload
sudo systemctl enable codigosH
sudo systemctl start codigosH
```

## 🔍 Verificación de la Instalación

```bash
# Ver estado del servicio
sudo systemctl status codigosH

# Ver logs en tiempo real
sudo journalctl -u codigosH -f

# Verificar puerto
sudo netstat -tlnp | grep 8080
```

## 🌐 Acceso

Una vez instalado, accede desde tu navegador:

- **URL**: `http://<IP_DEL_LXC>:8080`
- **Usuario**: `admin`
- **Contraseña**: `admin`

Para obtener la IP del contenedor:
```bash
pct exec <CONTAINER_ID> hostname -I
```

## ⚠️ Solución de Problemas

### El servicio no inicia

```bash
# Ver logs completos
sudo journalctl -u codigosH -n 100

# Intentar arranque manual para ver errores
cd /opt/CodigoSH
./codigosH

# Verificar permisos
ls -la /opt/CodigoSH/codigosH
```

### Puerto 8080 en uso

```bash
# Ver qué está usando el puerto
sudo netstat -tlnp | grep 8080

# O matar el proceso anterior
sudo pkill -9 codigosH
sudo systemctl restart codigosH
```

### Error de compilación: "cgo: gcc not found"

```bash
sudo apt-get install -y build-essential
sudo apt-get install -y gcc
```

### Error: "go: command not found"

Verifica la versión de Go:
```bash
go version

# Si no está instalado
sudo apt-get install -y golang-go

# Verificar que Go esté en PATH
echo $PATH
```

### Contenedor sin acceso a internet

Dentro del LXC:
```bash
# Verificar DNS
cat /etc/resolv.conf

# Probar conectividad
ping 8.8.8.8
curl https://github.com
```

Si no funciona, en Proxmox (en el nodo host):
```bash
# Ver configuración de red del LXC
cat /etc/pve/lxc/<CONTAINER_ID>.conf

# Puede ser necesario:
# - Añadir gateway
# - Configurar DNS
# - Habilitar NAT/Bridge networking
```

## 🔐 Configuración de Seguridad

### 1. Cambiar credenciales por defecto

Una vez instalado, accede a Settings y cambia las credenciales.

### 2. Habilitar HTTPS (Opcional)

```bash
# Instalar Nginx como reverse proxy
sudo apt-get install -y nginx certbot python3-certbot-nginx

# Configurar certificado SSL
sudo certbot certonly --standalone -d tu-dominio.com

# Crear configuración Nginx en /etc/nginx/sites-available/codigosh
# (Consulta la documentación de Nginx)
```

### 3. Firewall

```bash
# Habilitar UFW
sudo ufw enable
sudo ufw allow 22/tcp   # SSH
sudo ufw allow 8080/tcp # CodigoSH

# Verificar
sudo ufw status
```

## 🚀 Comandos Útiles

```bash
# Detener servicio
sudo systemctl stop codigosH

# Reiniciar servicio
sudo systemctl restart codigosH

# Habilitar/Deshabilitar inicio automático
sudo systemctl enable codigosH
sudo systemctl disable codigosH

# Ver últimas líneas de logs
sudo journalctl -u codigosH -n 50

# Limpiar base de datos (cuidado)
sudo rm /opt/CodigoSH/codigosH.db
sudo systemctl restart codigosH

# Actualizar a la última versión
cd /opt/CodigoSH
sudo git pull origin main
sudo bash -c 'export CGO_ENABLED=1 GOOS=linux && go build -o codigosH ./cmd/codigosH/main.go'
sudo systemctl restart codigosH
```

## 📞 Soporte

Si tienes problemas:

1. **Revisa los logs**: `sudo journalctl -u codigosH -n 100`
2. **Verifica conectividad**: `curl http://localhost:8080`
3. **Comprueba permisos**: `ls -la /opt/CodigoSH/`
4. **Abre un issue**: [GitHub Issues](https://github.com/kiwinh0/CodigoSH/issues)

---

**Documentación**: https://github.com/kiwinh0/CodigoSH  
**Versión**: v0.1.0-Beta  
**Última actualización**: 3 de enero de 2026
