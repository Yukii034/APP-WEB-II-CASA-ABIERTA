# Servicio: actividad física

## Responsable

Equipo: Melanie Anchundia, Lisbeth Aray 

## Qué hace este servicio

Lleva un registro de las actividades físicas realizadas por el adulto mayor,
permitiendo almacenar información como el tipo de actividad, duración,
intensidad, fecha y observaciones. Además, calcula automáticamente una
estimación de las calorías quemadas durante la actividad.

# Mejoras al servicio
Implementación de Middleware de Logging
Integracion de un sistema de registro de eventos (Logging) para mejorar la observabilidad del microservicio.

¿Qué hace?: Registra automáticamente cada petición HTTP que recibe el servicio, capturando información clave como el método, la ruta, el código de estado (status code) y el tiempo de respuesta.

Implementación: Se creó un middleware personalizado en internal/middleware/ que utiliza el patrón decorador para interceptar el ResponseWriter y extraer el código de respuesta (ej. 200, 201, 404).

Cambios en main.go: Se actualizaron las definiciones de las rutas para envolver los manejadores (handlers) con esta función de Logger, garantizando que todas las peticiones sean auditadas sin duplicar lógica en los controladores.

Beneficios:

Facilita la depuración (debugging) rápida al ver qué peticiones fallan directamente en la consola.

Proporciona una trazabilidad profesional del uso de la API.

Centraliza el manejo de logs para facilitar futuros cambios o integraciones.

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
