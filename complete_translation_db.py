#!/usr/bin/env python3
"""
Script comprehensivo para traducir TODO el proyecto CodigoSH
Extrae todos los textos de cada página y los mapea a traducciones
"""

import json
import os
import re

# Directorio de traducciones
i18n_dir = '/Users/kiwinho/Proyectos/CodigoSH/web/static/i18n'

# Diccionario completo de traducciones por página
translations_db = {
    'login': {
        'en': {
            'CodigoSH': 'CodigoSH',
            'Tu centro de comando digital': 'Your digital command center',
            'Usuario': 'Username',
            'Ingresa tu usuario': 'Enter your username',
            'Contraseña': 'Password',
            'Ingresa tu contraseña': 'Enter your password',
            'Mantenerme conectado': 'Keep me signed in',
            'Sesión activa por 30 días': 'Active session for 30 days',
            'Iniciar Sesión': 'Sign In',
        },
        'es': {
            'CodigoSH': 'CodigoSH',
            'Tu centro de comando digital': 'Tu centro de comando digital',
            'Usuario': 'Usuario',
            'Ingresa tu usuario': 'Ingresa tu usuario',
            'Contraseña': 'Contraseña',
            'Ingresa tu contraseña': 'Ingresa tu contraseña',
            'Mantenerme conectado': 'Mantenerme conectado',
            'Sesión activa por 30 días': 'Sesión activa por 30 días',
            'Iniciar Sesión': 'Iniciar Sesión',
        },
        'fr': {
            'CodigoSH': 'CodigoSH',
            'Tu centro de comando digital': 'Votre centre de commandement numérique',
            'Usuario': 'Nom d\'utilisateur',
            'Ingresa tu usuario': 'Entrez votre nom d\'utilisateur',
            'Contraseña': 'Mot de passe',
            'Ingresa tu contraseña': 'Entrez votre mot de passe',
            'Mantenerme conectado': 'Me garder connecté',
            'Sesión activa por 30 días': 'Session active pendant 30 jours',
            'Iniciar Sesión': 'Se connecter',
        },
        'de': {
            'CodigoSH': 'CodigoSH',
            'Tu centro de comando digital': 'Dein digitales Kommandozentrum',
            'Usuario': 'Benutzername',
            'Ingresa tu usuario': 'Geben Sie Ihren Benutzernamen ein',
            'Contraseña': 'Passwort',
            'Ingresa tu contraseña': 'Geben Sie Ihr Passwort ein',
            'Mantenerme conectado': 'Anmeldedaten speichern',
            'Sesión activa por 30 días': 'Aktive Sitzung für 30 Tage',
            'Iniciar Sesión': 'Anmelden',
        },
        'it': {
            'CodigoSH': 'CodigoSH',
            'Tu centro de comando digital': 'Il tuo centro di comando digitale',
            'Usuario': 'Nome utente',
            'Ingresa tu usuario': 'Inserisci il tuo nome utente',
            'Contraseña': 'Password',
            'Ingresa tu contraseña': 'Inserisci la tua password',
            'Mantenerme conectado': 'Mantienimi connesso',
            'Sesión activa por 30 días': 'Sessione attiva per 30 giorni',
            'Iniciar Sesión': 'Accedi',
        },
        'pt': {
            'CodigoSH': 'CodigoSH',
            'Tu centro de comando digital': 'Seu centro de comando digital',
            'Usuario': 'Nome de usuário',
            'Ingresa tu usuario': 'Digite seu nome de usuário',
            'Contraseña': 'Senha',
            'Ingresa tu contraseña': 'Digite sua senha',
            'Mantenerme conectado': 'Mantenha-me conectado',
            'Sesión activa por 30 días': 'Sessão ativa por 30 dias',
            'Iniciar Sesión': 'Entrar',
        },
        'ru': {
            'CodigoSH': 'CodigoSH',
            'Tu centro de comando digital': 'Ваш центр цифрового управления',
            'Usuario': 'Имя пользователя',
            'Ingresa tu usuario': 'Введите ваше имя пользователя',
            'Contraseña': 'Пароль',
            'Ingresa tu contraseña': 'Введите ваш пароль',
            'Mantenerme conectado': 'Держать меня в системе',
            'Sesión activa por 30 días': 'Активная сессия на 30 дней',
            'Iniciar Sesión': 'Войти',
        },
        'zh': {
            'CodigoSH': 'CodigoSH',
            'Tu centro de comando digital': '您的数字指挥中心',
            'Usuario': '用户名',
            'Ingresa tu usuario': '输入您的用户名',
            'Contraseña': '密码',
            'Ingresa tu contraseña': '输入您的密码',
            'Mantenerme conectado': '保持登录',
            'Sesión activa por 30 días': '30天内保持活跃会话',
            'Iniciar Sesión': '登录',
        },
        'ja': {
            'CodigoSH': 'CodigoSH',
            'Tu centro de comando digital': 'あなたのデジタルコマンドセンター',
            'Usuario': 'ユーザー名',
            'Ingresa tu usuario': 'ユーザー名を入力',
            'Contraseña': 'パスワード',
            'Ingresa tu contraseña': 'パスワードを入力',
            'Mantenerme conectado': 'ログインしたままにする',
            'Sesión activa por 30 días': '30日間有効なセッション',
            'Iniciar Sesión': 'サインイン',
        },
        'ko': {
            'CodigoSH': 'CodigoSH',
            'Tu centro de comando digital': '당신의 디지털 명령 센터',
            'Usuario': '사용자명',
            'Ingresa tu usuario': '사용자명을 입력하세요',
            'Contraseña': '비밀번호',
            'Ingresa tu contraseña': '비밀번호를 입력하세요',
            'Mantenerme conectado': '로그인 유지',
            'Sesión activa por 30 días': '30일 동안 활성 세션',
            'Iniciar Sesión': '로그인',
        },
    },
    'dashboard': {
        'en': {
            'Buscar servicios...': 'Search services...',
            'Agregar marcador': 'Add bookmark',
            'Agregar grupo': 'Add group',
            'Preferencias': 'Settings',
            'Acerca de': 'About',
            'Cerrar sesión': 'Logout',
        },
        'es': {
            'Buscar servicios...': 'Buscar servicios...',
            'Agregar marcador': 'Agregar marcador',
            'Agregar grupo': 'Agregar grupo',
            'Preferencias': 'Preferencias',
            'Acerca de': 'Acerca de',
            'Cerrar sesión': 'Cerrar sesión',
        },
    },
    'about': {
        'en': {
            'Versión 1.0.0': 'Version 1.0.0',
            'Una aplicación web moderna para gestionar tus marcadores y enlaces favoritos. Organiza, personaliza y accede rápidamente a tus recursos más importantes.': 'A modern web application to manage your bookmarks and favorite links. Organize, customize and quickly access your most important resources.',
            'Desarrollado con ❤️ usando Go y modern web technologies': 'Developed with ❤️ using Go and modern web technologies',
            '© 2026 CodigoSH. Todos los derechos reservados.': '© 2026 CodigoSH. All rights reserved.',
            'Ver en GitHub': 'View on GitHub',
        },
        'es': {
            'Versión 1.0.0': 'Versión 1.0.0',
            'Una aplicación web moderna para gestionar tus marcadores y enlaces favoritos. Organiza, personaliza y accede rápidamente a tus recursos más importantes.': 'Una aplicación web moderna para gestionar tus marcadores y enlaces favoritos. Organiza, personaliza y accede rápidamente a tus recursos más importantes.',
            'Desarrollado con ❤️ usando Go y modern web technologies': 'Desarrollado con ❤️ usando Go y modern web technologies',
            '© 2026 CodigoSH. Todos los derechos reservados.': '© 2026 CodigoSH. Todos los derechos reservados.',
            'Ver en GitHub': 'Ver en GitHub',
        },
    }
}

