# 🚀 Solicitudes App

Plataforma completa para la gestión de solicitudes ciudadanas. Incluye backend y frontend desacoplados, listos para desplegar en cualquier municipio o institución.

---

## 🧩 Estructura del Proyecto

```
Solicitudes-App/
├── Solicitudes-Backend/
│   └── Solicitudes.Api/       # Backend en Node.js
├── Solicitudes-Frontend/
│   └── Solicitudes-Front/     # Frontend en Vite + React
└── README.md                  # Este archivo
```

---

## ⚙️ Tecnologías Utilizadas

- **Frontend:** React + Vite + JavaScript
- **Backend:** Node.js + Express
- **Base de Datos:** PostgreSQL
- **APIs:** REST
- **Ambiente:** `.env` para configuración local

---

## 🔧 Configuración Inicial

1. Cloná el repo o descargalo como ZIP.
2. Configurá tus archivos `.env` en ambos proyectos usando los `.env.example` de referencia.
3. Instalá dependencias:

### 🖥️ Backend

```bash
cd Solicitudes-Backend/Solicitudes.Api
npm install
npm start
```

### 🌐 Frontend

```bash
cd Solicitudes-Frontend/Solicitudes-Front
npm install
npm run dev
```

---

## 📡 Estructura del Backend

El backend expone endpoints REST agrupados por recursos:

- `/api/v1/users` → Gestión de usuarios
- `/api/v1/login` → Autenticación
- `/api/v1/mailing` → Envío de mails
- `/api/v1/requests` → Solicitudes ciudadanas
- `/api/v1/manager` → Gestión interna

---

## 🖼️ Frontend

La interfaz permite:
- Login de usuarios
- Carga y seguimiento de solicitudes
- Panel de administración para operadores

Variables como las URLs de APIs están en `VITE_*` y se configuran por `.env`.

---

## 🛡️ Seguridad

- Las credenciales y claves fueron **removidas** del código.
- Usá los `.env.example` para tu configuración sin exponer información sensible.

---

## 🌍 Licencia

Este proyecto está publicado como **open source**. ¡Podés usarlo, adaptarlo y mejorarlo!

---

Hecho con ❤️ para modernizar la gestión pública.
