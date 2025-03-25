# Solicitudes Frontend

Este repositorio contiene el frontend del sistema de solicitudes.

## ⚙️ Requisitos

- Node.js
- Vite (opcional)

## 🔧 Configuración

Antes de ejecutar el proyecto, debés crear un archivo `.env` en la carpeta `Solicitudes-Front/` con las siguientes variables:

```env
VITE_BASE_USERS_API=
VITE_BASE_LOGIN_API=
VITE_BASE_MAILING_API=
VITE_BASE_REQUEST_API=
VITE_BASE_MANAGER_API=
```

Podés usar el archivo `.env.example` como guía:

```bash
cp Solicitudes-Front/.env.example Solicitudes-Front/.env
```

Completá los valores según tu entorno local.

## 🚀 Ejecución

```bash
cd Solicitudes-Front
npm install
npm run dev
```

---

Este repositorio está listo para hacerse público, sin exponer datos sensibles.
