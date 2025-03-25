# Solicitudes Backend

Este repositorio contiene el backend del sistema de solicitudes.

## ⚙️ Requisitos

- Node.js
- PostgreSQL
- Docker (opcional)

## 🔧 Configuración

Antes de ejecutar el proyecto, debés crear un archivo `.env` en la carpeta `Solicitudes.Api/` con las siguientes variables de entorno:

```env
API_VERSION=
RSA_PUBLIC_KEY="YOUR_PUBLIC_KEY"
AUTH_DELVE_PORT=
AUTH_WEB_SERVER_PORT=
POSTGRES_PASSWORD=
```

Podés usar el archivo `.env.example` como guía:

```bash
cp Solicitudes.Api/.env.example Solicitudes.Api/.env
```

Completá los valores según tu entorno local o de producción.

## 🚀 Ejecución

```bash
cd Solicitudes.Api
npm install
npm start
```

---

Este repositorio está preparado para uso público, sin exponer información sensible.
