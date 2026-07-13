# Etapa de compilación
FROM golang:1.22-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o informacion-salud .

# Etapa final (imagen liviana, sin herramientas de compilación)
FROM alpine:latest

WORKDIR /app
COPY --from=build /app/informacion-salud .

EXPOSE 8082

CMD ["./informacion-salud"]