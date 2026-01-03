#!/bin/bash

# CodigoSH - Script de Diagnóstico
# Ayuda a identificar problemas de instalación

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${GREEN}🔍 CodigoSH - Diagnóstico de Instalación${NC}\n"

# 1. Verificar si el servicio existe
echo "1️⃣  Verificando servicio systemd..."
if [ -f /etc/systemd/system/codigosH.service ]; then
  echo -e "${GREEN}✓${NC} Servicio encontrado"
else
  echo -e "${RED}✗${NC} Servicio NO encontrado en /etc/systemd/system/codigosH.service"
fi

# 2. Verificar directorio de instalación
echo -e "\n2️⃣  Verificando directorio de instalación..."
if [ -d /opt/CodigoSH ]; then
  echo -e "${GREEN}✓${NC} Directorio /opt/CodigoSH existe"
  if [ -f /opt/CodigoSH/codigosH ]; then
    echo -e "${GREEN}✓${NC} Binario ejecutable encontrado"
  else
    echo -e "${RED}✗${NC} Binario NO encontrado en /opt/CodigoSH/codigosH"
  fi
else
  echo -e "${RED}✗${NC} Directorio /opt/CodigoSH NO existe"
fi

# 3. Verificar configuración
echo -e "\n3️⃣  Verificando configuración..."
if [ -f /opt/CodigoSH/configs/config.yaml ]; then
  echo -e "${GREEN}✓${NC} Archivo config.yaml encontrado"
  echo "   Contenido:"
  cat /opt/CodigoSH/configs/config.yaml | sed 's/^/   /'
else
  echo -e "${RED}✗${NC} Archivo config.yaml NO encontrado"
fi

# 4. Verificar base de datos
echo -e "\n4️⃣  Verificando base de datos..."
if [ -f /opt/CodigoSH/codigosH.db ]; then
  echo -e "${GREEN}✓${NC} Base de datos encontrada"
  echo "   Tamaño: $(du -h /opt/CodigoSH/codigosH.db | awk '{print $1}')"
else
  echo -e "${YELLOW}⚠${NC} Base de datos NO encontrada (se creará al iniciar)"
fi

# 5. Verificar directorio web
echo -e "\n5️⃣  Verificando archivos web..."
if [ -d /opt/CodigoSH/web ]; then
  echo -e "${GREEN}✓${NC} Directorio web encontrado"
  [ -d /opt/CodigoSH/web/templates ] && echo -e "${GREEN}✓${NC} Templates encontrados" || echo -e "${RED}✗${NC} Templates NO encontrados"
  [ -d /opt/CodigoSH/web/static ] && echo -e "${GREEN}✓${NC} Estáticos encontrados" || echo -e "${RED}✗${NC} Estáticos NO encontrados"
else
  echo -e "${RED}✗${NC} Directorio web NO encontrado"
fi

# 6. Verificar estado del servicio
echo -e "\n6️⃣  Estado del servicio..."
if systemctl is-active --quiet codigosH; then
  echo -e "${GREEN}✓${NC} Servicio codigosH está ACTIVO"
  PID=$(systemctl show -p MainPID --value codigosH)
  echo "   PID: $PID"
else
  echo -e "${RED}✗${NC} Servicio codigosH está INACTIVO"
fi

# 7. Verificar puertos
echo -e "\n7️⃣  Verificando puertos..."
if netstat -tlnp 2>/dev/null | grep -q :8080 || ss -tlnp 2>/dev/null | grep -q :8080; then
  echo -e "${GREEN}✓${NC} Puerto 8080 está en ESCUCHA"
else
  echo -e "${RED}✗${NC} Puerto 8080 NO está en escucha"
fi

# 8. Ver logs recientes
echo -e "\n8️⃣  Últimos logs (últimas 30 líneas)..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
journalctl -u codigosH -n 30 --no-pager 2>/dev/null || echo -e "${YELLOW}⚠${NC} No se pueden leer logs (sin permisos o systemd no disponible)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# 9. Verificar connectividad
echo -e "\n9️⃣  Verificando conectividad al servicio..."
if curl -s http://localhost:8080/login > /dev/null 2>&1; then
  echo -e "${GREEN}✓${NC} Servicio responde en http://localhost:8080"
else
  echo -e "${RED}✗${NC} No hay respuesta en http://localhost:8080"
fi

# 10. Recomendaciones
echo -e "\n🔧 RECOMENDACIONES:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if ! systemctl is-active --quiet codigosH; then
  echo "• El servicio no está corriendo. Intenta:"
  echo "  sudo systemctl start codigosH"
  echo ""
fi

if ! [ -f /opt/CodigoSH/codigosH ]; then
  echo "• El binario no está compilado. Intenta:"
  echo "  cd /opt/CodigoSH && go build -o codigosH ./cmd/codigosH"
  echo ""
fi

if ! [ -f /opt/CodigoSH/configs/config.yaml ]; then
  echo "• Falta el archivo de configuración. Crea uno en /opt/CodigoSH/configs/config.yaml"
  echo ""
fi

echo -e "\n📚 Para más ayuda:"
echo "• Ver logs en tiempo real: sudo journalctl -u codigosH -f"
echo "• Reiniciar el servicio: sudo systemctl restart codigosH"
echo "• Parar el servicio: sudo systemctl stop codigosH"
echo "• Estado completo: sudo systemctl status codigosH"
