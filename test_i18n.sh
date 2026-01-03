#!/bin/bash

# Test script to verify i18n implementation across all pages

echo "🧪 Pruebas de i18n en CodigoSH"
echo "================================"

# Test 1: Login page should have lang="es"
echo -e "\n✅ Test 1: Verificar login.html tiene lang dinámico"
curl -s http://localhost:8080/login | grep -E '^<html lang=' | head -1

# Test 2: About page (unauthenticated) should have lang="es"
echo -e "\n✅ Test 2: Verificar about.html (sin autenticación) tiene lang dinámico"
curl -s http://localhost:8080/about | grep -E '^<html lang=' | head -1

# Test 3: Verify i18n.js is included in login
echo -e "\n✅ Test 3: Verificar que i18n.js está incluido en login"
curl -s http://localhost:8080/login | grep 'i18n.js' | head -1

# Test 4: Verify i18n.js is included in about
echo -e "\n✅ Test 4: Verificar que i18n.js está incluido en about"
curl -s http://localhost:8080/about | grep 'i18n.js' | head -1

# Test 5: Verify translation functions are in pages
echo -e "\n✅ Test 5: Verificar que translateLoginPage está en login.html"
curl -s http://localhost:8080/login | grep 'translateLoginPage' | head -1

echo -e "\n✅ Test 6: Verificar que translateAboutPage está en about.html"
curl -s http://localhost:8080/about | grep 'translateAboutPage' | head -1

# Test 7: Verify JSON files have login translations
echo -e "\n✅ Test 7: Verificar que español tiene traducciones de login"
cat /Users/kiwinho/Proyectos/CodigoSH/web/static/i18n/es.json | grep -A3 '"login"' | head -5

# Test 8: Verify JSON files have about translations  
echo -e "\n✅ Test 8: Verificar que español tiene traducciones de about"
cat /Users/kiwinho/Proyectos/CodigoSH/web/static/i18n/es.json | grep -A3 '"about"' | head -5

echo -e "\n✅ Todas las pruebas completadas!"
