# Ejecución con Docker

Este microservicio escucha internamente en el puerto `8080` y utiliza almacenamiento en memoria.

## Docker Compose

Desde la raíz del repositorio:

```bash
docker compose up --build -d monitoreo-signos-vitales
```

El servicio queda disponible en `http://localhost:8083` y su estado puede verificarse en `http://localhost:8083/health`. Dentro de la red Docker, los demás servicios deben usar `http://monitoreo-signos-vitales:8080`.

Para detenerlo:

```bash
docker compose down
```

## Docker sin Compose

```bash
docker build -t monitoreo-signos-vitales:local services/monitoreo-signos-vitales
docker run --rm -p 8083:8080 --name monitoreo-signos-vitales monitoreo-signos-vitales:local
```

La imagen ejecuta el proceso con un usuario sin privilegios e incorpora una comprobación de salud sobre `/health`.
