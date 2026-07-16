package main

import (
	"log"
	"net/http"

	"cuidabien/alimentacion/config"
	"cuidabien/alimentacion/handlers"
	"cuidabien/alimentacion/httpserver"
	"cuidabien/alimentacion/middleware"
	"cuidabien/alimentacion/service"
	"cuidabien/alimentacion/storage"
)

func main() {
	cfg := config.Load()

	store := storage.NewStore()
	svc := service.New(store)

	svc.ConfigurarHoraLimite("desayuno", cfg.DesayunoHasta)
	svc.ConfigurarHoraLimite("almuerzo", cfg.AlmuerzoHasta)
	svc.ConfigurarHoraLimite("cena", cfg.CenaHasta)

	h := handlers.New(svc)
	mux := httpserver.New(h)
	wrapped := middleware.Logger(mux)

	log.Printf("Servicio de alimentación corriendo en el puerto %s", cfg.Puerto)
	log.Fatal(http.ListenAndServe(":"+cfg.Puerto, wrapped))
}
