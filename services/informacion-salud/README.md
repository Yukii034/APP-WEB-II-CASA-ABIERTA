# Servicio: [Informacion de Sauld]

## Responsable
Equipo: [Nahim Simba, Jostin Alvarado, Daivelyn Pincay, Joseph Paredes, Cristina Cedeño]

## Qué hace este servicio
Microservicio encargado de gestionar la información de salud del adulto mayor, almacenando datos como diagnósticos, alergias, enfermedades crónicas y antecedentes médicos. Proporciona una API REST para consultar y actualizar esta información, permitiendo que otros microservicios accedan a ella cuando sea necesario.

## Puerto
Este servicio corre internamente en el puerto **8082** (dentro del contenedor).
Puerto expuesto al host: **8082**

## Endpoints

| Método | Ruta                            | Descripción                                              |
|--------|---------------------------------|-----------------------------------------------------------|
| GET    | /health                         | Verifica que el servicio esté vivo                        |
| GET    | /api/informacion-salud          | Lista todas las fichas de salud registradas               |
| POST   | /api/informacion-salud          | Crea una nueva ficha de salud                              |
| GET    | /api/informacion-salud/{id}     | Consulta la ficha de salud de un paciente específico       |
| PUT    | /api/informacion-salud/{id}     | Actualiza (parcialmente) la ficha de salud de un paciente  |

### Cuerpo esperado (POST / PUT)

```json
{
  "nombre_paciente": "María Pérez",
  "diagnosticos": ["hipertensión"],
  "alergias": ["penicilina"],
  "enfermedades_cronicas": ["diabetes tipo 2"],
  "antecedentes_medicos": ["cirugía de cadera 2019"]
}
```

En `PUT`, los campos que no se envían (o se envían como `null`) no se sobrescriben; se conserva el valor anterior. Esto permite actualizar, por ejemplo, solo las alergias sin tener que reenviar todo lo demás.

### Ejemplo de respuesta

```json
{
  "id": "1",
  "nombre_paciente": "María Pérez",
  "diagnosticos": ["hipertensión"],
  "alergias": ["penicilina"],
  "enfermedades_cronicas": ["diabetes tipo 2"],
  "antecedentes_medicos": ["cirugía de cadera 2019"],
  "actualizado_en": "2026-07-12T11:00:00Z"
}
```

## Variables de entorno

| Variable | Descripción | Ejemplo |
|----------|-------------|---------|
| PORT     | Puerto interno del servicio | 8082 |

## Persistencia

Los datos se guardan **en memoria** (protegidos con un mutex, detrás de una interfaz `repository.Repository`), siguiendo la limitación conocida del proyecto (ver `docs/arquitectura.md`). Se pierden si el contenedor se reinicia. Cambiar a una base de datos real implicaría solo agregar una nueva implementación de `Repository`, sin tocar `service` ni `handler`.

## Estructura del servicio

```
services/informacion-salud/
├── main.go                          # wiring: arma repo -> service -> handler y levanta el router
├── internal/
│   ├── model/informacion_salud.go       # structs InformacionSalud y EntradaInformacionSalud
│   ├── repository/memoria.go            # interfaz Repository + implementación en memoria
│   ├── service/informacion_salud.go     # lógica de negocio (funciones puras, testeadas sin mocks)
│   └── handler/informacion_salud.go     # handlers HTTP (decodifican/codifican JSON, llaman al service)
├── go.mod / go.sum
└── Dockerfile
```

## Cómo correrlo solo (sin docker-compose)

```bash
cd services/informacion-salud
go run main.go
```

## Cómo correrlo con Docker

```bash
docker build -t [informacion-salud] .
docker run -p 8082:8080 [informacion-salud]
```