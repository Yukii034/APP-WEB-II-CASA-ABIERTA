# Servicio: actividad física

## Responsable

Equipo: Melanie Anchundia, Lisbeth Aray 

## Qué hace este servicio

Lleva un registro de las actividades físicas realizadas por el adulto mayor,
permitiendo almacenar información como el tipo de actividad, duración,
intensidad, fecha y observaciones. Además, calcula automáticamente una
estimación de las calorías quemadas durante la actividad.

## Puerto

Este servicio corre internamente en el puerto **8080** (dentro del contenedor).

Puerto expuesto al host: **8090**

## Endpoints

| Método | Ruta                       | Descripción                                     |
| ------ | -------------------------- | ----------------------------------------------- |
| GET    | /health                    | Verifica que el servicio esté vivo              |
| GET    | /api/actividad-fisica      | Lista todas las actividades registradas         |
| POST   | /api/actividad-fisica      | Registra una nueva actividad física. Body abajo |
| GET    | /api/actividad-fisica/{id} | Obtiene una actividad por su ID                 |
| PUT    | /api/actividad-fisica/{id} | Actualiza una actividad existente               |
| DELETE | /api/actividad-fisica/{id} | Elimina una actividad registrada                |

`intensidad` acepta: `baja`, `moderada`, `alta`.

`estado` acepta: `pendiente`, `completada`, `cancelada`.

### POST /api/actividad-fisica

```json
{
  "nombre_paciente": "María López",
  "tipo_actividad": "Caminata",
  "duracion_minutos": 30,
  "intensidad": "moderada",
  "fecha": "2026-07-15",
  "estado": "completada",
  "observaciones": "Actividad realizada sin inconvenientes"
}
```

### Ejemplo de respuesta

```json
{
  "id": "1",
  "nombre_paciente": "María López",
  "tipo_actividad": "Caminata",
  "duracion_minutos": 30,
  "intensidad": "moderada",
  "fecha": "2026-07-15",
  "estado": "completada",
  "observaciones": "Actividad realizada sin inconvenientes",
  "calorias_estimadas": 150
}
```

`calorias_estimadas` es calculado automáticamente por el sistema según la
duración e intensidad de la actividad física.

## Variables de entorno

| Variable | Descripción                 | Ejemplo |
| -------- | --------------------------- | ------- |
| PORT     | Puerto interno del servicio | 8080    |

## Cómo correrlo solo (sin docker-compose)

```bash
cd services/actividad-fisica
go run main.go
```

## Cómo correrlo con Docker

```bash
docker build -t actividad-fisica .
docker run -p 8088:8080 actividad-fisica
```

## Tests

```bash
go test ./... -v
```
