package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"
	_ "time/tzdata"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"cuidabien/recordatorios-medicamentos/internal/handler"
	"cuidabien/recordatorios-medicamentos/internal/repository"
	"cuidabien/recordatorios-medicamentos/internal/service"
)

func main() {
	repo := repository.NuevaMemoriaRepository()

	repository.Sembrar(repo)

	svc := service.NuevoRecordatorioMedicamentoService(
		repo,
		os.Stdout,
	)

	h := handler.NuevoRecordatorioMedicamentoHandler(
		svc,
	)

	zonaHoraria := os.Getenv("TIMEZONE")

	if zonaHoraria == "" {
		zonaHoraria = "America/Guayaquil"
	}

	ubicacion, err := time.LoadLocation(
		zonaHoraria,
	)

	if err != nil {
		log.Fatalf(
			"Zona horaria inválida: %v",
			err,
		)
	}

	go svc.IniciarVerificador(
		context.Background(),
		30*time.Second,
		ubicacion,
	)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", handler.Health)

	r.Route(
		"/api/recordatorio-medicamentos",
		func(r chi.Router) {
			r.Get("/", h.Listar)
			r.Post("/", h.Crear)
			r.Post("/verificar", h.Verificar)

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", h.Obtener)
				r.Put("/", h.Actualizar)
				r.Delete("/", h.Eliminar)
				r.Patch("/estado", h.CambiarEstado)
			})
		},
	)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	log.Printf(
		"Servicio de recordatorio de medicamentos "+
			"corriendo en el puerto %s",
		port,
	)

	log.Printf(
		"Verificador automático usando la zona horaria %s",
		zonaHoraria,
	)

	log.Fatal(
		http.ListenAndServe(
			":"+port,
			r,
		),
	)
}
