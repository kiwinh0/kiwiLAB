#!/bin/bash

# CodigoSH - Script para diagnosticar problemas en LXC

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}╔════════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║          CodigoSH - Guía de Solución de Problemas             ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════════╝${NC}"

# Función para preguntar al usuario
ask_user() {
    local prompt="$1"
    read -p "$(echo -e ${BLUE}$prompt${NC}): " -n 1 -r
    echo
    [[ $REPLY =~ ^[Ss]$ ]]
}

# 1. Verificar Go
echo -e "\n${YELLOW}1. Verificando Go...${NC}"
if command -v go &> /dev/null; then
    GO_VERSION=$(go version | awk '{print $3}')
    echo -e "${GREEN}✓${NC} Go está instalado: $GO_VERSION"
else
    echo -e "${RED}✗${NC} Go NO está instalado"
    echo "Instalando Go..."
    sudo apt update && sudo apt install -y golang-go
fi

# 2. Verificar compilación
echo -e "\n${YELLOW}2. Verificando compilación...${NC}"
if [ -f /opt/CodigoSH/codigosH ]; then
    echo -e "${GREEN}✓${NC} Binario codigosH existe"
else
    echo -e "${RED}✗${NC} Binario NO existe"
    if ask_user "¿Intentar compilar ahora?"; then
        cd /opt/CodigoSH
        export CGO_ENABLED=1
        echo "Compilando..."
        go build -o codigosH ./cmd/codigosH
        if [ $? -eq 0 ]; then
            echo -e "${GREEN}✓ Compilación exitosa${NC}"
        else
            echo -e "${RED}✗ Error en compilación${NC}"
            echo "Ver logs: sudo journalctl -u codigosH -n 50"
        fi
    fi
fi

# 3. Verificar configuración
echo -e "\n${YELLOW}3. Verificando configuración...${NC}"
if [ -f /opt/CodigoSH/configs/config.yaml ]; then
    echo -e "${GREEN}✓${NC} Archivo config.yaml existe"
    HOST=$(grep -A2 "^server:" /opt/CodigoSH/configs/config.yaml | grep "host:" | awk '{print $2}')
    PORT=$(grep -A2 "^server:" /opt/CodigoSH/configs/config.yaml | grep "port:" | awk '{print $2}')
    echo "   Host: $HOST"
    echo "   Puerto: $PORT"
    
    if [ "$HOST" == "localhost" ] || [ "$HOST" == "127.0.0.1" ]; then
        echo -e "${YELLOW}⚠${NC} Host es localhost - no será accesible desde el host"
        if ask_user "¿Cambiar a 0.0.0.0 para permitir conexiones remotas?"; then
            sudo sed -i 's/host: "localhost"/host: "0.0.0.0"/' /opt/CodigoSH/configs/config.yaml
            sudo sed -i 's/host: "127.0.0.1"/host: "0.0.0.0"/' /opt/CodigoSH/configs/config.yaml
            echo -e "${GREEN}✓ Configuración actualizada${NC}"
            sudo systemctl restart codigosH
        fi
    fi
else
    echo -e "${RED}✗${NC} Archivo config.yaml NO existe"
    if ask_user "¿Crear configuración por defecto?"; then
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
        echo -e "${GREEN}✓ Configuración creada${NC}"
        sudo systemctl restart codigosH
    fi
fi

# 4. Verificar estado del servicio
echo -e "\n${YELLOW}4. Verificando estado del servicio...${NC}"
if systemctl is-active --quiet codigosH; then
    echo -e "${GREEN}✓${NC} Servicio codigosH está ACTIVO"
    PID=$(systemctl show -p MainPID --value codigosH)
    echo "   PID: $PID"
else
    echo -e "${RED}✗${NC} Servicio codigosH está INACTIVO"
    if ask_user "¿Intentar iniciar el servicio?"; then
        sudo systemctl start codigosH
        sleep 2
        if systemctl is-active --quiet codigosH; then
            echo -e "${GREEN}✓ Servicio iniciado exitosamente${NC}"
        else
            echo -e "${RED}✗ Error al iniciar${NC}"
            echo "Ver logs con: sudo journalctl -u codigosH -n 50"
        fi
    fi
