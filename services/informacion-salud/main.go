package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"cuidabien/informacion-salud/internal/handler"
	"cuidabien/informacion-salud/internal/repository"
	"cuidabien/informacion-salud/internal/service"
)

func main() {
	repo := repository.NuevaMemoriaRepository()
	repository.Sembrar(repo)

	svc := service.NuevoInformacionSaludService(repo)
	h := handler.NuevoInformacionSaludHandler(svc)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", handler.Health)

	r.Route("/api/informacion-salud", func(r chi.Router) {
		r.Get("/", h.Listar)
		r.Post("/", h.Crear)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.Obtener)
			r.Put("/", h.Actualizar)
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("Servicio de información de salud corriendo en el puerto %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
