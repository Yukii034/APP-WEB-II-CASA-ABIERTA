# Arquitectura de CuidaBien

## Visión general

CuidaBien está construido como un conjunto de **microservicios independientes**, cada uno responsable de una funcionalidad específica del sistema. Todos los servicios están escritos en **Go**, empaquetados en contenedores **Docker**, y se comunican entre sí mediante peticiones **HTTP/REST**.

## Diagrama de componentes

```
                    ┌─────────────────────┐
                    │  Frontend/Dashboard │
                    └──────────┬──────────┘
                               │
                               ▼
                        ┌─────────────┐
                        │   Gateway   │  (puerto 8080)
                        │  (proxy /   │
                        │  agregador) │
                        └──────┬──────┘
                               │
        ┌──────────────┬──────┼──────┬──────────────┐
        ▼              ▼      ▼      ▼              ▼
 ┌─────────────┐┌───────────┐ ┌───────────┐ ┌───────────────┐
 │Medicamentos ││Emergencia │ │Monitoreo  │ │ ... otros     │
 │(8081)       ││(8082)     │ │(8083)     │ │ servicios     │
 └─────────────┘└─────┬─────┘ └───────────┘ └───────────────┘
                       │
                       ▼
                ┌─────────────┐
                │  Contactos  │  (llamado directo,
                │  (8084)     │   sin pasar por Gateway)
                └─────────────┘
```

## Componentes

### Gateway
- Es el único punto de entrada para el frontend o para quien consuma el sistema desde afuera.
- No contiene lógica de negocio propia: su función es recibir peticiones, redirigirlas al microservicio correspondiente, y en algunos casos combinar datos de varios servicios antes de responder.
- Conoce las URLs de todos los demás servicios a través de variables de entorno.

### Microservicios (medicamentos, emergencia, monitoreo, etc.)
- Cada uno maneja su propia lógica y expone endpoints REST propios.
- Todos exponen un endpoint `/health` usado para verificar que el servicio esté activo.
- No tienen por qué conocerse entre sí, **salvo que exista una dependencia lógica directa** (ver siguiente sección).

## Comunicación entre servicios

### Vía Gateway (patrón por defecto)
El frontend nunca llama directamente a un microservicio. Siempre pasa por el Gateway, que redirige la petición:

```
Frontend → Gateway → Medicamentos → Gateway → Frontend
```

### Comunicación directa entre microservicios (caso especial)
Cuando un servicio necesita datos de otro por una dependencia lógica —por ejemplo, **Emergencia** necesita consultar a **Contactos** para saber a quién avisar—, se permite la llamada directa entre ellos, sin pasar por el Gateway:

```
Emergencia → Contactos
```

Esto se hace igual que cualquier llamada HTTP, usando el nombre del servicio definido en `docker-compose.yml` como si fuera un dominio:

```go
resp, err := http.Get(os.Getenv("CONTACTOS_URL") + "/api/contactos")
```

### Resolución de nombres
Docker Compose crea una red interna (`cuidabien-net`) donde cada servicio es accesible por el nombre que se le dio en el `docker-compose.yml`. No se usa `localhost` para comunicación entre contenedores; `localhost` dentro de un contenedor apunta al contenedor mismo, no a otro servicio.

### Variables de entorno
Ninguna URL entre servicios se escribe fija en el código. Se definen en `docker-compose.yml` y se leen en Go con `os.Getenv()`. Esto permite cambiar direcciones o puertos sin tocar el código fuente.

## Justificación del diseño

- **Independencia entre equipos**: al tener cada servicio su propio módulo Go y Dockerfile, los 30 estudiantes pueden trabajar en paralelo sin pisarse el código.
- **Simplicidad**: no se usa un bus de mensajes ni colas complejas, solo HTTP/REST, adecuado para el tiempo disponible (una semana) y el nivel del curso.
- **Escalabilidad futura**: si el proyecto creciera, este mismo patrón permite reemplazar Docker Compose por Kubernetes, o cambiar llamadas HTTP directas por una cola de mensajes, sin rediseñar todo el sistema desde cero.

## Limitaciones conocidas (versión actual)

- Los datos se manejan en memoria dentro de cada servicio (no hay base de datos persistente).
- No hay autenticación ni cifrado entre servicios.
- Las notificaciones (medicamentos, emergencias) son simuladas, no se envían SMS/correo reales.

Estas limitaciones son aceptables para el alcance de la casa abierta, ya que el objetivo es demostrar el funcionamiento de la arquitectura de microservicios, no un producto listo para producción.