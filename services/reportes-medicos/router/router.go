package router

import (
	"cuidabien/reportes-medicos/handlers"
	"net/http"
	"strings"
)

func New(h *handlers.Handlers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if path == "/health" {
			h.HealthHandler(w, r)
			return
		}

		if path == "/api/reportes-medicos/resumen" {
			h.ResumenHandler(w, r)
			return
		}

		if path == "/api/reportes-medicos/semanal" {
			h.SemanalHandler(w, r)
			return
		}

		if strings.HasPrefix(path, "/api/reportes-medicos/paciente/") {
			h.PacienteHandler(w, r)
			return
		}

		if path == "/api/reportes-medicos" {
			h.SemanalHandler(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"Ruta no encontrada"}`))
	}
}
