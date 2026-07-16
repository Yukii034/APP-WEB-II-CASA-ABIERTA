# Servicio: Recordatorio de Medicamentos

## Responsables

Equipo:

- Manuel Intriago
- Madelyn Zambrano
- Michelle Salazar

## Qué hace este servicio

Microservicio encargado de registrar los medicamentos de los
adultos mayores, incluyendo paciente, dosis, frecuencia y horario.

Cuando llega la hora programada, el sistema genera una alerta
simulada en la consola. También existe un endpoint para verificar
manualmente una hora durante la demostración.

## Puerto

- Puerto interno del contenedor: `8080`
- Puerto expuesto al host: `8088`

## Endpoints

| Método | Ruta | Descripción |
|---|---|---|
| GET | `/health` | Verifica que el servicio esté activo |
| GET | `/api/recordatorio-medicamentos` | Lista los recordatorios |
| POST | `/api/recordatorio-medicamentos` | Crea un recordatorio |
| POST | `/api/recordatorio-medicamentos/verificar` | Simula una hora |
| GET | `/api/recordatorio-medicamentos/{id}` | Consulta por ID |
| PUT | `/api/recordatorio-medicamentos/{id}` | Actualiza parcialmente |
| DELETE | `/api/recordatorio-medicamentos/{id}` | Elimina |
| PATCH | `/api/recordatorio-medicamentos/{id}/estado` | Activa o desactiva |

## Crear un recordatorio

```json
{
  "adulto_mayor_id": "AM-003",
  "nombre_paciente": "Carmen Torres",
  "medicamento": "Aspirina",
  "dosis": "100 mg",
  "hora": "09:00",
  "frecuencia": "diaria",
  "activo": true
}