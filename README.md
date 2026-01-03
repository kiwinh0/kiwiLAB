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
* **Bajo consumo:** ~50MB RAM en operación normal

---

## 🛠️ Instalación

### Requisitos Mínimos
* **OS:** Debian 11+, Ubuntu 20.04+ o compatible con Linux
* **CPU:** 1 core mínimo (2+ recomendado)
* **RAM:** 256MB mínimo (512MB recomendado)
* **Almacenamiento:** 500MB disponible
* **Net:** Acceso a internet (solo para descargar el script)

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

## ⚙️ Configuración

Edita `~/.codigosh/config.yaml` si necesitas cambios avanzados:

```yaml
server:
  host: "0.0.0.0"
  port: "8080"
```

---

## 📝 Licencia

Este proyecto está disponible bajo licencia MIT.

## Créditos

Iconos proporcionados por [selfhst/icons](https://github.com/selfhst/icons).
