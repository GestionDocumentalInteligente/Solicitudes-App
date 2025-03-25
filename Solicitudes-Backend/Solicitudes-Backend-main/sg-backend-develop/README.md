### Para usar el proyecto:

```bash
# Desarrollo
make sg_backend-dev-up
make sg_backend-dev-build
make sg_backend-dev-down

# Producción
make sg_backend-prod-up
make sg_backend-prod-build
make sg_backend-prod-down
```

### Configuración de Git Flow

Para configurar Git Flow en este proyecto, ejecuta el siguiente comando en la raíz del repositorio:

````bash
git flow init
```

### Configuración de Git Flow
```bash
# AWS
docker-compose --profile sg_backend -f config/docker-compose.prod.yml up --build -d
````
