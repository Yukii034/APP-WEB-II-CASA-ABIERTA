# Servicio: Reportes medicos

## Responsable

Equipo:

- Anthony Mendoza - Deimuz
- Holguin Nathaly Jasmin - NathalyLucas11
- Cedeño Geovanny Alexander - alex167j

Nota: el desarrollo fue colaborativo, pero el equipo centralizo los commits desde una sola computadora para evitar conflictos de ramas e integracion.

## Que hace este servicio

Genera reportes medicos simulados para adultos mayores. Resume informacion semanal sobre citas, alimentacion, alertas de salud y adherencia a medicinas.

Cuando `CITAS_URL` esta configurada, consulta el microservicio de citas medicas para calcular cuantas citas tiene cada paciente y cuantas estan completadas.

Los datos se manejan en memoria para la demostracion de casa abierta.

## Puerto

Este servicio corre internamente en el puerto **8080** dentro del contenedor.

Puerto expuesto al host: **8086**.

## Endpoints

| Metodo | Ruta | Descripcion |
|--------|------|-------------|
| GET | `/health` | Verifica que el servicio este vivo |
| GET | `/api/reportes-medicos/semanal` | Lista reportes semanales simulados |
| GET | `/api/reportes-medicos/resumen` | Muestra resumen general de reportes |
| GET | `/api/reportes-medicos/paciente/{id}` | Consulta reporte por paciente |

## Estructura del servicio

```txt
services/reportes-medicos/
├── main.go
├── models/
│   └── models.go
├── storage/
│   ├── storage.go
│   └── storage_test.go
├── handlers/
│   └── handlers.go
└── router/
    └── router.go
```

## Variables de entorno

| Variable | Descripcion | Ejemplo |
|----------|-------------|---------|
| PORT | Puerto interno del servicio | 8080 |
| CITAS_URL | URL interna del servicio de citas medicas | http://citas:8080 |

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
