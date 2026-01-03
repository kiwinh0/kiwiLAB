# 🚀 CodigoSH

**CodigoSH** es un dashboard de marcadores minimalista, rápido y profesional, diseñado para centralizar el acceso a tus servicios autohospedados (Self-hosted) con una estética moderna, limpia y funcional.

---

## ✨ Características

### Interfaz & Diseño
* **UI Moderna:** Diseño glassmorphism con efectos de desenfoque y gradientes animados
* **Modo Oscuro/Claro:** Soporte nativo para ambos temas con persistencia
* **Responsive:** Totalmente adaptable a dispositivos móviles y desktop
* **Animaciones Fluidas:** Transiciones suaves y efectos visuales profesionales

### Funcionalidad
* **Asistente de Instalación:** Configuración inicial guiada e intuitiva
* **Multi-idioma:** Soporte para 20 idiomas principales
* **Gestión de Marcadores:** Añadir, editar, eliminar y reordenar con drag & drop
* **Búsqueda en Tiempo Real:** Filtra servicios instantáneamente
* **Iconos Integrados:** Librería de iconos de `selfhst`
* **Perfiles de Usuario:** Avatares personalizables y preferencias individuales

### Seguridad & Rendimiento
* **Autenticación JWT:** Sesiones seguras con tokens HTTP-only
* **Bcrypt:** Hashing seguro de contraseñas
* **SQLite:** Base de datos ligera e integrada
* **Aplicación Monolítica:** Single binary, sin dependencias externas
* **Bajo consumo:** ~20-30MB RAM en reposo, ~50MB bajo carga

---

## 🛠️ Instalación

### Requisitos Mínimos
* **OS:** Debian 11+, Ubuntu 20.04+ o compatible con Linux
* **CPU:** 1 core (x86_64 o ARM64)
* **RAM:** 128MB disponible (64MB en uso típico)
* **Almacenamiento:** 100MB disponible (50MB binario + 50MB espacio de trabajo)
* **Net:** Puerto 8080 libre (configurable)

### ⚡ Instalación Rápida

```bash
curl -sSL "https://raw.githubusercontent.com/kiwinh0/CodigoSH/main/scripts/install.sh" | sudo bash
```

El script instala automáticamente todo lo necesario. Al terminar:
- ✅ Dependencias del sistema
- ✅ Binario compilado de CodigoSH  
- ✅ Servicio systemd para autoarranque
- ✅ Base de datos inicializada

**Acceso:** http://IP_DEL_SERVIDOR:8080

El asistente de instalación te guiará para crear tu usuario y configurar la aplicación.

---

## 🔄 Actualizaciones Automáticas

CodigoSH incluye un sistema de actualizaciones automáticas que:

1. **Detecta Nuevas Versiones:** Verifica GitHub automáticamente cada 24 horas
2. **Notifica al Usuario:** Muestra badge y menú cuando hay actualización disponible
3. **Compila desde Fuente:** Descarga el código fuente y compila localmente en tu servidor
4. **Reemplaza Binario:** Actualiza automáticamente el ejecutable con backup
5. **Reinicia Servicio:** Reinicia automáticamente el servicio systemd

### Cómo Actualizar

* Click en el icono del usuario → "Actualización disponible"
* Click en el botón "Actualizar"
* Espera a que compile (30-60 segundos)
* El servicio se reiniciará automáticamente

### Requisitos para Actualizaciones

Las actualizaciones requieren que tu servidor tenga:
- **Git:** Instalado (incluido en install.sh)
- **Go 1.24+:** Instalado (incluido en install.sh)
- **Build Tools:** gcc, make (incluido en install.sh)

---


---

## 📝 Licencia

Este proyecto está disponible bajo licencia MIT.

## Créditos

Iconos proporcionados por [selfhst/icons](https://github.com/selfhst/icons).
