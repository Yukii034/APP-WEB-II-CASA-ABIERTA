# CuidaBien

Sistema de microservicios enfocado en el cuidado y bienestar de adultos mayores, desarrollado en **Go** y desplegado con **Docker**. Proyecto de casa abierta para la materia de Aplicaciones Web 2.

## Descripción

CuidaBien está compuesto por múltiples microservicios independientes (recordatorio de medicamentos, contacto de emergencia, monitoreo de signos vitales, entre otros) que se comunican entre sí a través de peticiones REST, coordinados por un API Gateway central.

## Arquitectura

```
Frontend / Dashboard
        │
        ▼
   ┌─────────┐
   │ Gateway │  ← punto de entrada único
   └────┬────┘
        │
   ┌────┼────┬──────────┬───────────┐
   ▼    ▼    ▼           ▼
Medicamentos Emergencia Monitoreo  ... (otros servicios)
```

- Cada microservicio vive en su propia carpeta dentro de `services/`, con su propio `go.mod` y `Dockerfile`.
- Los servicios se comunican entre sí usando su **nombre** definido en `docker-compose.yml` (Docker resuelve esto automáticamente, no se usa `localhost` entre contenedores).
- El **Gateway** es el punto de entrada que consulta a los demás servicios y expone los datos hacia el frontend.
- Las URLs de otros servicios **nunca se hardcodean**, se pasan por variables de entorno.

Ver más detalle en [`docs/arquitectura.md`](docs/arquitectura.md).

## Cómo levantar el proyecto completo

Requisitos: tener Docker y Docker Compose instalados.

```bash
git clone https://github.com/Yukii034/APP-WEB-II-CASA-ABIERTA.git
docker-compose up --build
```

Esto levanta todos los servicios definidos en `docker-compose.yml`. El gateway queda disponible en `http://localhost:8080`.

## Cómo agregar un nuevo microservicio

1. **Copiar la plantilla base:**
   ```bash
   cp -r services/_template services/nombre-de-tu-servicio
   ```

2. **Renombrar el módulo** en `services/nombre-de-tu-servicio/go.mod`:
   ```
   module cuidabien/nombre-de-tu-servicio

   go 1.22
   ```

3. **Reemplazar la lógica de ejemplo** en `main.go` (el struct `Item` y el handler `itemsHandler`) por la lógica real de tu servicio. **No elimines el endpoint `/health`**, es usado para verificar que el servicio esté vivo.

4. **Agregar tu servicio al `docker-compose.yml`** de la raíz, copiando el bloque de `medicamentos` como referencia:
```yaml
   nombre-de-tu-servicio:
     build: ./services/nombre-de-tu-servicio
     environment:
       - PORT=8080
     networks:
       - cuidabien-net
     ports:
       - "808X:8080"   # asignar un puerto libre, ver tabla abajo
```

5. **Agregar tu servicio a la matrix del CI**, en `.github/workflows/ci.yml`:
```yaml
   service:
     - medicamentos
     - gateway
     - nombre-de-tu-servicio   # agregar aquí
```
   ⚠️ **Importante:** solo agrégalo aquí cuando la carpeta de tu servicio ya exista en el repo con su `main.go`, `go.mod` y `Dockerfile` funcionando. Si agregas el nombre antes de crear la carpeta, el CI va a fallar porque no encuentra la ruta.

6. **Si tu servicio necesita datos de otro**, sigue el patrón usado en `services/gateway/main.go`...

7. **Completar el `README.md`** dentro de tu carpeta de servicio con: qué hace, endpoints disponibles, puerto y variables de entorno.

8. **Abrir un Pull Request hacia `main`**. No se puede mergear directo, debe pasar el CI (build, vet y test) y al menos una revisión aprobada.

## Convención de puertos

| Servicio      | Puerto host | Estado |
|---------------|-------------|--------|
| gateway       | 8080        | ✅ activo |
| informacion salud    | 8082        | ✅ activo |
| monitoreo     | 8083        | pendiente |
| alimentacion   | 8084        | ✅ activo |
| citas médicas | 8085        | ✅ activo |
| reportes médicos | 8086        | ✅ activo |
| estado animo     | 8087 | ✅ activo |
| Recordatorio de medicamentos     | 8089 | ✅ activo  |

> Antes de asignarte un puerto, revisa esta tabla y actualízala en tu PR para evitar choques con otro equipo.

## Convenciones del repositorio

**Ramas:**
- `main`: protegida, solo se actualiza vía Pull Request aprobado y con CI en verde.
- `feature/nombre-del-servicio`: una rama por servicio/equipo, ej. `feature/emergencia`.

**Commits:** mensajes cortos y descriptivos en español, ej:
```
feat: agrega endpoint de recordatorio de medicamentos
fix: corrige puerto expuesto en docker-compose
docs: actualiza README de emergencia
```

**Pull Requests:**
- Deben pasar el pipeline de CI (`go build`, `go vet`, `go test`) antes de poder mergearse.
- Requieren al menos 1 aprobación.
- Describir brevemente qué se agregó y cómo probarlo.

## Estructura del repositorio

```
cuidabien/
├── .github/workflows/ci.yml    ← pipeline de CI (build, vet, test)
├── services/
│   ├── _template/                ← plantilla base, copiar para cada servicio nuevo
│   ├── gateway/                  ← API Gateway, punto de entrada único
│   ├── medicamentos/             ← ejemplo funcional
│   └── ...                       ← nuevos servicios de cada equipo
├── docker-compose.yml
├── .env.example
├── .gitignore
└── docs/
    └── arquitectura.md
```

## Equipo

| Servicio | Responsable(s) |
|----------|----------------|
| Gateway / base del repo | [Pierina Peñaherrera] |
| Informacion de salud | [Nahim Simba, Jostin Alvarado, Daivelyn Pincay, Joseph Paredes, Cristina Cedeño] |
| Emergencia | [nombre] |
| Alimentos | [Eduardo Lopez, Pierina Peñaherrera, José Manuel Castillo, Néstor Gallegos] |
| Reportes médicos | [Anthony Mendoza - Deimuz, Holguin Nathaly Jasmin, Cedeño Geovanny Alexander] |
| Citas médicas | [Anthony Mendoza - Deimuz, Holguin Nathaly Jasmin, Cedeño Geovanny Alexander] |
| Estado animo |  [Danny Zambrano, Michael Cedeño] |
| Recordatorio de medicamentos | [Manuel Intriago, Madelyn Zambrano, Michelle Salazar] |

