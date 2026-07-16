# Servicio: Contacto de emergencia

## Responsable

Equipo: Luis (Luisao)

## Que hace este servicio

Administra los contactos de emergencia (familiares/cuidadores) de cada paciente y gestiona las alertas de emergencia. Permite:

- Registrar, consultar, actualizar y eliminar contactos (nombre, telefono, parentesco, prioridad)
- Consultar los contactos de un paciente ordenados por prioridad de notificacion
- Crear una alerta de emergencia para un paciente, lo que simula el envio de una notificacion a sus contactos respetando el orden de prioridad
- Atender o cancelar una alerta activa
- Consultar el historial de una alerta y metricas generales del servicio

Los datos se guardan en memoria, por lo que se reinician al reiniciar el servicio.

## Puerto

Este servicio corre internamente en el puerto **8080** dentro del contenedor.

Puerto expuesto al host: **8086**.

## Endpoints

| Metodo | Ruta | Descripcion |
|--------|------|-------------|
| GET | `/health` | Verifica que el servicio este vivo |
| GET | `/api/contacts` | Lista todos los contactos (o filtra con `?paciente_id=`) |
| POST | `/api/contacts` | Crea un contacto de emergencia |
| GET | `/api/contacts/{id}` | Obtiene un contacto por ID |
| PUT | `/api/contacts/{id}` | Actualiza un contacto |
| DELETE | `/api/contacts/{id}` | Elimina un contacto |
| POST | `/api/alerts` | Crea una alerta y notifica (simulado) a los contactos del paciente |
| GET | `/api/alerts` | Lista todas las alertas |
| GET | `/api/alerts/{id}` | Obtiene una alerta por ID |
| PATCH | `/api/alerts/{id}/attend` | Marca una alerta activa como atendida |
| DELETE | `/api/alerts/{id}` | Cancela una alerta activa |
| GET | `/api/alerts/history/{id}` | Lista el historial de una alerta |
| GET | `/api/metrics` | Muestra metricas del servicio |

## Variables de entorno

| Variable | Descripcion | Ejemplo |
|----------|-------------|---------|
| PORT | Puerto interno del servicio | 8080 |
| API_KEY | Clave opcional para proteger endpoints | demo123 |

Si `API_KEY` no esta configurada, el servicio permite peticiones sin autenticacion.

## Ejemplo: crear un contacto

```bash
curl -X POST http://localhost:8086/api/contacts \
  -H "Content-Type: application/json" \
  -d '{"paciente_id":"P001","nombre":"Sofia Garcia","telefono":"555-0102","parentesco":"Hija","prioridad":1,"principal":true}'
```

## Ejemplo: activar una alerta de emergencia

```bash
curl -X POST http://localhost:8086/api/alerts \
  -H "Content-Type: application/json" \
  -d '{"paciente_id":"P001","mensaje":"Caida detectada en la sala","nivel":"critico"}'
```

## Como correrlo solo

```bash
cd services/contacto-emergencia
go run main.go
```

## Como correrlo con Docker Compose

Desde la raiz del repositorio:

```bash
docker compose up --build contacto-emergencia
```

## Como probarlo

```bash
cd services/contacto-emergencia
go test ./...
```