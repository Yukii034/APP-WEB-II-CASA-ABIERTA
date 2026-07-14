# Servicio: alimentación

## Responsable
Equipo: Eduardo Lopez, Pierina Peñaherrera, José Manuel Castillo, Néstor Gallegos

## Qué hace este servicio
Lleva un registro de las comidas del día (desayuno, almuerzo, cena) del
adulto mayor y avisa si alguna se saltó, es decir, si ya pasó la hora
límite esperada para esa comida y no fue registrada.

## Puerto
Este servicio corre internamente en el puerto **8080** (dentro del contenedor).
Puerto expuesto al host: **8084**

## Endpoints

| Método | Ruta                        | Descripción                                                          |
|--------|-----------------------------|------------------------------------------------------------------------|
| GET    | /health                     | Verifica que el servicio esté vivo                                     |
| GET    | /api/alimentacion           | Lista las comidas registradas hoy                                      |
| POST   | /api/alimentacion           | Registra una comida. Body: `{"tipo_comida":"almuerzo","descripcion":"sopa"}` |
| GET    | /api/alimentacion/resumen   | Estado de desayuno/almuerzo/cena de hoy y si alguna se saltó            |
| POST   | /api/alimentacion/reset     | Borra los registros de hoy (solo para pruebas/demo)                    |

`tipo_comida` acepta: `desayuno`, `almuerzo`, `cena` (también se puede usar
`merienda`, aunque no tiene hora límite configurada por defecto).

### Ejemplo de respuesta de `/api/alimentacion/resumen`

```json
{
  "comidas": [
    {"tipo_comida": "desayuno", "registrada": true,  "saltada": false, "hora_limite": "10:00"},
    {"tipo_comida": "almuerzo", "registrada": false, "saltada": true,  "hora_limite": "15:00"},
    {"tipo_comida": "cena",     "registrada": false, "saltada": false, "hora_limite": "21:00"}
  ],
  "comidas_hechas": 1,
  "comidas_total": 3,
  "hay_saltadas": true,
  "mensaje": "Hay una o más comidas que no se registraron a tiempo hoy."
}
```

## Variables de entorno

| Variable         | Descripción                                   | Ejemplo |
|-------------------|------------------------------------------------|---------|
| PORT              | Puerto interno del servicio                     | 8080    |
| DESAYUNO_HASTA    | Hora límite (HH:MM) para considerar el desayuno saltado | 10:00 |
| ALMUERZO_HASTA    | Hora límite (HH:MM) para el almuerzo             | 15:00   |
| CENA_HASTA        | Hora límite (HH:MM) para la cena                 | 21:00   |

## Cómo correrlo solo (sin docker-compose)

```bash
cd services/alimentacion
go run main.go
```

## Cómo correrlo con Docker

```bash
docker build -t alimentacion .
docker run -p 8084:8080 alimentacion
```

## Tests

```bash
go test ./... -v
```
