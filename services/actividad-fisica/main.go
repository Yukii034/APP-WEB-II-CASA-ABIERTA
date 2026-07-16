package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"cuidabien/actividad-fisica/internal/handler"
	"cuidabien/actividad-fisica/internal/repository"
	"cuidabien/actividad-fisica/internal/service"
	"cuidabien/actividad-fisica/internal/middleware"
)

func main() {
	repo := repository.NuevaMemoriaRepository()
	repository.Sembrar(repo)
	svc := service.NuevoActividadFisicaService(repo)
	h := handler.NuevoActividadFisicaHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.Health)
	
	mux.HandleFunc("/api/actividad-fisica", middleware.Logger(func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            h.Listar(w, r)
        case http.MethodPost:
            h.Crear(w, r)
        default:
            http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
        }
    }))

	mux.HandleFunc("/api/actividad-fisica/",  middleware.Logger(func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/api/actividad-fisica/")
		if id == "" || strings.Contains(id, "/") {
			http.NotFound(w, r)
			return
		}
		switch r.Method {
		case http.MethodGet:
			h.Obtener(w, id)
		case http.MethodPut:
			h.Actualizar(w, r, id)
		case http.MethodDelete:
			h.Eliminar(w, id)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	}))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Servicio de actividad física corriendo en el puerto %s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
