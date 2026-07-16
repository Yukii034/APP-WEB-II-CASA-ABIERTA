# syntax=docker/dockerfile:1

FROM golang:1.22-alpine AS builder

WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/monitoreo-signos-vitales ./cmd/API

FROM alpine:3.20

RUN addgroup -S app \
    && adduser -S -G app app

WORKDIR /app

COPY --from=builder /out/monitoreo-signos-vitales ./monitoreo-signos-vitales

USER app

ENV PORT=8080

EXPOSE 8080

ENTRYPOINT ["/app/monitoreo-signos-vitales"]
