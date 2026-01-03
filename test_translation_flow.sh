#!/bin/bash

echo "🔍 Test de Flujo Completo de Traducción"
echo "========================================"
echo ""

# Crear archivo de cookies
COOKIES="/tmp/codigosh_cookies.txt"
rm -f "$COOKIES"

# Test 1: Obtener página de login (sin cookies previas)
echo "1️⃣ GET /login (idioma por defecto)..."
RESPONSE=$(curl -s -c "$COOKIES" http://localhost:8080/login)
if echo "$RESPONSE" | grep -q 'lang="es"'; then
    echo "   ✅ Login en español (por defecto)"
else
    echo "   ⚠️ Verificar idioma por defecto"
fi

# Ver las cookies guardadas
echo ""
echo "2️⃣ Cookies después de GET /login:"
cat "$COOKIES" 2>/dev/null | grep -v "^#" || echo "   (No hay cookies aún)"

# Test 2: Hacer login (simular)
echo ""
echo "3️⃣ POST /login (iniciar sesión)..."
# Nota: Esto fallará porque el usuario no existe, pero podemos ver el flujo
RESPONSE=$(curl -s -b "$COOKIES" -c "$COOKIES" -X POST http://localhost:8080/login \
  -d "username=testuser&password=password123&remember_me=on" -L)

echo "   Response: $(echo "$RESPONSE" | head -5 | tr '\n' ' ')..."

# Test 3: Verificar estructura de archivos JSON
echo ""
echo "4️⃣ Verificar estructura JSON en es.json..."
jq '.login | keys | length' /Users/kiwinho/Proyectos/CodigoSH/web/static/i18n/es.json 2>/dev/null && echo "   ✅ Login translations presente" || echo "   ❌ Error en JSON"

echo ""
echo "5️⃣ Verificar estructura JSON en en.json..."
jq '.login | keys | length' /Users/kiwinho/Proyectos/CodigoSH/web/static/i18n/en.json 2>/dev/null && echo "   ✅ English login translations presente" || echo "   ❌ Error en JSON"

# Test 4: Verificar que los métodos existen en JavaScript
echo ""
echo "6️⃣ Verificar métodos JavaScript..."
if grep -q "translate(key)" /Users/kiwinho/Proyectos/CodigoSH/web/static/i18n/i18n.js; then
    echo "   ✅ Método translate() existe"
else
    echo "   ❌ Método translate() NO existe"
fi

if grep -q "persistLanguage" /Users/kiwinho/Proyectos/CodigoSH/web/static/i18n/i18n.js; then
    echo "   ✅ Método persistLanguage() existe"
else
    echo "   ❌ Método persistLanguage() NO existe"
fi

# Test 5: Ver la base de datos
echo ""
echo "7️⃣ Verificar usuario de prueba en BD..."
sqlite3 /Users/kiwinho/Proyectos/CodigoSH/codigosH.db "SELECT username, language FROM users WHERE username='testuser';" 2>/dev/null | grep -q "testuser" && echo "   ✅ Usuario testuser en BD con idioma" || echo "   ⚠️ Usuario no encontrado"

echo ""
echo "========================================"
echo "✅ Tests completados"
