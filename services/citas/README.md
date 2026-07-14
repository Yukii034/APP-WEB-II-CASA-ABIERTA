# Servicio: Citas medicas

## Responsable

Equipo: Deimuzh

## Que hace este servicio

Administra citas medicas simuladas para adultos mayores. Permite crear, consultar, actualizar, cancelar, confirmar y completar citas. Tambien expone recordatorios simulados, historial de cambios, metricas, pacientes y doctores de ejemplo.

Los datos se guardan en memoria, por lo que se reinician al reiniciar el servicio.

## Puerto

Este servicio corre internamente en el puerto **8080** dentro del contenedor.

Puerto expuesto al host: **8085**.

## Endpoints

| Metodo | Ruta | Descripcion |
|--------|------|-------------|
| GET | `/health` | Verifica que el servicio este vivo |
| GET | `/api/appointments` | Lista citas con filtros y paginacion |
| POST | `/api/appointments` | Crea una cita |
| GET | `/api/appointments/{id}` | Obtiene una cita por ID |
| PUT | `/api/appointments/{id}` | Actualiza fecha, hora, prioridad o motivo |
| DELETE | `/api/appointments/{id}` | Cancela una cita |
| PATCH | `/api/appointments/{id}/confirm` | Confirma una cita pendiente |
| PATCH | `/api/appointments/{id}/complete` | Marca una cita confirmada como completada |
| PATCH | `/api/appointments/{id}/notes` | Agrega notas medicas |
| GET | `/api/appointments/patient/{id}` | Lista citas de un paciente |
| GET | `/api/appointments/history/{id}` | Lista el historial de una cita |
| GET | `/api/appointments/reminders` | Lista recordatorios simulados |
| GET | `/api/appointments/metrics` | Muestra metricas del servicio |
| POST | `/api/appointments/recurring` | Crea citas recurrentes |
| GET | `/api/patients` | Lista pacientes simulados |
| GET | `/api/doctors` | Lista doctores simulados |

## Variables de entorno

| Variable | Descripcion | Ejemplo |
|----------|-------------|---------|
| PORT | Puerto interno del servicio | 8080 |
| API_KEY | Clave opcional para proteger endpoints | demo123 |

Si `API_KEY` no esta configurada, el servicio permite peticiones sin autenticacion.

## Ejemplo de creacion de cita

```bash
curl -X POST http://localhost:8085/api/appointments \
  -H "Content-Type: application/json" \
  -d '{"paciente_id":"P001","doctor_id":"D001","fecha":"2030-01-01","hora":"10:00","motivo":"Control general"}'
```

## Como correrlo solo

```bash
cd services/citas
go run main.go
```

## Como correrlo con Docker Compose

Desde la raiz del repositorio:

```bash
docker compose up --build citas
```

## Como probarlo

```bash
cd services/citas
go test ./...
```
