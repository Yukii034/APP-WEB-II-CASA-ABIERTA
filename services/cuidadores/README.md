# Servicio: cuidadores

## Responsable
Jeremy (S3N-SHI)

## Qué hace este servicio
Administra el directorio de cuidadores responsables del cuidado de los
adultos mayores: datos de contacto, relación con el paciente, horario
disponible, nivel de responsabilidad y qué pacientes tiene asignados.
No guarda información clínica ni de medicamentos; eso pertenece a otros
servicios (`informacion-salud`, `medicamentos`).

## Puerto
Este servicio corre internamente en el puerto **8080** (dentro del contenedor).
Puerto expuesto al host: **8099**

## Endpoints

| Método | Ruta                                   | Descripción                              |
|--------|-----------------------------------------|-------------------------------------------|
| GET    | /health                                 | Verifica que el servicio esté vivo        |
| POST   | /api/cuidadores                         | Crea un cuidador                          |
| GET    | /api/cuidadores                         | Lista todos los cuidadores                |
| GET    | /api/cuidadores/{id}                    | Consulta un cuidador por ID               |
| PUT    | /api/cuidadores/{id}                    | Actualiza un cuidador                     |
| DELETE | /api/cuidadores/{id}                    | Elimina un cuidador                       |
| GET    | /api/cuidadores/paciente/{pacienteId}   | Lista los cuidadores de un paciente       |

## Ejemplo de solicitud (POST /api/cuidadores)

```json
{
  "nombre": "Jeremy",
  "telefono": "0999999999",
  "email": "jeremy@example.com",
  "relacion": "Hijo",
  "horario_disponible": "Lunes a Viernes 08:00-18:00",
  "pacientes": ["paciente-1"],
  "nivel_responsabilidad": "principal"
}
```

## Variables de entorno

| Variable | Descripción                  | Ejemplo |
|----------|-------------------------------|---------|
| PORT     | Puerto interno del servicio   | 8080    |

## Cómo correrlo solo (sin docker-compose)

```bash
cd services/cuidadores
go mod download
go run .
```

## Cómo correrlo con Docker

```bash
docker build -t cuidadores .
docker run -p 8099:8080 cuidadores
```

## Pruebas

```bash
go test ./...
```
