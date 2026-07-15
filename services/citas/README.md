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
| GET | `/api/cita-medica` | Lista citas con filtros y paginacion |
| POST | `/api/cita-medica` | Crea una cita |
| GET | `/api/cita-medica/{id}` | Obtiene una cita por ID |
| PUT | `/api/cita-medica/{id}` | Actualiza fecha, hora, prioridad o motivo |
| DELETE | `/api/cita-medica/{id}` | Cancela una cita |
| PATCH | `/api/cita-medica/{id}/confirmar` | Confirma una cita pendiente |
| PATCH | `/api/cita-medica/{id}/completar` | Marca una cita confirmada como completada |
| PATCH | `/api/cita-medica/{id}/notas` | Agrega notas medicas |
| GET | `/api/cita-medica/{id}/detalle` | Consulta la cita con paciente, doctor e informacion de salud |
| GET | `/api/cita-medica/paciente/{id}` | Lista citas de un paciente |
| GET | `/api/cita-medica/historial/{id}` | Lista el historial de una cita |
| GET | `/api/cita-medica/recordatorios` | Lista recordatorios simulados |
| GET | `/api/cita-medica/metricas` | Muestra metricas del servicio |
| POST | `/api/cita-medica/recurrentes` | Crea citas recurrentes |
| GET | `/api/cita-medica/pacientes` | Lista pacientes simulados |
| GET | `/api/cita-medica/doctores` | Lista doctores simulados |

## Variables de entorno

| Variable | Descripcion | Ejemplo |
|----------|-------------|---------|
| PORT | Puerto interno del servicio | 8080 |
| API_KEY | Clave opcional para proteger endpoints | demo123 |
| INFORMACION_SALUD_URL | URL interna del servicio de informacion-salud | http://informacion-salud:8080 |

Si `API_KEY` no esta configurada, el servicio permite peticiones sin autenticacion.

El endpoint de detalle usa un mapeo interno para relacionar pacientes de citas con fichas de informacion-salud:

| Paciente citas | Ficha informacion-salud |
|----------------|--------------------------|
| P001 | 1 |
| P002 | 2 |
| P003 | 3 |

## Ejemplo de creacion de cita

```bash
curl -X POST http://localhost:8085/api/cita-medica \
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
