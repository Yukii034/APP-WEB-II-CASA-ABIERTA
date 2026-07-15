# Microservicio de actividad física

Registra, consulta, actualiza y elimina actividades físicas de adultos mayores.

## Rutas

- `GET /health`
- `GET /api/actividad-fisica`
- `POST /api/actividad-fisica`
- `GET /api/actividad-fisica/{id}`
- `PUT /api/actividad-fisica/{id}`
- `DELETE /api/actividad-fisica/{id}`

## JSON de ejemplo

```json
{
  "nombre_paciente": "Melanie Anchundia",
  "tipo_actividad": "Caminata",
  "duracion_minutos": 30,
  "intensidad": "moderada",
  "fecha": "2026-07-15",
  "estado": "completada",
  "observaciones": "Actividad realizada correctamente"
}
```
