#!/bin/bash

# CodigoSH - Script de Diagnóstico
# Uso: sudo bash /path/to/diagnostico.sh

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║      CodigoSH v0.2.7-Beta - Script de Diagnóstico       ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════╝${NC}\n"

check_status() {
  if [ $1 -eq 0 ]; then
    echo -e "${GREEN}✅${NC}"
  else
    echo -e "${RED}❌${NC}"
  fi
}

# 1. Sistema Operativo
echo -e "${YELLOW}🖥️  Sistema Operativo:${NC}"
if [ -f /etc/os-release ]; then
  . /etc/os-release
  echo "   Distribución: $NAME $VERSION_ID"
  check_status 0
else
  echo "   ❌ No se pudo identificar"
  check_status 1
fi

# 2. Conectividad a Internet
echo -e "\n${YELLOW}🌐 Conectividad a Internet:${NC}"
if ping -c 1 -W 2 8.8.8.8 &>/dev/null; then
  echo "   ✅ Conectividad externa OK"
else
  echo "   ❌ Sin conectividad a internet"
fi

# 3. DNS
echo -e "\n${YELLOW}📡 Configuración DNS:${NC}"
grep nameserver /etc/resolv.conf 2>/dev/null | head -2 || echo "   ❌ No hay nameservers configurados"

# 4. Almacenamiento
echo -e "\n${YELLOW}💾 Almacenamiento:${NC}"
df -h / | tail -1 | awk '{printf "   Disponible: %s / Total: %s\n", $4, $2}'

# 5. Memoria RAM
echo -e "\n${YELLOW}🧠 Memoria RAM:${NC}"
free -h | grep Mem | awk '{printf "   Disponible: %s / Total: %s\n", $7, $2}'

# 6. CPU
echo -e "\n${YELLOW}⚙️  Procesador:${NC}"
nproc &>/dev/null && echo "   Cores: $(nproc)" || echo "   No disponible"

# 7. Dependencias
echo -e "\n${YELLOW}📦 Dependencias:${NC}"

echo -n "   Git: "
git --version &>/dev/null && echo -e "${GREEN}✅${NC} $(git --version | awk '{print $3}')" || echo -e "${RED}❌${NC}"

echo -n "   Go: "
go version &>/dev/null && echo -e "${GREEN}✅${NC} $(go version | awk '{print $3}')" || echo -e "${RED}❌${NC}"

echo -n "   GCC: "
gcc --version &>/dev/null && echo -e "${GREEN}✅${NC}" || echo -e "${RED}❌${NC}"

echo -n "   SQLite3: "
sqlite3 --version &>/dev/null && echo -e "${GREEN}✅${NC} $(sqlite3 --version)" || echo -e "${RED}❌${NC}"

echo -n "   Curl: "
curl --version &>/dev/null && echo -e "${GREEN}✅${NC}" || echo -e "${RED}❌${NC}"

# 8. Puertos
echo -e "\n${YELLOW}🔌 Puertos:${NC}"

echo -n "   Puerto 8080: "
if netstat -tlnp 2>/dev/null | grep -q ":8080 "; then
  echo -e "${YELLOW}⚠️  EN USO${NC}"
  netstat -tlnp 2>/dev/null | grep ":8080 " | awk '{print "      ("$NF")"}'
else
  echo -e "${GREEN}✅ Disponible${NC}"
fi

echo -n "   Puerto 22 (SSH): "
netstat -tlnp 2>/dev/null | grep -q ":22 " && echo -e "${GREEN}✅${NC}" || echo -e "${RED}❌${NC}"

# 9. Servicio CodigoSH
echo -e "\n${YELLOW}🚀 Servicio CodigoSH:${NC}"

echo -n "   Binario: "
if [ -f /opt/CodigoSH/codigosH ]; then
  echo -e "${GREEN}✅${NC} /opt/CodigoSH/codigosH"
  echo -n "   Permisos: "
  [ -x /opt/CodigoSH/codigosH ] && echo -e "${GREEN}✅ Ejecutable${NC}" || echo -e "${RED}❌ No ejecutable${NC}"
