#!/bin/bash

# Script de diagnóstico para CodigoSH en LXC/Proxmox

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}════════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}🔍 Diagnóstico CodigoSH en LXC/Proxmox${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════${NC}\n"

# 1. Verificar si el binario existe
echo -e "${YELLOW}1. Verificando binario...${NC}"
if [ -f /opt/CodigoSH/codigosH ]; then
  echo -e "${GREEN}✓${NC} Binario encontrado"
  file /opt/CodigoSH/codigosH
else
  echo -e "${RED}✗${NC} Binario NO encontrado en /opt/CodigoSH"
  echo "  Solución: Re-ejecuta el script de instalación"
  exit 1
fi

# 2. Verificar permisos
echo -e "\n${YELLOW}2. Verificando permisos...${NC}"
if [ -x /opt/CodigoSH/codigosH ]; then
  echo -e "${GREEN}✓${NC} Binario es ejecutable"
else
  echo -e "${RED}✗${NC} Binario NO es ejecutable"
  echo "  Arreglando..."
  sudo chmod +x /opt/CodigoSH/codigosH
fi

# 3. Verificar puerto 8080
echo -e "\n${YELLOW}3. Verificando puerto 8080...${NC}"
if sudo netstat -tuln 2>/dev/null | grep -q ":8080" || sudo ss -tuln 2>/dev/null | grep -q ":8080"; then
  echo -e "${GREEN}✓${NC} Puerto 8080 está en uso"
  sudo netstat -tuln 2>/dev/null | grep 8080 || sudo ss -tuln 2>/dev/null | grep 8080
else
  echo -e "${YELLOW}⚠${NC} Puerto 8080 no está en uso (normal si el servicio no está corriendo)"
fi

# 4. Verificar servicio systemd
echo -e "\n${YELLOW}4. Verificando servicio systemd...${NC}"
if systemctl is-enabled codigosH 2>/dev/null; then
  echo -e "${GREEN}✓${NC} Servicio habilitado al arranque"
else
  echo -e "${YELLOW}⚠${NC} Servicio NO está habilitado"
  echo "  Ejecuta: sudo systemctl enable codigosH"
fi

if systemctl is-active --quiet codigosH 2>/dev/null; then
  echo -e "${GREEN}✓${NC} Servicio está ACTIVO"
  systemctl status codigosH | grep "Active:"
else
  echo -e "${RED}✗${NC} Servicio NO está activo"
  echo "  Intentando iniciar..."
  sudo systemctl start codigosH
  sleep 2
  if systemctl is-active --quiet codigosH 2>/dev/null; then
    echo -e "${GREEN}✓${NC} Servicio iniciado exitosamente"
  else
    echo -e "${RED}✗${NC} Fallo al iniciar el servicio"
  fi
fi

# 5. Revisar logs
echo -e "\n${YELLOW}5. Últimos logs del servicio:${NC}"
journalctl -u codigosH -n 10 2>/dev/null | tail -10 || echo "No hay logs disponibles"

# 6. Verificar acceso
echo -e "\n${YELLOW}6. Probando acceso...${NC}"
IP=$(hostname -I 2>/dev/null | awk '{print $1}' || echo "127.0.0.1")
if curl -s http://localhost:8080/login > /dev/null 2>&1; then
  echo -e "${GREEN}✓${NC} Servidor responde en http://localhost:8080"
  echo -e "${GREEN}✓${NC} Acceso: http://${IP}:8080"
else
  echo -e "${RED}✗${NC} Servidor NO responde"
  
  # Intentar diagnóstico más profundo
  echo -e "\n${YELLOW}7. Intentando diagnóstico profundo...${NC}"
  
  # Verificar si el contenedor tiene permisos de red
  if [ -f /proc/net/tcp ]; then
    echo -e "${GREEN}✓${NC} Networking disponible en LXC"
  else
    echo -e "${RED}✗${NC} Problema de networking en LXC"
  fi
  
  # Intentar ejecutar manualmente
  echo -e "\n${YELLOW}Intentando ejecución manual:${NC}"
  cd /opt/CodigoSH
  timeout 5 /opt/CodigoSH/codigosH 2>&1 || true
fi

# 8. Resumen y recomendaciones
echo -e "\n${BLUE}════════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}📋 RESUMEN Y RECOMENDACIONES${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════${NC}\n"

echo "Si CodigoSH no está funcionando:"
echo ""
echo "1️⃣  Reinicia el servicio:"
echo "   ${YELLOW}sudo systemctl restart codigosH${NC}"
echo ""
echo "2️⃣  Comprueba los logs en tiempo real:"
echo "   ${YELLOW}sudo journalctl -u codigosH -f${NC}"
echo ""
echo "3️⃣  Si aún no funciona, ejecuta manualmente:"
echo "   ${YELLOW}cd /opt/CodigoSH && ./codigosH${NC}"
echo ""
echo "4️⃣  Acceso en navegador:"
echo "   ${YELLOW}http://${IP}:8080${NC}"
echo "   Usuario: ${YELLOW}admin${NC}"
echo "   Contraseña: ${YELLOW}admin${NC}"
echo ""
echo "5️⃣  Si vuelve a fallar, reconstruye:"
echo "   ${YELLOW}cd /opt/CodigoSH${NC}"
echo "   ${YELLOW}go build -o codigosH ./cmd/codigosH/main.go${NC}"
echo "   ${YELLOW}sudo systemctl restart codigosH${NC}"
echo ""
