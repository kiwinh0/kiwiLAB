#!/bin/bash

# CodigoSH - Script de Instalación Automática
# Desarrollado para Debian/Ubuntu

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

echo -e "${GREEN}🚀 Iniciando instalación de CodigoSH...${NC}"

# 1. Comprobar permisos
if [ "$EUID" -ne 0 ]; then 
  echo -e "${RED}Por favor, ejecuta este script con sudo.${NC}"
  exit
fi

# 2. Instalar dependencias
echo -e "📦 Instalando dependencias de sistema (Go, SQLite, GCC)..."
apt update && apt install -y git golang-go sqlite3 build-essential ca-certificates

# 3. Limpiar instalaciones previas y clonar
echo -e "📂 Configurando directorio en /opt/CodigoSH..."
rm -rf /opt/CodigoSH
git clone https://github.com/kiwinh0/CodigoSH.git /opt/CodigoSH
cd /opt/CodigoSH

# 4. Compilar binario
echo -e "🛠️ Compilando CodigoSH (CGO habilitado)..."
export CGO_ENABLED=1
go mod download
go build -o codigosH ./cmd/codigosH/main.go

# 5. Crear servicio Systemd
echo -e "⚙️ Configurando servicio de sistema..."
cat <<EOF > /etc/systemd/system/codigosH.service
[Unit]
Description=CodigoSH Dashboard
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/opt/CodigoSH
ExecStart=/opt/CodigoSH/codigosH
Restart=always
Environment=NODE_ENV=production

[Install]
WantedBy=multi-user.target
EOF

# 6. Activar y arrancar
echo -e "🏁 Finalizando configuración..."
systemctl daemon-reload
systemctl enable --now codigosH

IP_LOCAL=$(hostname -I | awk '{print $1}')
echo -e "${GREEN}✅ ¡CodigoSH instalado con éxito!${NC}"
echo -e "🌐 Puedes acceder en: ${RED}http://$IP_LOCAL:8080${NC}"