fi

# 5. Verificar puerto
echo -e "\n${YELLOW}5. Verificando puerto 8080...${NC}"
if ss -tlnp 2>/dev/null | grep -q :8080 || netstat -tlnp 2>/dev/null | grep -q :8080; then
    echo -e "${GREEN}✓${NC} Puerto 8080 está en ESCUCHA"
    LISTENING_ADDR=$(ss -tlnp 2>/dev/null | grep :8080 | awk '{print $4}' || netstat -tlnp 2>/dev/null | grep :8080 | awk '{print $4}')
    echo "   Escuchando en: $LISTENING_ADDR"
else
    echo -e "${RED}✗${NC} Puerto 8080 NO está en escucha"
    echo "El servicio puede no estar corriendo correctamente"
fi

# 6. Verificar conectividad local
echo -e "\n${YELLOW}6. Verificando conectividad local...${NC}"
if curl -s http://localhost:8080/login > /dev/null 2>&1; then
    echo -e "${GREEN}✓${NC} Responde en http://localhost:8080"
else
    echo -e "${RED}✗${NC} No responde en http://localhost:8080"
fi

# 7. Obtener IP del LXC
echo -e "\n${YELLOW}7. Información de conectividad...${NC}"
IP_LOCAL=$(hostname -I | awk '{print $1}')
if [ -n "$IP_LOCAL" ]; then
    echo -e "${GREEN}✓${NC} IP del LXC: $IP_LOCAL"
    echo "   Puedes acceder desde el host en: http://$IP_LOCAL:8080"
    
    # Verificar si es accesible desde el host
    if ask_user "¿Verificar accesibilidad desde el host (si estás en el LXC)?"; then
        echo "Ejecuta esto en el HOST:"
        echo "   curl -I http://$IP_LOCAL:8080/login"
    fi
else
    echo -e "${RED}✗${NC} No se pudo obtener IP del LXC"
fi

# 8. Ver logs
echo -e "\n${YELLOW}8. Últimos logs del servicio...${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
sudo journalctl -u codigosH -n 20 --no-pager
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# 9. Resumen y recomendaciones
echo -e "\n${BLUE}═════════════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}RESUMEN Y RECOMENDACIONES${NC}"
echo -e "${BLUE}═════════════════════════════════════════════════════════════════${NC}"

ISSUES=0

if ! command -v go &> /dev/null; then
    echo -e "${RED}✗${NC} Go no está instalado"
    ((ISSUES++))
fi

if [ ! -f /opt/CodigoSH/codigosH ]; then
    echo -e "${RED}✗${NC} Binario no existe"
    ((ISSUES++))
fi

if [ ! -f /opt/CodigoSH/configs/config.yaml ]; then
    echo -e "${RED}✗${NC} Configuración no existe"
    ((ISSUES++))
fi

if ! systemctl is-active --quiet codigosH; then
    echo -e "${RED}✗${NC} Servicio no está corriendo"
    ((ISSUES++))
fi

if [ $ISSUES -eq 0 ]; then
    echo -e "\n${GREEN}✅ TODO PARECE ESTAR BIEN!${NC}"
    echo -e "\nAccede a: ${BLUE}http://$IP_LOCAL:8080${NC}"
    echo -e "Usuario por defecto: ${BLUE}admin${NC}"
    echo -e "Contraseña: ${BLUE}admin${NC}"
else
    echo -e "\n${YELLOW}⚠ Se encontraron $ISSUES problema(s)${NC}"
    echo -e "\n📚 Consulta la guía de troubleshooting:"
    echo -e "   ${BLUE}https://github.com/kiwinh0/CodigoSH/blob/main/TROUBLESHOOTING.md${NC}"
fi

echo -e "\n${BLUE}═════════════════════════════════════════════════════════════════${NC}"
echo -e "\n📚 Comandos útiles:"
echo "   Ver estado: sudo systemctl status codigosH"
echo "   Reiniciar: sudo systemctl restart codigosH"
echo "   Ver logs vivos: sudo journalctl -u codigosH -f"
echo "   Compilar manualmente: cd /opt/CodigoSH && export CGO_ENABLED=1 && go build -o codigosH ./cmd/codigosH"
echo ""
