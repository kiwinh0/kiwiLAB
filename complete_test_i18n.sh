#!/bin/bash

echo "🧪 Pruebas Completas de Traducción en CodigoSH"
echo "=============================================="
echo ""

# Test 1: Verificar que el servidor está corriendo
echo "1️⃣ Verificar servidor..."
curl -s http://localhost:8080/login > /dev/null 2>&1 && echo "✅ Servidor respondiendo" || echo "❌ Servidor no responde"

# Test 2: Verificar que login tiene lang="es"
echo ""
echo "2️⃣ Verificar login.html..."
LANG_ATTR=$(curl -s http://localhost:8080/login | grep -oP 'lang="\K[^"]+' | head -1)
echo "Lang attribute: $LANG_ATTR"
[ "$LANG_ATTR" = "es" ] && echo "✅ lang correcto" || echo "⚠️ Esperaba 'es', obtuve '$LANG_ATTR'"

# Test 3: Verificar que i18n.js está incluido
echo ""
echo "3️⃣ Verificar inclusión de i18n.js..."
curl -s http://localhost:8080/login | grep -q "i18n.js" && echo "✅ i18n.js incluido" || echo "❌ i18n.js NO encontrado"

# Test 4: Verificar que la función de traducción existe
echo ""
echo "4️⃣ Verificar función translateLoginPage..."
curl -s http://localhost:8080/login | grep -q "translateLoginPage" && echo "✅ Función definida" || echo "❌ Función NO encontrada"

# Test 5: Verificar archivo de traducción en español
echo ""
echo "5️⃣ Verificar archivo es.json..."
[ -f "/Users/kiwinho/Proyectos/CodigoSH/web/static/i18n/es.json" ] && echo "✅ Archivo existe" || echo "❌ Archivo NO existe"

# Test 6: Verificar que i18n.js tiene el método translate
echo ""
echo "6️⃣ Verificar método translate() en i18n.js..."
grep -q "translate(key)" /Users/kiwinho/Proyectos/CodigoSH/web/static/i18n/i18n.js && echo "✅ Método existe" || echo "❌ Método NO existe"

# Test 7: Verificar que persistLanguage guarda cookies
echo ""
echo "7️⃣ Verificar guardado de cookies en persistLanguage..."
grep -q "document.cookie" /Users/kiwinho/Proyectos/CodigoSH/web/static/i18n/i18n.js && echo "✅ Cookie guardado" || echo "❌ Cookie NO se guarda"

# Test 8: Verificar dashboard.html
echo ""
echo "8️⃣ Verificar dashboard.html..."
curl -s http://localhost:8080/dashboard > /dev/null 2>&1 && echo "⚠️ Dashboard disponible (requiere autenticación)" || true

# Test 9: Ver todos los idiomas en JSON files
echo ""
echo "9️⃣ Verificar archivos de idioma..."
ls /Users/kiwinho/Proyectos/CodigoSH/web/static/i18n/*.json | wc -l
echo "   archivos JSON encontrados"

# Test 10: Verificar estructura de traducciones en login
echo ""
echo "🔟 Verificar traducciones en es.json..."
jq '.login' /Users/kiwinho/Proyectos/CodigoSH/web/static/i18n/es.json | grep -c "tagline\|submit_button" && echo "✅ Estructura de login presente" || echo "❌ Falta estructura"

echo ""
echo "=============================================="
echo "✅ Pruebas completadas"
