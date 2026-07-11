# Servicio: Recordatorio de medicamentos

## Responsable
Equipo: [nombres]

## Qué hace este servicio
Guarda horarios y "avisa" (puede imprimir en consola o mandar notificación simulada) cuándo tomar cada medicamento

## Puerto
Este servicio corre internamente en el puerto **8080** (dentro del contenedor).
Puerto expuesto al host: **[completar, ej: 8082]**

## Endpoints

| Método | Ruta              | Descripción                  |
|--------|-------------------|-------------------------------|
| GET    | /health           | Verifica que el servicio esté vivo |
| GET    | /api/items        | [reemplazar con endpoint real] |

## Variables de entorno

| Variable | Descripción | Ejemplo |
|----------|-------------|---------|
| PORT     | Puerto interno del servicio | 8080 |

Si este servicio necesita llamar a otro, agregar aquí la variable, ej:
| OTRO_SERVICIO_URL | URL del servicio X | http://otro-servicio:8080 |

## Cómo correrlo solo (sin docker-compose)

```bash
cd services/[nombre-del-servicio]
go run main.go
```

## Cómo correrlo con Docker

```bash
docker build -t [nombre-del-servicio] .
docker run -p 8082:8080 [nombre-del-servicio]
```