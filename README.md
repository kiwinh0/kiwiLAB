# 🚀 CodigoSH

**Versión:** v0.1.0-Beta

**CodigoSH** es un dashboard de marcadores minimalista, rápido y profesional, diseñado para centralizar el acceso a tus servicios autohospedados (Self-hosted) con una estética moderna, limpia y funcional.

> **ESTADO DEL PROYECTO:** CodigoSH v0.1.0-Beta es la primera versión beta pública. El proyecto está en constante evolución con nuevas características en desarrollo.

---

## 📁 Arquitectura del Proyecto
El proyecto sigue una estructura estándar de Go:

```
CodigoSH/
├── cmd/codigosH/          # Punto de entrada de la aplicación
├── internal/              # Código privado
│   ├── config/            # Configuración
│   ├── db/                # Capa de base de datos
│   ├── handlers/          # Handlers HTTP
│   ├── models/            # Estructuras de datos
│   └── middleware/        # Middlewares (auth, logging)
├── web/                   # Activos web
│   ├── static/            # CSS, JS, imágenes
│   │   ├── i18n/          # Archivos de internacionalización
│   │   └── uploads/       # Avatares de usuario
│   └── templates/         # Plantillas HTML
├── configs/               # Archivos de configuración YAML
├── scripts/               # Scripts de instalación
├── .github/workflows/     # CI/CD
└── Makefile               # Automatización de builds
```

---

## ✨ Características

### Interfaz & Diseño
* **UI Moderna:** Diseño glassmorphism con efectos de desenfoque y gradientes animados
* **Modo Oscuro:** Soporte nativo para temas claros y oscuros con persistencia
* **Responsive:** Totalmente adaptable a dispositivos móviles y desktop
* **Animaciones Fluidas:** Transiciones suaves y efectos visuales profesionales

### Funcionalidad
* **Wizard de Setup:** Configuración inicial guiada al primer uso
* **Multi-idioma:** Soporte para 20 idiomas principales del mundo
* **Gestión de Marcadores:** Añadir, editar, eliminar y reordenar con drag & drop
* **Búsqueda en Tiempo Real:** Filtra servicios instantáneamente
* **Iconos Integrados:** Librería de iconos de `selfhst`
* **Perfiles de Usuario:** Avatares personalizables y preferencias individuales

### Seguridad & Backend
* **Autenticación JWT:** Sesiones seguras con tokens HTTP-only
* **Bcrypt:** Hashing seguro de contraseñas
* **SQLite:** Base de datos ligera y eficiente
* **Logging Estructurado:** Monitoreo con Logrus
* **Middleware:** Autenticación y logging centralizado

---

## 🛠️ Instalación

### Requisitos
* Go 1.21+
* SQLite3
* GCC (para CGO)



### Instalación Automática (Debian/Ubuntu)
Script completo con configuración de systemd:

```bash
curl -sSL "https://raw.githubusercontent.com/kiwinh0/CodigoSH/main/scripts/install.sh" | sudo bash
```

### Instalación Manual

```bash
# Clonar repositorio
git clone https://github.com/kiwinh0/CodigoSH.git
cd CodigoSH

# Instalar dependencias (Debian/Ubuntu)
sudo apt update && apt install -y build-essential gcc sqlite3

# Instalar dependencias Go
make deps

# Compilar
make build

# Ejecutar
make run
```

### Docker

```bash
# Con docker-compose
docker-compose up -d

# O con Docker directo
docker build -t codigosh .
docker run -p 8080:8080 -v ./codigosH.db:/root/codigosH.db codigosh
```

La aplicación estará disponible en `http://localhost:8080`

---

## 🎯 Primer Uso

1. **Wizard de Setup:** Al acceder por primera vez, se mostrará un asistente de configuración
2. **Selecciona tu idioma:** Elige entre 20 idiomas disponibles
3. **Crea tu usuario admin:** Username, contraseña, tema y avatar opcional
4. **¡Listo!** Accede y comienza a agregar tus marcadores

### Gestión de Usuarios

Para agregar más usuarios manualmente:
```sql
sqlite3 codigosH.db
INSERT INTO users (username, password, role, language, theme) 
VALUES ('usuario', '$2a$10$...hash...', 'user', 'es', 'dark');
```

---

## ⚙️ Configuración

Edita `configs/config.yaml`:

```yaml
server:
  host: "0.0.0.0"
  port: "8080"

database:
  path: "./codigosH.db"

logging:
  level: "info"  # debug, info, warn, error
```

---

## 🚀 Desarrollo

```bash
# Compilar
make build

# Ejecutar tests
make test

# Ejecutar en desarrollo
make run

# Limpiar binarios
make clean
```

---

## 📝 Roadmap

- [x] Autenticación JWT
- [x] Multi-idioma (20 idiomas)
- [x] Wizard de setup inicial
- [x] Glassmorphism UI
- [x] Docker & Docker Compose
- [ ] API REST completa
- [ ] Importar/Exportar marcadores
- [ ] Temas personalizados
- [ ] Dashboard de estadísticas
- [ ] Integración con servicios externos

[x] Estructura Profesional: Separación de capas y configuración externa.

[x] Autenticación: Sistema de login seguro con JWT y bcrypt.

[ ] Seguridad: Implementar HTTPS y rate limiting.

[ ] Organización: Categorización de marcadores por grupos o etiquetas.

[ ] Backup: Soporte para copias de seguridad automáticas de la base de datos.

---

## Créditos
Este proyecto utiliza la magnífica librería de iconos de selfhst/icons.
