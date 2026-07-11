# Servicio: [nombre del servicio]

## Responsable
Equipo: [nombres]

## Qué hace este servicio
Microservicio encargado de gestionar la información de salud del adulto mayor, almacenando datos como diagnósticos, alergias, enfermedades crónicas y antecedentes médicos. Proporciona una API REST para consultar y actualizar esta información, permitiendo que otros microservicios accedan a ella cuando sea necesario.

## Puerto
Este servicio corre internamente en el puerto **8080** (dentro del contenedor).
Puerto expuesto al host: **8082**

## Endpoints

| Método | Ruta              | Descripción                  |
|--------|-------------------|-------------------------------|
| GET    | /health           | Verifica que el servicio esté vivo |

## Variables de entorno

| Variable | Descripción | Ejemplo |
|----------|-------------|---------|
| PORT     | Puerto interno del servicio | 8082 |

Si este servicio necesita llamar a otro, agregar aquí la variable, ej:
| OTRO_SERVICIO_URL | URL del servicio X | http://otro-servicio:8080 |

## Cómo correrlo solo (sin docker-compose)

```bash
cd services/informacion-salud
go run main.go
```

## Cómo correrlo con Docker

```bash
docker build -t [informacion-salud] .
docker run -p 8082:8080 [informacion-salud]
```