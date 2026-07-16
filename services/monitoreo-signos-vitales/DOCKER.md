# Ejecución con Docker

Este microservicio expone el puerto `8083` y utiliza almacenamiento en memoria.

## Docker Compose

Desde este directorio:

```bash
docker compose up --build -d
```

El servicio queda disponible en `http://localhost:8083` y su estado puede verificarse en `http://localhost:8083/health`.

Para detenerlo:

```bash
docker compose down
```

## Docker sin Compose

```bash
docker build -t monitoreo-signos-vitales:local .
docker run --rm -p 8083:8083 --name monitoreo-signos-vitales monitoreo-signos-vitales:local
```

La imagen ejecuta el proceso con un usuario sin privilegios e incorpora una comprobación de salud sobre `/health`.
