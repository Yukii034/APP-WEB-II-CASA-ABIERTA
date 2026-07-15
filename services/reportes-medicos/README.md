# Servicio: Reportes medicos

## Responsable

Equipo: Deimuzh

## Que hace este servicio

Genera reportes medicos simulados para adultos mayores. Resume informacion semanal sobre citas, alimentacion, alertas de salud y adherencia a medicinas.

Los datos se manejan en memoria para la demostracion de casa abierta.

## Puerto

Este servicio corre internamente en el puerto **8080** dentro del contenedor.

Puerto expuesto al host: **8086**.

## Endpoints

| Metodo | Ruta | Descripcion |
|--------|------|-------------|
| GET | `/health` | Verifica que el servicio este vivo |
| GET | `/api/reportes` | Lista reportes semanales simulados |
| GET | `/api/reportes/semanal` | Lista reportes semanales simulados |
| GET | `/api/reportes/resumen` | Muestra resumen general de reportes |
| GET | `/api/reportes/paciente/{id}` | Consulta reporte por paciente |

## Variables de entorno

| Variable | Descripcion | Ejemplo |
|----------|-------------|---------|
| PORT | Puerto interno del servicio | 8080 |

## Como correrlo solo

```bash
cd services/reportes-medicos
go run main.go
```

## Como probarlo

```bash
cd services/reportes-medicos
go test ./...
```
