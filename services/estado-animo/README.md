# Servicio de Estado de Ánimo

## Responsable
Equipo: Zambrano Mera Danny : Cedeño Pincay Michael

## Qué hace este servicio
Registra diariamente el estado de ánimo del adulto mayor utilizando una escala numérica del 1 al 5, una etiqueta emocional y comentarios de texto opcionales. Cuenta con lógica interna capaz de emitir alertas si se registra un decaimiento anímico persistente en los últimos dos días.

## Puerto
Este servicio corre internamente en el puerto *8080* (dentro del contenedor).
Puerto expuesto al host: *8087*

## Endpoints

| Método | Ruta                      | Descripción                                                  |
|--------|---------------------------|--------------------------------------------------------------|
| GET    | /health                   | Verifica que el servicio esté vivo                           |
| GET    | /api/estado-animo         | Obtiene el historial de registros de estado de ánimo         |
| POST   | /api/estado-animo         | Registra un nuevo estado de ánimo diario                     |
| GET    | /api/estado-animo/alertas | Analiza el historial y devuelve una alerta si hay desánimo   |

## Variables de entorno

| Variable | Descripción | Ejemplo |
|----------|-------------|---------|
| PORT     | Puerto interno del servicio | 8080 |

## Pruebas de Endpoints (Manuales)
Puedes probar manualmente el flujo del servicio ejecutando estos comandos en tu terminal:

### 1. Health Check
bash
curl -X GET http://localhost:8087/health

### Obtener historial

curl -X GET http://localhost:8087/api/estado-animo


### Registrar Estado de Ánimo (POST)

curl -X POST http://localhost:8087/api/estado-animo \
     -H "Content-Type: application/json" \
     -d '{"nivel": 1, "emocion": "Triste", "comentario": "Me siento muy desanimado"}'


### Consultar Activación de Alertas

curl -X GET http://localhost:8087/api/estado-animo/alertas


### Pruebas Unitarias (Automáticas)

go test ./... -v


## Cómo correrlo solo (sin docker-compose)

bash
cd services/estado-animo
go run main.go


## Cómo correrlo solo con docker-compose

docker build -t estado-animo .
docker run -p 8087:8080 estado-animo