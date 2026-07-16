# Primera etapa: compilar el microservicio
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copia primero el archivo del módulo
COPY go.mod ./

# Descarga las dependencias declaradas
RUN go mod download

# Copia todos los archivos del microservicio
COPY . .

# Compila el proyecto para Linux
RUN CGO_ENABLED=0 GOOS=linux go build -o recordatorio-medicamentos .

# Segunda etapa: imagen ligera para ejecutar
FROM alpine:3.20

WORKDIR /app

# Copia únicamente el ejecutable compilado
COPY --from=builder /app/recordatorio-medicamentos ./recordatorio-medicamentos

# Puerto interno utilizado por el microservicio
EXPOSE 8080

# Comando que inicia el servicio
CMD ["./recordatorio-medicamentos"]