else
  echo -e "${RED}❌ No encontrado${NC}"
fi

echo -n "   Systemd unit: "
if [ -f /etc/systemd/system/codigosH.service ]; then
  echo -e "${GREEN}✅${NC}"
else
  echo -e "${RED}❌ No encontrado${NC}"
fi

echo -n "   Estado del servicio: "
if systemctl is-active --quiet codigosH 2>/dev/null; then
  echo -e "${GREEN}✅ CORRIENDO${NC}"
  systemctl is-enabled codigosH &>/dev/null && echo "      Autostart: ${GREEN}✅${NC}" || echo "      Autostart: ${RED}❌${NC}"
else
  echo -e "${RED}❌ DETENIDO${NC}"
  systemctl is-enabled codigosH &>/dev/null && echo "      Autostart: ${GREEN}✅${NC}" || echo "      Autostart: ${RED}❌${NC}"
fi

# 10. Base de datos
echo -e "\n${YELLOW}💾 Base de datos:${NC}"

echo -n "   Archivo: "
if [ -f /opt/CodigoSH/codigosH.db ]; then
  SIZE=$(du -h /opt/CodigoSH/codigosH.db | awk '{print $1}')
  echo -e "${GREEN}✅${NC} ($SIZE)"
else
  echo -e "${YELLOW}⚠️  No existe (se creará al iniciar)${NC}"
fi

# 11. Logs
echo -e "\n${YELLOW}📋 Logs:${NC}"

echo -n "   Archivo log: "
if [ -f /var/log/codigosH.log ]; then
  SIZE=$(du -h /var/log/codigosH.log | awk '{print $1}')
  LINES=$(wc -l < /var/log/codigosH.log)
  echo -e "${GREEN}✅${NC} ($SIZE, $LINES líneas)"
  echo -e "\n   Últimas 5 líneas:"
  tail -5 /var/log/codigosH.log | sed 's/^/      /'
else
  echo -e "${YELLOW}⚠️  No existe${NC}"
fi

# 12. Información del directorio
echo -e "\n${YELLOW}📁 Directorio /opt/CodigoSH:${NC}"

if [ -d /opt/CodigoSH ]; then
  echo "   Contenido:"
  ls -lh /opt/CodigoSH | grep -E "^(-|d)" | awk '{printf "      %s %10s %s\n", $1, $5, $9}' | head -10
else
  echo -e "   ${RED}❌ Directorio no existe${NC}"
fi

# 13. Recomendaciones
echo -e "\n${YELLOW}💡 Recomendaciones:${NC}"

ISSUES=0

if ! git --version &>/dev/null; then
  echo "   • Instala Git: sudo apt-get install -y git"
  ISSUES=$((ISSUES+1))
fi

if ! go version &>/dev/null; then
  echo "   • Instala Go: sudo apt-get install -y golang-go"
  ISSUES=$((ISSUES+1))
fi

if ! gcc --version &>/dev/null; then
  echo "   • Instala GCC: sudo apt-get install -y build-essential"
  ISSUES=$((ISSUES+1))
fi

if [ $ISSUES -eq 0 ]; then
  echo "   ✅ Sistema listo para instalar CodigoSH"
  echo "   Ejecuta: curl -sSL https://raw.githubusercontent.com/kiwinh0/CodigoSH/main/scripts/install.sh | sudo bash"
else
  echo "   ⚠️  Faltan $ISSUES dependencias. Instálalas antes de continuar."
fi

# 14. Comandos útiles
echo -e "\n${BLUE}📚 Comandos útiles:${NC}"
echo "   Ver estado: sudo systemctl status codigosH"
echo "   Ver logs: sudo journalctl -u codigosH -f"
echo "   Reiniciar: sudo systemctl restart codigosH"
echo "   Detener: sudo systemctl stop codigosH"

echo -e "\n${GREEN}════════════════════════════════════════════════════════════${NC}\n"
