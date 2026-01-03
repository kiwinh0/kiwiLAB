#!/bin/bash
set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log_error() { echo -e "${RED}❌ $1${NC}" >&2; }
log_success() { echo -e "${GREEN}✅ $1${NC}"; }
log_info() { echo -e "${BLUE}ℹ️  $1${NC}"; }

echo -e "${GREEN}════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}🚀 CodigoSH Installation Script${NC}"
echo -e "${GREEN}════════════════════════════════════════════════════════════${NC}\n"

[ "$EUID" -eq 0 ] || { log_error "Ejecuta con: sudo bash $0"; exit 1; }

log_info "Actualizando repositorios..."
DEBIAN_FRONTEND=noninteractive apt-get update -qq 2>/dev/null || true

log_info "Instalando dependencias (puede tardar 2-3 minutos)..."
apt-get install -y -qq git golang-go sqlite3 build-essential ca-certificates curl wget net-tools 2>/dev/null || {
  log_error "Error instalando dependencias"
  exit 1
}
log_success "Dependencias instaladas"

log_info "Clonando repositorio..."
rm -rf /opt/CodigoSH 2>/dev/null || true
mkdir -p /opt
git clone --depth=1 https://github.com/kiwinh0/CodigoSH.git /opt/CodigoSH 2>&1 | grep -v "Cloning\|remote:" || true
cd /opt/CodigoSH || exit 1
log_success "Repositorio descargado"

log_info "Compilando CodigoSH (puede tardar 2-5 minutos)..."
export CGO_ENABLED=1 GOOS=linux GOARCH=amd64
go mod download 2>/dev/null || { log_error "Error descargando dependencias Go"; exit 1; }
go build -o codigosH ./cmd/codigosH/main.go 2>/dev/null || { log_error "Error compilando"; exit 1; }
[ -f codigosH ] || { log_error "Binario no generado"; exit 1; }
chmod +x codigosH
log_success "Compilación completada"

log_info "Configurando sistema..."
mkdir -p /opt/CodigoSH/configs

cat > /etc/systemd/system/codigosH.service <<'SERVICE_EOF'
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
Environment="PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"

[Install]
WantedBy=multi-user.target
SERVICE_EOF

systemctl daemon-reload 2>/dev/null || true
systemctl enable codigosH 2>/dev/null || log_info "Advertencia: systemctl enable falló"

log_info "Iniciando servicio..."
if systemctl start codigosH 2>/dev/null; then
  sleep 3
  
  if systemctl is-active --quiet codigosH 2>/dev/null; then
    IP=$(hostname -I 2>/dev/null | awk '{print $1}' || echo "127.0.0.1")
    log_success "¡Instalación completada!"
    echo ""
    echo -e "${GREEN}════════════════════════════════════════════════════════════${NC}"
    echo -e "📍 Acceso: ${YELLOW}http://$IP:8080${NC}"
    echo -e "${GREEN}════════════════════════════════════════════════════════════${NC}"
    exit 0
  else
    log_error "Servicio falló. Logs:"
    journalctl -u codigosH -n 20 2>/dev/null || echo "No hay logs disponibles"
    echo ""
    log_info "Intentando arranque manual para diagnosticar..."
    cd /opt/CodigoSH && /opt/CodigoSH/codigosH 2>&1 | head -20
    exit 1
  fi
else
  log_error "No se pudo iniciar el servicio con systemctl"
  log_info "Intentando arranque manual..."
  cd /opt/CodigoSH
  /opt/CodigoSH/codigosH &
  sleep 2
  if pgrep -f "codigosH" > /dev/null; then
    IP=$(hostname -I 2>/dev/null | awk '{print $1}' || echo "127.0.0.1")
    log_success "CodigoSH iniciado en modo manual"
    echo -e "📍 Acceso: ${YELLOW}http://$IP:8080${NC}"
    exit 0
  else
    log_error "Error al iniciar. Revisa la configuración."
    exit 1
  fi
fi
