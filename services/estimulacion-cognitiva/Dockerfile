FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN go build -o estimulacion-cognitiva .

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/estimulacion-cognitiva .

EXPOSE 8080

CMD ["./estimulacion-cognitiva"]
