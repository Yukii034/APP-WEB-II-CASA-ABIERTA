package main

import (
	"cuidabien/reportes-medicos/handlers"
	"cuidabien/reportes-medicos/router"
	"cuidabien/reportes-medicos/storage"
	"log"
	"net/http"
	"os"
)

func main() {
	store := storage.NewStore()
	h := handlers.New(store)
	r := router.New(h)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servicio de reportes medicos corriendo en el puerto %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