# Agregar traducciones de dashboard completas a todos los idiomas
dashboard_full = {
    'en': {
        'Buscar servicios...': 'Search services...',
        'Agregar marcador': 'Add bookmark',
        'Agregar grupo': 'Add group',
        'Preferencias': 'Settings',
        'Acerca de': 'About',
        'Cerrar sesión': 'Logout',
        'Dashboard': 'Dashboard',
        'Grupo sin nombre': 'Unnamed group',
        'Configurar Marcador': 'Configure Bookmark',
        'Nombre': 'Name',
        'Ej: Plex': 'E.g.: Plex',
        'URL': 'URL',
        'https://...': 'https://...',
        'Icono': 'Icon',
        'Buscar icono...': 'Search icon...',
        'Online': 'Online',
        'Guardar Cambios': 'Save Changes',
        'Cancelar': 'Cancel',
        '¿Eliminar marcador?': 'Delete bookmark?',
        'Estás a punto de borrar': 'You are about to delete',
        'Confirmar': 'Confirm',
    },
    'es': {
        'Buscar servicios...': 'Buscar servicios...',
        'Agregar marcador': 'Agregar marcador',
        'Agregar grupo': 'Agregar grupo',
        'Preferencias': 'Preferencias',
        'Acerca de': 'Acerca de',
        'Cerrar sesión': 'Cerrar sesión',
    },
    'fr': {
        'Buscar servicios...': 'Chercher des services...',
        'Agregar marcador': 'Ajouter un signet',
        'Agregar grupo': 'Ajouter un groupe',
        'Preferencias': 'Paramètres',
        'Acerca de': 'À propos',
        'Cerrar sesión': 'Déconnexion',
    },
    'de': {
        'Buscar servicios...': 'Dienste durchsuchen...',
        'Agregar marcador': 'Lesezeichen hinzufügen',
        'Agregar grupo': 'Gruppe hinzufügen',
        'Preferencias': 'Einstellungen',
        'Acerca de': 'Über',
        'Cerrar sesión': 'Abmelden',
    },
    'it': {
        'Buscar servicios...': 'Cerca servizi...',
        'Agregar marcador': 'Aggiungi segnalibro',
        'Agregar grupo': 'Aggiungi gruppo',
        'Preferencias': 'Impostazioni',
        'Acerca de': 'Informazioni',
        'Cerrar sesión': 'Esci',
    },
}

# Print useful info
print("📚 Diccionario de traducciones cargado")
print(f"✅ Idiomas configurados: {list(translations_db.keys())}")
for page in translations_db:
    print(f"\n📄 {page.upper()}:")
    for lang in translations_db[page]:
        print(f"  - {lang}: {len(translations_db[page][lang])} traducciones")
