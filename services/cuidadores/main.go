package main

import (
	"log"
	"net/http"
	"os"

	"cuidabien/cuidadores/internal/handler"
	"cuidabien/cuidadores/internal/repository"
	"cuidabien/cuidadores/internal/router"
	"cuidabien/cuidadores/internal/service"
)

func main() {
	repo := repository.NuevaMemoriaRepository()
	svc := service.NuevoCuidadorService(repo)
	h := handler.NuevoCuidadorHandler(svc)
	mux := router.Nuevo(h)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // puerto interno por defecto dentro del contenedor
	}

	log.Printf("Servicio de cuidadores corriendo en el puerto %s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
