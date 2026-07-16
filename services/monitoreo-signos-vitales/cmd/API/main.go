package main

import (
	"log"
	"net/http"

	"monitoreo-signos-vitales/internal/config"
	"monitoreo-signos-vitales/internal/handlers"
	"monitoreo-signos-vitales/internal/httpserver"
	"monitoreo-signos-vitales/internal/service"
	"monitoreo-signos-vitales/internal/storage"
)

func main() {
	repo := storage.NuevaMemoriaRepository()
	storage.Sembrar(repo)
	servicio := service.NuevoSignosVitalesService(repo)
	router := httpserver.NuevoRouter(handlers.NuevoSignosVitalesHandler(servicio))
	puerto := config.Puerto()
	log.Printf("Monitoreo de signos vitales escuchando en el puerto %s", puerto)
	log.Fatal(http.ListenAndServe(":"+puerto, router))
}
