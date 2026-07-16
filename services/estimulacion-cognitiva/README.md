# Servicio: Estimulación cognitiva

## Responsable
Equipo: 

López Cedeño Bryan Steeven
Sornoza Leon Isaac Arturo
Garcia Flores Eduardo Antonio

## Qué hace este servicio
Registra los ejercicios de estimulación cognitiva que completa el adulto
mayor (memoria, trivia, sopa de letras, etc.) y calcula cuántos días han
pasado desde el último ejercicio realizado. Si pasan **2 o más días** sin
ningún ejercicio registrado, activa una alerta de inactividad — la idea es
detectar señales tempranas de desinterés o deterioro cognitivo.

Al arrancar el servicio se precargan datos de ejemplo (seed) que dejan la
**alerta ya activa a propósito**: así se puede mostrar en vivo cómo se
resuelve la alerta al registrar un nuevo ejercicio, sin depender de que la
demo funcione "a la primera intentona" con datos improvisados.

## Puerto
Este servicio corre internamente en el puerto **8080** (dentro del contenedor).
Puerto expuesto al host: **8095** (verificar que no choque con otro equipo
antes de agregarlo a `docker-compose.yml` y actualizar la tabla de puertos
del README general).

## Endpoints

| Método | Ruta                        | Descripción                                                     |
|--------|-----------------------------|----------------------------------------------------------------|
| GET    | /health                     | Verifica que el servicio esté vivo                              |
| GET    | /api/ejercicios             | Lista todos los ejercicios registrados                           |
| POST   | /api/ejercicios             | Registra un ejercicio completado. Body: `{"tipo":"memoria"}`     |
| GET    | /api/ejercicios/resumen     | Ejercicios de hoy, días desde el último, y si hay alerta         |
| POST   | /api/ejercicios/reset       | Borra los ejercicios registrados (solo para pruebas/demo)        |

`tipo` es libre (ej. `memoria`, `trivia`, `sopa_letras`, `rompecabezas`), solo
no puede venir vacío.

### Ejemplo de respuesta de `/api/ejercicios/resumen`

```json
{
  "ejercicios": [
    {"id": "1", "tipo": "trivia", "fecha": "2026-07-11T10:00:00Z"},
    {"id": "2", "tipo": "memoria", "fecha": "2026-07-13T10:00:00Z"}
  ],
  "total": 2,
  "ejercicios_hoy": 0,
  "ultimo_ejercicio": {"id": "2", "tipo": "memoria", "fecha": "2026-07-13T10:00:00Z"},
  "dias_desde_ultimo": 2,
  "hay_alerta": true,
  "mensaje": "Han pasado varios días sin actividad cognitiva registrada."
}
```

## Variables de entorno

| Variable    | Descripción                                                                       | Ejemplo |
|-------------|---------------------------------------------------------------------------------------|---------|
| PORT        | Puerto interno del servicio                                                            | 8080    |
| SEED_DATOS  | Si es `false`, desactiva la precarga de datos de ejemplo al arrancar. Por defecto está activada. | false   |

## Persistencia

Los datos se guardan **en memoria** (protegidos con un mutex, dentro de
`storage/`), siguiendo la limitación conocida del proyecto (ver
`docs/arquitectura.md`). Se pierden si el contenedor se reinicia.

## Estructura del servicio

```
services/estimulacion-cognitiva/
├── main.go                    # wiring: arma storage -> router y levanta el servidor
├── models/ejercicio.go        # structs Ejercicio y Resumen
├── storage/storage.go         # persistencia en memoria + lógica de negocio (validaciones, cálculo de días de inactividad)
├── storage/storage_test.go    # tests del storage
├── handlers/handlers.go       # handlers HTTP (decodifican/codifican JSON, llaman al storage)
├── router/router.go           # arma las rutas y aplica los middlewares
├── middleware/middleware.go   # logging de requests + recuperación de panics
├── logger/logger.go           # logger simple con niveles Info/Error
├── go.mod
└── Dockerfile
```

## Cómo correrlo solo (sin docker-compose)

```bash
cd services/estimulacion-cognitiva
go run main.go
```

## Cómo correrlo con Docker

```bash
docker build -t estimulacion-cognitiva .
docker run -p 8090:8080 estimulacion-cognitiva
```

## Tests

```bash
go test ./... -v
```
