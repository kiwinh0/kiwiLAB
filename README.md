#  kiwiLAB

**kiwiLAB** es un dashboard de marcadores minimalista, rápido y profesional, diseñado para centralizar el acceso a tus servicios autohospedados (Self-hosted) con una estética moderna, limpia y funcional.

>  **ESTADO DEL PROYECTO:** kiwiLAB se encuentra actualmente en **fase activa de crecimiento y desarrollo**. Es un proyecto joven y en constante evolución, por lo que pueden producirse cambios significativos mientras avanzamos hacia una versión estable.

---

##  Características actuales
* **Interfaz Premium:** Diseño compacto con efectos de desenfoque y tipografía profesional.
* **Modo Oscuro Inteligente:** Soporte nativo para temas claros y oscuros.
* **Buscador en tiempo real:** Filtra tus servicios instantáneamente mientras escribes.
* **Gestión de Iconos:** Integración directa con la librería de iconos de `selfhst`.
* **Personalización:** Modales integrados para añadir, editar o eliminar marcadores sin salir de la web.
* **Orden Dinámico:** Reorganiza tus marcadores con tecnología Drag & Drop.

---

##  Requisitos Mínimos
Para garantizar el correcto funcionamiento de kiwiLAB de forma nativa:
* **Sistema Operativo:** Linux (Probado en Debian 12/13 y Ubuntu 22.04+).
* **Lenguaje:** Go (Golang) 1.21 o superior.
* **Dependencias de sistema:** `gcc`, `musl-dev` y `sqlite3` (necesarios para la compilación de la base de datos).
* **Hardware:** Muy ligero (64MB de RAM y 1 vCPU son suficientes para el binario).

---

##  Instalación y Despliegue

### 1. Instalación Automática (Script Bash)
Si utilizas una distribución basada en Debian o Ubuntu, puedes realizar la instalación completa (incluyendo dependencias y configuración del servicio de sistema `systemd`) con este comando:

```bash
curl -sSL [https://raw.githubusercontent.com/kiwinh0/kiwiLAB/main/install.sh](https://raw.githubusercontent.com/kiwinh0/kiwiLAB/main/install.sh) | sudo bash
```

### 2. Instalación Manual (Compilar)

####  Clonar el repositorio oficial
git clone [https://github.com/kiwinh0/kiwiLAB.git](https://github.com/kiwinh0/kiwiLAB.git)
cd kiwiLAB

####  Instalar dependencias de compilación (Debian/Ubuntu)
sudo apt update && sudo apt install -y build-essential gcc sqlite3

####  Compilar el binario habilitando CGO para el soporte de SQLite
CGO_ENABLED=1 go build -o kiwilab ./cmd/kiwilab/main.go

####  Ejecutar la aplicación
./kiwilab

### 3. La aplicación será accesible por defecto en http://localhost:8080

## Próximos Pasos (Roadmap)
[x] Dockerización: Implementación de Dockerfile y Docker Compose.

[ ] Seguridad: Sistema de autenticación de usuario (Login).

[ ] Organización: Categorización de marcadores por grupos o etiquetas.

[ ] Backup: Soporte para copias de seguridad automáticas de la base de datos.

## Créditos
Este proyecto utiliza la magnífica librería de iconos de selfhst/icons.

## kiwiLAB Project
