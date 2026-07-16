# Etapa de compilación
FROM golang:1.22-alpine AS build

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN go build -o servicio .

# Etapa final (imagen liviana, sin herramientas de compilación)
FROM alpine:latest

WORKDIR /app
COPY --from=build /app/servicio .

EXPOSE 8080

CMD ["./servicio"]