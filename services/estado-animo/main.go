package main

import (
	"log"
	"net/http"
	"os"

	"cuidabien/estado-animo/handlers"
	"cuidabien/estado-animo/storage"
)

func main() {

	storage.InicializarDatos()

	http.HandleFunc("/health", handlers.HealthHandler)
	http.HandleFunc("/api/estado-animo", handlers.EstadoAnimoHandler)
	http.HandleFunc("/api/estado-animo/alertas", handlers.AlertasHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8087"
	}

	log.Printf("Servicio de Estado de Ánimo corriendo en el puerto %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}