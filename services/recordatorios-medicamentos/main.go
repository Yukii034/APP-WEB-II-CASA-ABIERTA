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
	// Crea el repositorio en memoria.
	repo := repository.NuevaMemoriaRepository()

	// Carga recordatorios iniciales para la demostración.
	repository.Sembrar(repo)

	// Crea la capa de servicio.
	svc := service.NuevoRecordatorioMedicamentoService(repo, os.Stdout)

	// Crea la capa handler.
	h := handler.NuevoRecordatorioMedicamentoHandler(svc)

	// Obtiene la zona horaria desde las variables de entorno.
	zonaHoraria := os.Getenv("TIMEZONE")
	if zonaHoraria == "" {
		zonaHoraria = "America/Guayaquil"
	}

	ubicacion, err := time.LoadLocation(zonaHoraria)
	if err != nil {
		log.Fatalf("Zona horaria inválida: %v", err)
	}

	// Inicia la verificación automática de los horarios.
	go svc.IniciarVerificador(
		context.Background(),
		30*time.Second,
		ubicacion,
	)

	// Crea el router.
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Endpoint de salud.
	r.Get("/health", handler.Health)

	// Endpoints de recordatorios de medicamentos.
	r.Get("/api/recordatorios-medicamentos", h.Listar)
	r.Post("/api/recordatorios-medicamentos", h.Crear)
	r.Post("/api/recordatorios-medicamentos/verificar", h.Verificar)
	r.Get("/api/recordatorios-medicamentos/{id}", h.Obtener)
	r.Put("/api/recordatorios-medicamentos/{id}", h.Actualizar)
	r.Delete("/api/recordatorios-medicamentos/{id}", h.Eliminar)
	r.Patch("/api/recordatorios-medicamentos/{id}/estado", h.CambiarEstado)

	// Obtiene el puerto desde las variables de entorno.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servicio de recordatorios de medicamentos corriendo en el puerto %s", port)
	log.Printf("Verificador automático usando la zona horaria %s", zonaHoraria)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
